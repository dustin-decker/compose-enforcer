package main

import "testing"

var validations = Validations{
	StoragePath:       "/dockerstorage",
	QualifiedService:  "some/service",
	DockerWrite:       []string{"SwarmCommand"},
	Secrets:           []string{"secret1", "secret2"},
	Networks:          []string{"network1"},
	MemoryLimit:       "3G",
	CPULimit:          "4",
	MemoryReservation: "3G",
	CPUReservation:    "4",
}

func TestValidateVolumes(t *testing.T) {
	config, err := LoadConfigFile("test.yml")
	if err != nil {
		t.Error("Failed to load test.yml compose file")
	}

	t.Run("Docker sock must be read only", func(t *testing.T) {
		for _, Service := range config.Services {
			Service.Volumes[0].Source = "/var/run/docker.sock"
			Service.Volumes[0].ReadOnly = false
			err := ValidateVolumes(validations, Service)
			if err == nil {
				t.Errorf("Failed")
			}
		}
	})

	t.Run("Docker sock must be read only", func(t *testing.T) {
		for _, Service := range config.Services {
			Service.Volumes[0].Source = "/var/run/docker.sock"
			Service.Volumes[0].ReadOnly = true
			err := ValidateVolumes(validations, Service)
			if err != nil {
				t.Error(err)
			}
		}
	})

	t.Run("Volumes cannot have a relative path", func(t *testing.T) {
		for _, Service := range config.Services {
			Service.Volumes[0].Source = "run/docker.sock"
			err := ValidateVolumes(validations, Service)
			if err == nil {
				t.Error("Failed")
			}
		}
	})

	t.Run("Must use allow volume mount path", func(t *testing.T) {
		for _, Service := range config.Services {
			Service.Volumes[0].Source = "/some/path"
			err := ValidateVolumes(validations, Service)
			if err == nil {
				t.Error("Failed")
			}
		}
	})

}

func TestValidateNetworks(t *testing.T) {
	config, err := LoadConfigFile("test.yml")
	if err != nil {
		t.Error("Failed to load test.yml compose file")
	}

	t.Run("Network must be whitelisted for the service", func(t *testing.T) {
		for _, Service := range config.Services {
			validations.Networks = []string{"notYoNetwork"}
			err := ValidateNetworks(validations, Service)
			if err == nil {
				t.Errorf("Failed")
			}
		}
	})

	t.Run("Network must be whitelisted for the service", func(t *testing.T) {
		for _, Service := range config.Services {
			validations.Networks = []string{"network1"}
			err := ValidateNetworks(validations, Service)
			if err != nil {
				t.Error(err)
			}
		}
	})

}
