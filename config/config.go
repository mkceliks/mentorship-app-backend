package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Environment struct {
		Staging    string `yaml:"staging"`
		Production string `yaml:"production"`
		AppName    string `yaml:"app_name"`
	} `yaml:"environment"`
	Context struct {
		Account string `yaml:"account"`
		Region  string `yaml:"region"`
	} `yaml:"context"`
	BucketName string `yaml:"bucket_name"`
	RouteName  string `yaml:"route_name"`
	Lambda     struct {
		Handler string `yaml:"handler"`
	}
}

func LoadConfig() (*Config, error) {
	file, err := os.Open("config/config.yaml")
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
