package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/docker/docker/cli/compose/loader"
	"github.com/docker/docker/cli/compose/types"
	"github.com/spf13/viper"
)

func main() {
	if len(os.Args) == 3 {
		var config *types.Config
		config, err := LoadConfigFile(os.Args[2])
		if err != nil {
			log.Fatal(err.Error())
		}

		// Get config
		viper.SetConfigName(os.Args[1]) // name of config file (without extension)
		viper.AddConfigPath("/")
		viper.AddConfigPath("$HOME")
		viper.AddConfigPath(".")
		err = viper.ReadInConfig() // yaml, toml, json, ini, whatever
		if err != nil {
			log.Fatal("Error loading the config file")
		}

		validations := Validations{
			StoragePath:       viper.GetString("StoragePath"),
			QualifiedService:  viper.GetString("QualifiedService"),
			DockerWrite:       viper.GetStringSlice("DockerWrite"),
			Secrets:           viper.GetStringSlice("Secrets"),
			Networks:          viper.GetStringSlice("Networks"),
			MemoryLimit:       viper.GetString("MemoryLimit"),
			CPULimit:          viper.GetString("CPULimit"),
			MemoryReservation: viper.GetString("MemoryReservation"),
			CPUReservation:    viper.GetString("CPUReservation"),
		}

		err = ValidateConfig(validations, config)
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Print("Success!")
	} else {
		log.Fatal("ComposeEnforcer <validation.yml> <docker-compose.yml>")
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
