package main

import "testing"

var validations = Validations{
	StoragePath:       "/dockerstorage",
	QualifiedService:  "some/service",
	DockerWrite:       []string{"SwarmCommand"},
	Secrets:           []string{"secret1", "secret2"},
	Networks:          []string{"network1"},
	MemoryLimit:       "4G",
	CPULimit:          "4",
	MemoryReservation: "4G",
	CPUReservation:    "4",
}

func TestLoadConfigFile(t *testing.T) {
	_, err := LoadConfigFile("sdfjsdf.yml")
	if err == nil {
		t.Error("Loading non-existent file should fail")
	}
}

func TestLoadConfig(t *testing.T) {
	var empty []byte
	_, err := LoadConfig(empty)
	if err == nil {
		t.Error("Loading empty file should fail")
	}
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

func TestValidateSecrets(t *testing.T) {
	config, err := LoadConfigFile("test.yml")
	if err != nil {
		t.Error("Failed to load test.yml compose file")
	}

	t.Run("Network must be whitelisted for the service", func(t *testing.T) {
		for _, Service := range config.Services {
			validations.Secrets = []string{"notYoSecret"}
			err := ValidateSecrets(validations, Service)
			if err == nil {
				t.Errorf("Failed")
			}
		}
	})

	t.Run("Network must be whitelisted for the service", func(t *testing.T) {
		for _, Service := range config.Services {
			validations.Secrets = []string{"secret1", "secret2"}
			err := ValidateSecrets(validations, Service)
			if err != nil {
				t.Error(err)
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

func TestValidateResourcess(t *testing.T) {
	config, err := LoadConfigFile("test.yml")
	if err != nil {
		t.Error("Failed to load test.yml compose file")
	}

	t.Run("Memory limit", func(t *testing.T) {
		for _, Service := range config.Services {
			validations.MemoryLimit = "1M"
			err := ValidateResources(validations, Service)
			if err == nil {
				t.Errorf("Failed")
			}
		}
	})

	t.Run("Memory limit", func(t *testing.T) {
		for _, Service := range config.Services {
			validations.MemoryLimit = "999G"
			err := ValidateResources(validations, Service)
			if err != nil {
				t.Error(err)
			}
		}
	})

	t.Run("Memory reservation", func(t *testing.T) {
		for _, Service := range config.Services {
			validations.MemoryReservation = "1M"
			err := ValidateResources(validations, Service)
			if err == nil {
				t.Errorf("Failed")
			}
		}
	})

	t.Run("Memory reservation", func(t *testing.T) {
		for _, Service := range config.Services {
			validations.MemoryReservation = "999G"
			err := ValidateResources(validations, Service)
			if err != nil {
				t.Error(err)
			}
		}
	})

	t.Run("CPU limit", func(t *testing.T) {
		for _, Service := range config.Services {
			validations.CPULimit = "0.1"
			err := ValidateResources(validations, Service)
			if err == nil {
				t.Errorf("Failed")
			}
		}
	})

	t.Run("CPU limit", func(t *testing.T) {
		for _, Service := range config.Services {
			validations.CPULimit = "999"
			err := ValidateResources(validations, Service)
			if err != nil {
				t.Error(err)
			}
		}
	})

	t.Run("CPU reservation", func(t *testing.T) {
		for _, Service := range config.Services {
			validations.CPUReservation = "0.1"
			err := ValidateResources(validations, Service)
			if err == nil {
				t.Errorf("Failed")
			}
		}
	})

	t.Run("CPU reservation", func(t *testing.T) {
		for _, Service := range config.Services {
			validations.CPUReservation = "999"
			err := ValidateResources(validations, Service)
			if err != nil {
				t.Error(err)
			}
		}
	})
}

func TestConfig(t *testing.T) {
	config, err := LoadConfigFile("test.yml")
	if err != nil {
		t.Error("Failed to load test.yml compose file")
	}

	t.Run("Config", func(t *testing.T) {
		err := ValidateConfig(validations, config)
		if err != nil {
			t.Error(err)
		}
	})
}
