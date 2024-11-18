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

func ResendHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req entity.ResendRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return errorpackage.ClientError(http.StatusBadRequest, "Invalid request body")
	}

	log.Printf("Received resend request for email: %s", req.Email)

	if err := validator.ValidateEmail(req.Email); err != nil {
		return errorpackage.ClientError(http.StatusBadRequest, "Email validation failed")
	}

	client := config.CognitoClient()
	_, err := client.ResendConfirmationCode(context.TODO(), &cognitoidentityprovider.ResendConfirmationCodeInput{
		ClientId: &cfg.CognitoClientID,
		Username: &req.Email,
	})
	if err != nil {
		return errorpackage.ServerError(fmt.Sprintf("Failed to resend confirmation code: %s", err.Error()))
	}

	response := map[string]string{
		"message": "Confirmation code resent successfully",
	}
	responseBody, err := json.Marshal(response)
	if err != nil {
		return errorpackage.ServerError("Failed to marshal resend response")
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

	lambda.Start(wrapper.HandlerWrapper(ResendHandler, "#auth-cognito", "ResendHandler"))
}
