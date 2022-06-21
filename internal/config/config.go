package config

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`

	Storage struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Database string `yaml:"database"`
	} `yaml:"storage"`
}

func NewConfig() (*Config, error) {
	config := &Config{}

	configPath, exists := os.LookupEnv("CONFIG_PATH")

	if !exists {
		return nil, errors.New("specify the path to the configuration file in .env file")
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
