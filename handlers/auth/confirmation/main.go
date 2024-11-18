package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mentorship-app-backend/components/errorpackage"
	"mentorship-app-backend/config"
	"mentorship-app-backend/entity"
	"mentorship-app-backend/handlers/validator"
	"mentorship-app-backend/handlers/wrapper"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

var (
	cfg         config.Config
	environment = os.Getenv("ENVIRONMENT")
)

func ConfirmHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req entity.ConfirmRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return errorpackage.ClientError(http.StatusBadRequest, "Invalid request body")
	}

	log.Printf("Received confirmation request for email: %s", req.Email)

	if err := validator.ValidateEmail(req.Email); err != nil {
		return errorpackage.ClientError(http.StatusBadRequest, "Email validation failed")
	}

	if req.Code == "" {
		return errorpackage.ClientError(http.StatusBadRequest, "Confirmation code is required")
	}

	client := config.CognitoClient()
	_, err := client.ConfirmSignUp(context.TODO(), &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         &cfg.CognitoClientID,
		Username:         &req.Email,
		ConfirmationCode: &req.Code,
	})
	if err != nil {
		if errorpackage.IsInvalidConfirmationCodeError(err) {
			return errorpackage.ClientError(http.StatusBadRequest, "Invalid confirmation code")
		}
		if errorpackage.IsExpiredConfirmationCodeError(err) {
			return errorpackage.ClientError(http.StatusBadRequest, "Confirmation code expired")
		}
		return errorpackage.ServerError(fmt.Sprintf("Failed to confirm sign-up with Cognito: %s", err.Error()))
	}

	response := map[string]string{
		"message": "Email confirmed successfully",
	}
	responseBody, err := json.Marshal(response)
	if err != nil {
		return errorpackage.ServerError("Failed to marshal confirmation response")
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    wrapper.SetHeadersPost(),
		Body:       string(responseBody),
	}, nil
}

func main() {
	log.Printf("Loading configuration for environment: %s", environment)
	var err error
	cfg, err = config.LoadConfig(environment)
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	err = config.InitAWSConfig(cfg)
	if err != nil {
		log.Fatalf("failed to initialize AWS config: %v", err)
	}

	lambda.Start(wrapper.HandlerWrapper(ConfirmHandler, "#auth-cognito", "ConfirmHandler"))
}
