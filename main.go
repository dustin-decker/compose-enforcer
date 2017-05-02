package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/docker/docker/cli/compose/loader"
	"github.com/docker/docker/cli/compose/types"
	units "github.com/docker/go-units"
	util "github.com/sevoma/goutil"
)

func main() {
	var config *types.Config
	config, err := LoadConfig("dc.yml")
	if err != nil {
		log.Fatal(err.Error())
	}

	dockerWrite := []string{"SwarmCommand"}
	Validations := Validations{
		StoragePath:       "/dockerstorage",
		QualifiedService:  "some/service",
		DockerWrite:       dockerWrite,
		Secrets:           []string{"secret1", "secret2"},
		Networks:          []string{"lol"},
		MemoryLimit:       "3G",
		CPULimit:          "4",
		MemoryReservation: "3G",
		CPUReservation:    "4",
	}
	_, err = ValidateConfig(Validations, config)
	if err != nil {
		log.Fatal(err.Error())
	}

}

// Validations struct provides details used for whitelisting configurations
type Validations struct {
	StoragePath       string
	QualifiedService  string
	DockerWrite       []string
	Secrets           []string
	Networks          []string
	MemoryLimit       string
	CPULimit          string
	MemoryReservation string
	CPUReservation    string
}

// ValidateVolumes validates a service's volumes to ensure they do what we allow
func ValidateVolumes(Validations Validations, Service types.ServiceConfig) (bool, error) {
	for _, Volume := range Service.Volumes {

		// Only allow reads to docker.sock
		// Be aware that this still allows inspecting env vars - which may contain secrets
		// We don't use env vars for secrets - use docker secrets management instead
		if Volume.Source == "/var/run/docker.sock" {
			if Volume.ReadOnly == false {
				return false, errors.New("Docker socket mount must be read only")
			}
			continue
		}

		// Don't allow relative paths
		if !path.IsAbs(Volume.Source) {
			return false, errors.New("Volume paths must be absolute")
		}

		// Ensure that any other bind mounts are where they're allowed for that
		// service.network
		allowedVolumePath := path.Join("/dockerstorage", "service-name")
		// Collapses traversals ../../
		collapsedVolumePath := path.Join(Volume.Source)
		if !strings.HasPrefix(collapsedVolumePath, allowedVolumePath) {
			return false, fmt.Errorf("Volume mounts must be in '%s'", allowedVolumePath)
		}

	}
	return true, nil
}

// ValidateSecrets validates that a service's secrets are permitted for the service
func ValidateSecrets(Validations Validations, Service types.ServiceConfig) (bool, error) {
	for _, Secret := range Service.Secrets {
		if !util.StringInSlice(Secret.Source, Validations.Secrets) {
			return false, fmt.Errorf("Secret '%s' not in the whitelist", Secret.Source)
		}
	}
	return true, nil
}

// ValidateNetworks validates that a service's secrets are permitted for the service
func ValidateNetworks(Validations Validations, Service types.ServiceConfig) (bool, error) {
	for Network := range Service.Networks {
		if !util.StringInSlice(Network, Validations.Networks) {
			return false, fmt.Errorf("Network '%s' not in the whitelist", Network)
		}
	}
	return true, nil
}

// ValidateResources validates that a service's resources and limits specified are sensible
func ValidateResources(Validations Validations, Service types.ServiceConfig) (bool, error) {

	// Ensure mem limit does not exceed our service max configured
	if err := meetsMemoryConstraint(Service.Deploy.Resources.Limits.MemoryBytes,
		Validations.MemoryLimit); err != nil {
		return false, err
	}

	// Ensure mem reservation does not exceed our service max configured
	if err := meetsMemoryConstraint(Service.Deploy.Resources.Reservations.MemoryBytes,
		Validations.MemoryReservation); err != nil {
		return false, err
	}

	return true, nil
}

func meetsMemoryConstraint(mem types.UnitBytes, memAllowed string) error {
	memAllowedBytes, err := units.RAMInBytes(memAllowed)
	if err != nil {
		return err
	}
	if mem > types.UnitBytes(memAllowedBytes) {
		return fmt.Errorf("Please keep memory limit <= %s", memAllowed)
	}
	return nil
}

// ValidateConfig ensures that the provided config follows our rules
func ValidateConfig(Validations Validations, config *types.Config) (bool, error) {
	for _, Service := range config.Services {

		_, err := ValidateVolumes(Validations, Service)
		if err != nil {
			return false, err
		}

		_, err = ValidateSecrets(Validations, Service)
		if err != nil {
			return false, err
		}

		_, err = ValidateNetworks(Validations, Service)
		if err != nil {
			return false, err
		}

		_, err = ValidateResources(Validations, Service)
		if err != nil {
			return false, err
		}
		fmt.Println(Service.NetworkMode)

	}

	return true, nil
}

// LoadConfig lol
func LoadConfig(filePath string) (*types.Config, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	homeDir := "/home/foo"
	env := map[string]string{"HOME": homeDir, "QUX": "qux_from_environment"}
	config, err := loadYAMLWithEnv(string(bytes), env)
	if err != nil {
		return nil, err
	}

	return config, err
}

func buildConfigDetails(source map[string]interface{}, env map[string]string) types.ConfigDetails {
	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return types.ConfigDetails{
		WorkingDir: workingDir,
		ConfigFiles: []types.ConfigFile{
			{Filename: "filename.yml", Config: source},
		},
		Environment: env,
	}
}

func loadYAML(yaml string) (*types.Config, error) {
	return loadYAMLWithEnv(yaml, nil)
}

func loadYAMLWithEnv(yaml string, env map[string]string) (*types.Config, error) {
	dict, err := loader.ParseYAML([]byte(yaml))
	if err != nil {
		return nil, err
	}

	return loader.Load(buildConfigDetails(dict, env))
}
