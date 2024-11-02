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
}

var AppConfig Config

func LoadConfig() error {
	data, err := os.ReadFile("./config.yaml")
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, &AppConfig)
	if err != nil {
		return err
	}
	return nil
}
