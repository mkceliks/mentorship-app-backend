package config

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Environment     string `yaml:"ENVIRONMENT"`
	Account         string `yaml:"ACCOUNT"`
	AppName         string `yaml:"APP_NAME"`
	Region          string `yaml:"REGION"`
	CognitoPoolArn  string `yaml:"COGNITO_POOL_ARN"`
	CognitoClientID string `yaml:"COGNITO_CLIENT_ID"`
	BucketName      string `yaml:"BUCKET_NAME"`
}

var AppConfig Config

func LoadConfig(environment string, filePath ...string) (Config, error) {
	configPath := "./config/config.yaml"
	if len(filePath) > 0 {
		configPath = filePath[0]
	}

	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file at %s: %w", configPath, err)
	}

	configData := make(map[string]Config)
	if err = yaml.Unmarshal(configFile, &configData); err != nil {
		return Config{}, fmt.Errorf("failed to parse YAML in config file: %w", err)
	}

	envConfig, exists := configData[environment]
	if !exists {
		return Config{}, fmt.Errorf("environment %s not found in config file", environment)
	}

	AppConfig = envConfig

	return envConfig, nil
}

var cognitoClient *cognitoidentityprovider.Client

func InitCognitoClient(cfg Config) error {
	awsConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(cfg.Region))
	if err != nil {
		return fmt.Errorf("failed to load AWS configuration: %w", err)
	}

	cognitoClient = cognitoidentityprovider.NewFromConfig(awsConfig)
	return nil
}

func CognitoClient() *cognitoidentityprovider.Client {
	if cognitoClient == nil {
		log.Fatal("Cognito client not initialized. Call InitCognitoClient first.")
	}
	return cognitoClient
}
