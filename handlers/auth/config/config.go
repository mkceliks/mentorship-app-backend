package config

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

var (
	cognitoClient *cognitoidentityprovider.Client
	clientID      string
	region        string
)

func Init() {
	if cognitoClient != nil && clientID != "" && region != "" {
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	cognitoClient = cognitoidentityprovider.NewFromConfig(cfg)

	clientID = os.Getenv("COGNITO_CLIENT_ID")
	if clientID == "" {
		log.Fatal("COGNITO_CLIENT_ID environment variable is not set")
	}

	region = os.Getenv("AWS_REGION")
	if region == "" {
		log.Fatal("AWS_REGION environment variable is not set")
	}

	log.Printf("Initialized Cognito with Client ID: %s and Region: %s", clientID, region)
}

func CognitoClient() *cognitoidentityprovider.Client {
	Init()
	return cognitoClient
}

func CognitoClientID() string {
	Init()
	return clientID
}
