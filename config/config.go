package config

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"gopkg.in/yaml.v3"
)

type CognitoConfig struct {
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

	log.Printf("Loaded config from %s", configPath)

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(AppConfig.Context.Region))
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	cognitoClient = cognitoidentityprovider.NewFromConfig(cfg)
	return nil
}

func CognitoClient() *cognitoidentityprovider.Client {
	return cognitoClient
}

func CognitoClientID() string {
	switch AppConfig.Environment.Staging {
	case "production":
		return AppConfig.Environment.Cognito.ProductionClientID
	case "staging":
		return AppConfig.Environment.Cognito.StagingClientID
	default:
		log.Fatal("Environment must be set to 'staging' or 'production' in config.yaml")
		return ""
	}
}
