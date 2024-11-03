package config

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type CognitoConfig struct {
	StagingPoolArn     string `yaml:"staging_pool_arn"`
	ProductionPoolArn  string `yaml:"production_pool_arn"`
	StagingClientID    string `yaml:"staging_client_id"`
	ProductionClientID string `yaml:"production_client_id"`
}

type EnvironmentConfig struct {
	Staging    string        `yaml:"staging"`
	Production string        `yaml:"production"`
	AppName    string        `yaml:"app_name"`
	Cognito    CognitoConfig `yaml:"cognito"`
}

type ContextConfig struct {
	Account string `yaml:"account"`
	Region  string `yaml:"region"`
}

type Config struct {
	Environment EnvironmentConfig `yaml:"environment"`
	Context     ContextConfig     `yaml:"context"`
	BucketName  string            `yaml:"bucket_name"`
	RouteName   string            `yaml:"route_name"`
}

var (
	AppConfig     Config
	cognitoClient *cognitoidentityprovider.Client
)

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

	log.Printf("Loaded config: %+v", AppConfig)
	return nil
}

func CognitoClient() *cognitoidentityprovider.Client {
	return cognitoClient
}

func GetCognitoClientID(environment string) (string, error) {
	switch environment {
	case AppConfig.Environment.Staging:
		return AppConfig.Environment.Cognito.StagingClientID, nil
	case AppConfig.Environment.Production:
		return AppConfig.Environment.Cognito.ProductionClientID, nil
	default:
		return "", fmt.Errorf("unknown environment: %s", environment)
	}
}
