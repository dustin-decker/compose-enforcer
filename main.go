package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/docker/docker/cli/compose/loader"
	"github.com/docker/docker/cli/compose/types"
)

func main() {
	var config *types.Config
	config, err := LoadConfigFile("dc.yml")
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

// Jacked pretty much the rest of this file from docker/cli/loader_test.go

// LoadConfigFile loads config from file!
func LoadConfigFile(filePath string) (*types.Config, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	config, err := LoadConfig(bytes)
	if err != nil {
		return nil, err
	}

	return config, err
}

// LoadConfig loads config from byte slice!
func LoadConfig(bytes []byte) (*types.Config, error) {

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
