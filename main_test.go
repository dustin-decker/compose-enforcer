package main

import "testing"

var validations = Validations{
	StoragePath:       "/dockerstorage",
	QualifiedService:  "some/service",
	DockerWrite:       []string{"SwarmCommand"},
	Secrets:           []string{"secret1", "secret2"},
	Networks:          []string{"lol"},
	MemoryLimit:       "3G",
	CPULimit:          "4",
	MemoryReservation: "3G",
	CPUReservation:    "4",
}

func TestVolumes(t *testing.T) {
	config, err := LoadConfigFile("test.yml")
	if err != nil {
		t.Error("Failed to load test.yml compose file")
	}

	t.Run("A=1", func(t *testing.T) {
		for _, Service := range config.Services {
			Service.Volumes[0].Source = "/var/run/docker.sock:/var/run/docker.sock:ro"
			ValidateVolumes(validations, Service)
		}
	})
}
