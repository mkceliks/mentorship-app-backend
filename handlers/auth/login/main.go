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
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

var (
	cfg         config.Config
	environment = os.Getenv("ENVIRONMENT")
)

func LoginHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req entity.AuthRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return errorpackage.ClientError(http.StatusBadRequest, "Invalid request body")
	}

	if err := validator.ValidateEmail(req.Email); err != nil {
		return errorpackage.ClientError(http.StatusBadRequest, "Email validation failed")
	}

	client := config.CognitoClient()
	resp, err := client.InitiateAuth(context.TODO(), &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		ClientId: &cfg.CognitoClientID,
		AuthParameters: map[string]string{
			"USERNAME": req.Email,
			"PASSWORD": req.Password,
		},
	})
	if err != nil {
		if errorpackage.IsInvalidCredentialsError(err) {
			return errorpackage.ClientError(http.StatusUnauthorized, "Invalid credentials")
		}
		return errorpackage.ServerError(fmt.Sprintf("Failed to authenticate with Cognito provider: %s", err.Error()))
	}

	if resp.AuthenticationResult == nil {
		return errorpackage.ServerError("Authentication failed: empty authentication result from Cognito")
	}

	tokens := map[string]string{}
	if resp.AuthenticationResult.AccessToken != nil {
		tokens["access_token"] = *resp.AuthenticationResult.AccessToken
	}
	if resp.AuthenticationResult.IdToken != nil {
		tokens["id_token"] = *resp.AuthenticationResult.IdToken
	}
	if resp.AuthenticationResult.RefreshToken != nil {
		tokens["refresh_token"] = *resp.AuthenticationResult.RefreshToken
	}

	responseBody, err := json.Marshal(tokens)
	if err != nil {
		return errorpackage.ServerError("Failed to marshal authentication tokens")
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

	lambda.Start(wrapper.HandlerWrapper(LoginHandler, "#auth-cognito", "LoginHandler"))
}
