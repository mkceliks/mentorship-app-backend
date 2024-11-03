package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

var cognitoClient *cognitoidentityprovider.Client

func InitCognitoClient() error {
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
