package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment           string `yaml:"environment"`
	Account               string `yaml:"account"`
	AppName               string `yaml:"app_name"`
	Region                string `yaml:"region"`
	CognitoPoolArn        string `yaml:"cognito_pool_arn"`
	CognitoClientID       string `yaml:"cognito_client_id"`
	BucketName            string `yaml:"bucket_name"`
	SlackWebhookSecretARN string `yaml:"slack_webhook_secret_arn"`
	EndpointBaseURL       string `yaml:"endpoint_base_url"`
}

var (
	AppConfig     Config
	awsConfig     aws.Config
	cognitoClient *cognitoidentityprovider.Client
	dynamoClient  *dynamodb.Client
)

func LoadConfig(environment string, filePath ...string) (Config, error) {
	if environment == "" {
		return Config{}, fmt.Errorf("environment variable is empty. Please specify the environment (e.g., 'staging' or 'production')")
	}

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
		return Config{}, fmt.Errorf("environment '%s' not found in config file", environment)
	}

	AppConfig = envConfig
	return envConfig, nil
}

func InitAWSConfig(cfg Config) error {
	var err error
	awsConfig, err = config.LoadDefaultConfig(context.TODO(), config.WithRegion(cfg.Region))
	if err != nil {
		return fmt.Errorf("failed to load AWS configuration: %w", err)
	}

	cognitoClient = cognitoidentityprovider.NewFromConfig(awsConfig)
	dynamoClient = dynamodb.NewFromConfig(awsConfig)

	return nil
}

func AWSConfig() aws.Config {
	return awsConfig
}

func CognitoClient() *cognitoidentityprovider.Client {
	if cognitoClient == nil {
		log.Fatal("Cognito client not initialized. Call InitAWSConfig first.")
	}
	return cognitoClient
}

func DynamoDBClient() *dynamodb.Client {
	if dynamoClient == nil {
		log.Fatal("DynamoDB client not initialized. Call InitAWSConfig first.")
	}
	return dynamoClient
}
