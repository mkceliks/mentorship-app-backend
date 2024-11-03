package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/joho/godotenv"
)

var cognitoClient *cognitoidentityprovider.Client

func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, continuing with system environment variables")
	} else {
		log.Println(".env file loaded successfully")
	}

	requiredEnvVars := []string{
		"ACCOUNT", "REGION", "STAGING_POOL_ARN", "PRODUCTION_POOL_ARN",
		"STAGING_CLIENT_ID", "PRODUCTION_CLIENT_ID", "BUCKET_NAME",
	}

	for _, envVar := range requiredEnvVars {
		if value := os.Getenv(envVar); value == "" {
			return fmt.Errorf("missing required environment variable: %s", envVar)
		}
	}

	log.Println("All required environment variables are set.")
	return nil
}

func InitCognitoClient() error {
	if err := LoadEnv(); err != nil {
		return err
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("REGION")))
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %v", err)
	}
	cognitoClient = cognitoidentityprovider.NewFromConfig(cfg)
	return nil
}

func CognitoClient() *cognitoidentityprovider.Client {
	return cognitoClient
}

func GetCognitoClientID(environment string) (string, error) {
	clientID := os.Getenv(fmt.Sprintf("%s_CLIENT_ID", strings.ToUpper(environment)))
	if clientID == "" {
		return "", fmt.Errorf("Cognito Client ID is missing for environment: %s", environment)
	}
	log.Printf("Fetched Cognito Client ID for %s: %s", environment, clientID)
	return clientID, nil
}
