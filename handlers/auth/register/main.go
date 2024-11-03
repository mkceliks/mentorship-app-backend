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

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func RegisterHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req entity.AuthRequest

	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return errorpackage.ClientError(http.StatusBadRequest, "Invalid request body")
	}

	if err := validator.ValidateEmail(req.Email); err != nil {
		return errorpackage.ClientError(http.StatusBadRequest, "Email validation failed")
	}

	client := config.CognitoClient()
	environment := config.AppConfig.Environment.Staging
	clientID, err := config.GetCognitoClientID(environment)
	if err != nil {
		log.Printf("Error getting Cognito Client ID: %v", err)
		return errorpackage.ServerError("Internal server error")
	}

	_, err = client.SignUp(context.TODO(), &cognitoidentityprovider.SignUpInput{
		ClientId: &clientID,
		Username: &req.Email,
		Password: &req.Password,
		UserAttributes: []types.AttributeType{
			{Name: aws.String("email"), Value: &req.Email},
		},
	})
	if err != nil {
		log.Printf("Error during SignUp: %v\n", err)
		return errorpackage.ServerError(fmt.Sprintf("Failed to register user: %v", err.Error()))
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Headers:    wrapper.SetHeadersPost(),
		Body:       `{"message":"User registered successfully"}`,
	}, nil
}

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	lambda.Start(RegisterHandler)
}
