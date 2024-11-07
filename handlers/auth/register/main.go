package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mentorship-app-backend/components/notifier"
	"mentorship-app-backend/components/secrets"
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
	cfg config.Config
)

var (
	slackWebhookARN = os.Getenv("SLACK_WEBHOOK_SECRET_ARN")
	clientID        = os.Getenv("COGNITO_CLIENT_ID")
)

func RegisterHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Received registration request: %v", request)

	slackWebhookURL, err := secrets.GetSecretValue(slackWebhookARN)
	if err != nil {
		log.Printf("Failed to retrieve Slack webhook URL: %v", err)
		return errorpackage.ServerError(fmt.Sprintf("Internal server error : %v", err))
	}

	var req entity.AuthRequest
	if err = json.Unmarshal([]byte(request.Body), &req); err != nil {
		log.Printf("Invalid request body: %v", err)
		notifyErr := notifier.SendSlackNotification(slackWebhookURL, fmt.Sprintf("Invalid request body: %v", err))
		if notifyErr != nil {
			log.Printf("Failed to send Slack notification: %v", notifyErr)
		}
		return errorpackage.ClientError(http.StatusBadRequest, "Invalid request body")
	}

	log.Printf("Validated request body. Email: %s", req.Email)
	notifyErr := notifier.SendSlackNotification(slackWebhookURL, fmt.Sprintf("Received registration request for email: %s", req.Email))
	if notifyErr != nil {
		log.Printf("Failed to send Slack notification: %v", notifyErr)
	}

	if err = validator.ValidateEmail(req.Email); err != nil {
		log.Printf("Email validation failed for %s: %v", req.Email, err)
		notifyErr = notifier.SendSlackNotification(slackWebhookURL, fmt.Sprintf("Email validation failed for %s: %v", req.Email, err))
		if notifyErr != nil {
			log.Printf("Failed to send Slack notification: %v", notifyErr)
		}
		return errorpackage.ClientError(http.StatusBadRequest, "Email validation failed")
	}

	client := config.CognitoClient()

	log.Printf("Attempting Cognito SignUp for user: %s", req.Email)
	_, err = client.SignUp(context.TODO(), &cognitoidentityprovider.SignUpInput{
		ClientId: &clientID,
		Username: &req.Email,
		Password: &req.Password,
		UserAttributes: []types.AttributeType{
			{Name: aws.String("email"), Value: &req.Email},
		},
	})
	if err != nil {
		log.Printf("Error during SignUp for %s: %v", req.Email, err)
		errorMessage := fmt.Sprintf("Failed to register user: %v", err.Error())
		notifyErr = notifier.SendSlackNotification(slackWebhookURL, fmt.Sprintf("Error during SignUp for %s: %v", req.Email, err))
		if notifyErr != nil {
			log.Printf("Failed to send Slack notification: %v", notifyErr)
		}
		return errorpackage.ServerError(errorMessage)
	}

	log.Printf("User %s registered successfully", req.Email)
	notifyErr = notifier.SendSlackNotification(slackWebhookURL, fmt.Sprintf("User %s registered successfully", req.Email))
	if notifyErr != nil {
		log.Printf("Failed to send Slack notification: %v", notifyErr)
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
