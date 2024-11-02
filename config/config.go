package config

import (
	"gopkg.in/yaml.v3"
	"log"
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
}

var AppConfig Config

func LoadConfig() error {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/config.yaml"
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &AppConfig)
	if err != nil {
		return err
	}

	log.Printf("Loaded config from %s", configPath)
	return nil
}
