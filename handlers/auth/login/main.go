package main

import (
	"context"
	"encoding/json"
	"mentorship-app-backend/entity"
	"mentorship-app-backend/handlers/auth/config"
	"mentorship-app-backend/handlers/errorpackage"
	"mentorship-app-backend/handlers/validator"
	"mentorship-app-backend/handlers/wrapper"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
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
	clientID := config.CognitoClientID()

	resp, err := client.InitiateAuth(context.TODO(), &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		ClientId: &clientID,
		AuthParameters: map[string]string{
			"USERNAME": req.Email,
			"PASSWORD": req.Password,
		},
	})
	if err != nil {
		if errorpackage.IsInvalidCredentialsError(err) {
			return errorpackage.ClientError(http.StatusUnauthorized, "Invalid credentials")
		}
		return errorpackage.ServerError("Failed to authenticate user")
	}

	tokens := map[string]string{
		"access_token":  *resp.AuthenticationResult.AccessToken,
		"id_token":      *resp.AuthenticationResult.IdToken,
		"refresh_token": *resp.AuthenticationResult.RefreshToken,
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
	lambda.Start(LoginHandler)
}
