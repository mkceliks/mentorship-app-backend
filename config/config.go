package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
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

	if data, err := os.ReadFile(configPath); err == nil {
		if err := yaml.Unmarshal(data, &AppConfig); err != nil {
			return fmt.Errorf("failed to unmarshal config.yaml: %v", err)
		}
		log.Printf("Loaded config from config.yaml: %+v", AppConfig)
	} else {
		log.Println("Config file not found, loading from environment variables")
		LoadConfigFromEnv()
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(AppConfig.Context.Region))
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %v", err)
	}
	cognitoClient = cognitoidentityprovider.NewFromConfig(cfg)

	return nil
}

func LoadConfigFromEnv() {
	AppConfig = Config{
		Environment: EnvironmentConfig{
			Staging:    os.Getenv("STAGING_ENVIRONMENT"),
			Production: os.Getenv("PRODUCTION_ENVIRONMENT"),
			AppName:    os.Getenv("APP_NAME"),
			Cognito: CognitoConfig{
				StagingPoolArn:     os.Getenv("STAGING_POOL_ARN"),
				ProductionPoolArn:  os.Getenv("PRODUCTION_POOL_ARN"),
				StagingClientID:    os.Getenv("STAGING_CLIENT_ID"),
				ProductionClientID: os.Getenv("PRODUCTION_CLIENT_ID"),
			},
		},
		Context: ContextConfig{
			Account: os.Getenv("ACCOUNT"),
			Region:  os.Getenv("REGION"),
		},
		BucketName: os.Getenv("BUCKET_NAME"),
		RouteName:  os.Getenv("ROUTE_NAME"),
	}
}

func CognitoClient() *cognitoidentityprovider.Client {
	return cognitoClient
}

func GetCognitoClientID(environment string) (string, error) {
	var clientID string
	switch environment {
	case AppConfig.Environment.Staging:
		clientID = AppConfig.Environment.Cognito.StagingClientID
	case AppConfig.Environment.Production:
		clientID = AppConfig.Environment.Cognito.ProductionClientID
	default:
		return "", fmt.Errorf("unknown environment: %s", environment)
	}

	if clientID == "" {
		log.Printf("Warning: Client ID for environment '%s' is empty", environment)
		return "", fmt.Errorf("client ID for environment '%s' is not set", environment)
	}

	log.Printf("Using Cognito Client ID: %s", clientID)
	return clientID, nil
}
