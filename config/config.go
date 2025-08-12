package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Host struct {
	Name string `yaml:"name"`
	Addr string `yaml:"addr"`
}

type Settings struct {
	TimeoutSeconds int64 `yaml:"timeout_seconds"`
}

type Config struct {
	Hosts    []Host
	Settings Settings
}

func Parse(path string) Config {
	filename, err := filepath.Abs(path)

	if err != nil {
		fmt.Println("Error parsing filepath: ", err)
		panic(err)
	}

	data, err := os.ReadFile(filename)

	if err != nil {
		fmt.Println("Error reading file:", err)
		panic(err)
	}

	var config Config

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("There is an error with the format of the config file:", err)
		panic(err)
	}

	if !(config.Settings.TimeoutSeconds > 0) {
		fmt.Println("No timeout value was supplied, using default 1 second")
		config.Settings.TimeoutSeconds = 1
	}

	return config

}
