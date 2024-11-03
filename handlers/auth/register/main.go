package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mentorship-app-backend/config"
	"mentorship-app-backend/entity"
	"mentorship-app-backend/handlers/errorpackage"
	"mentorship-app-backend/handlers/validator"
	"mentorship-app-backend/handlers/wrapper"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

var cfg config.Config

func RegisterHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req entity.AuthRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return errorpackage.ClientError(http.StatusBadRequest, "Invalid request body")
	}

	if err := validator.ValidateEmail(req.Email); err != nil {
		return errorpackage.ClientError(http.StatusBadRequest, "Email validation failed")
	}

	client := config.CognitoClient()
	clientID := cfg.CognitoClientID

	_, err := client.SignUp(context.TODO(), &cognitoidentityprovider.SignUpInput{
		ClientId: &clientID,
		Username: &req.Email,
		Password: &req.Password,
		UserAttributes: []types.AttributeType{
			{Name: aws.String("email"), Value: &req.Email},
		},
	})
	if err != nil {
		log.Printf("Error during SignUp: %v\n", err)
		errorMessage := fmt.Sprintf("Failed to register user: %v", err.Error())
		return errorpackage.ServerError(errorMessage)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Headers:    wrapper.SetHeadersPost(),
		Body:       `{"message":"User registered successfully"}`,
	}, nil
}

func main() {
	var err error

	environment := os.Getenv("ENVIRONMENT")
	log.Printf("Loading configuration for environment: %s", environment)
	
	cfg, err = config.LoadConfig(environment)
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	err = config.InitCognitoClient(cfg)
	if err != nil {
		log.Fatalf("failed to initialize Cognito client: %v", err)
	}

	lambda.Start(RegisterHandler)
}
