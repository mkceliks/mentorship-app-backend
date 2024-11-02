package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
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

func RegisterHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req entity.AuthRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return errorpackage.ClientError(http.StatusBadRequest, "Invalid request body")
	}

	if err := validator.ValidateEmail(req.Email); err != nil {
		return errorpackage.ClientError(http.StatusBadRequest, "Email validation failed")
	}

	client := config.CognitoClient()
	clientID := config.CognitoClientID()

	_, err := client.SignUp(context.TODO(), &cognitoidentityprovider.SignUpInput{
		ClientId: &clientID,
		Username: &req.Email,
		Password: &req.Password,
		UserAttributes: []types.AttributeType{
			{Name: aws.String("email"), Value: &req.Email},
		},
	})
	if err != nil {
		if errorpackage.IsUserAlreadyExistsError(err) {
			return errorpackage.ClientError(http.StatusConflict, "User already exists")
		}
		return errorpackage.ServerError("Failed to register user")
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Headers:    wrapper.SetHeadersPost(),
		Body:       `{"message":"User registered successfully"}`,
	}, nil
}

func main() {
	lambda.Start(RegisterHandler)
}
