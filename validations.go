package main

import (
	"errors"
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/docker/cli/cli/compose/types"

	units "github.com/docker/go-units"
	"github.com/sevoma/goutil"
)

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
func ValidateVolumes(Validations Validations, Service types.ServiceConfig) error {
	for _, Volume := range Service.Volumes {

		// Only allow reads to docker.sock
		// Be aware that this still allows inspecting env vars - which may contain secrets
		// We don't use env vars for secrets - use docker secrets management instead
		if Volume.Source == "/var/run/docker.sock" {
			if Volume.ReadOnly == false {
				return errors.New("Docker socket mount must be read only")
			}
			continue
		}

		// Don't allow relative paths
		if !path.IsAbs(Volume.Source) {
			return errors.New("Volume paths must be absolute")
		}

		// Ensure that any other bind mounts are where they're allowed for that
		// service.network
		allowedVolumePath := path.Join("/dockerstorage", "service-name")
		// Collapses traversals ../../
		collapsedVolumePath := path.Join(Volume.Source)
		if !strings.HasPrefix(collapsedVolumePath, allowedVolumePath) {
			return fmt.Errorf("Volume mounts must be in '%s'", allowedVolumePath)
		}

	}
	return nil
}

// ValidateSecrets validates that a service's secrets are permitted for the service
func ValidateSecrets(Validations Validations, Service types.ServiceConfig) error {
	for _, Secret := range Service.Secrets {
		if !goutil.StringInSlice(Secret.Source, Validations.Secrets) {
			return fmt.Errorf("Secret '%s' not in the whitelist", Secret.Source)
		}
	}
	return nil
}

// ValidateNetworks validates that a service's secrets are permitted for the service
func ValidateNetworks(Validations Validations, Service types.ServiceConfig) error {
	for Network := range Service.Networks {
		if !goutil.StringInSlice(Network, Validations.Networks) {
			return fmt.Errorf("Network '%s' not in the whitelist", Network)
		}
	}
	return nil
}

// ValidateResources validates that a service's resources and limits specified are sensible
func ValidateResources(Validations Validations, Service types.ServiceConfig) error {

	// Ensure mem limit does not exceed our service max configured
	if err := meetsMemoryConstraint(Service.Deploy.Resources.Limits.MemoryBytes,
		Validations.MemoryLimit); err != nil {
		return err
	}

	// Ensure mem reservation does not exceed our service max configured
	if err := meetsMemoryConstraint(Service.Deploy.Resources.Reservations.MemoryBytes,
		Validations.MemoryReservation); err != nil {
		return err
	}

	// Ensure CPU limit does not exceed our service max configured
	if err := meetsCPUConstraint(Service.Deploy.Resources.Limits.NanoCPUs,
		Validations.CPULimit); err != nil {
		return err
	}

	// Ensure CPU reservation does not exceed our service max configured
	if err := meetsCPUConstraint(Service.Deploy.Resources.Reservations.NanoCPUs,
		Validations.CPUReservation); err != nil {
		return err
	}

	return nil
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

func meetsCPUConstraint(cpu string, cpuAllowed string) error {
	cpuAllowedFloat, err := strconv.ParseFloat(cpuAllowed, 64)
	if err != nil {
		return fmt.Errorf("Invalid CPUs provided")
	}

	cpuFloat, err := strconv.ParseFloat(cpu, 64)
	if err != nil {
		return fmt.Errorf("Invalid CPUs provided")
	}

	// nanoCPU := cpuAllowedFloat / math.Pow(10, -9)

	if cpuFloat > cpuAllowedFloat {
		return fmt.Errorf("Please keep memory limit <= %s", cpuAllowed)
	}
	return nil
}

// ValidateConfig ensures that the provided config follows our rules
func ValidateConfig(Validations Validations, config *types.Config) error {
	for _, Service := range config.Services {

		err := ValidateVolumes(Validations, Service)
		if err != nil {
			return err
		}

		err = ValidateSecrets(Validations, Service)
		if err != nil {
			return err
		}

		err = ValidateNetworks(Validations, Service)
		if err != nil {
			return err
		}

		err = ValidateResources(Validations, Service)
		if err != nil {
			return err
		}
		fmt.Println(Service.NetworkMode)

	}

	return nil
}
