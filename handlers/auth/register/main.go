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

var (
	cfg         config.Config
	clientID    = os.Getenv("COGNITO_CLIENT_ID")
	environment = os.Getenv("ENVIRONMENT")
)

func RegisterHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req entity.AuthRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		log.Printf("Invalid request body: %v", err)
		return errorpackage.ClientError(http.StatusBadRequest, "Invalid request body")
	}

	if err := validator.ValidateEmail(req.Email); err != nil {
		log.Printf("Email validation failed for %s: %v", req.Email, err)
		return errorpackage.ClientError(http.StatusBadRequest, "Email validation failed")
	}

	client := config.CognitoClient()
	_, err := client.SignUp(context.TODO(), &cognitoidentityprovider.SignUpInput{
		ClientId: &clientID,
		Username: &req.Email,
		Password: &req.Password,
		UserAttributes: []types.AttributeType{
			{Name: aws.String("email"), Value: &req.Email},
		},
	})
	if err != nil {
		log.Printf("Error during Signup for %s: %v", req.Email, err)
		errorMessage := fmt.Sprintf("Failed to register user: %v", err.Error())
		return errorpackage.ServerError(errorMessage)
	}

	log.Printf("User %s registered successfully", req.Email)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Headers:    wrapper.SetHeadersPost(),
		Body:       `{"message":"User registered successfully"}`,
	}, nil
}

func main() {
	log.Printf("Loading configuration for environment: %s", environment)

	var err error
	cfg, err = config.LoadConfig(environment)
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	err = config.InitCognitoClient(cfg)
	if err != nil {
		log.Fatalf("failed to initialize Cognito client: %v", err)
	}

	lambda.Start(wrapper.HandlerWrapper(RegisterHandler, "#auth-cognito", "RegisterHandler"))
}
