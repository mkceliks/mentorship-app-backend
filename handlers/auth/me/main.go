package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mentorship-app-backend/components/errorpackage"
	"mentorship-app-backend/config"
	"mentorship-app-backend/handlers/validator"
	"mentorship-app-backend/handlers/wrapper"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var (
	cfg         config.Config
	environment = os.Getenv("ENVIRONMENT")
	tableName   = os.Getenv("DDB_TABLE_NAME")
)

func MeHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	idToken, err := validator.ValidateAuthorizationHeader(request.Headers["Authorization"])
	if err != nil {
		return errorpackage.ClientError(http.StatusUnauthorized, err.Error())
	}

	payload, err := validator.DecodeAndValidateIDToken(idToken)
	if err != nil {
		return errorpackage.ClientError(http.StatusUnauthorized, err.Error())
	}

	if payload.Email == "" {
		return errorpackage.ClientError(http.StatusUnauthorized, "Email is missing in the token")
	}

	if validator.ValidateEmail(payload.Email) != nil {
		return errorpackage.ClientError(http.StatusBadRequest, "Invalid email format")
	}

	userDetails, err := fetchUserProfile(payload.Email)
	if err != nil {
		if errorpackage.IsDynamoDBNotFoundError(err) {
			return errorpackage.ClientError(http.StatusNotFound, "User profile not found")
		}
		return errorpackage.ServerError(err.Error())
	}

	responseBody, err := json.Marshal(userDetails)
	if err != nil {
		return errorpackage.ServerError("Failed to marshal user details")
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    wrapper.SetHeadersGet(""),
		Body:       string(responseBody),
	}, nil
}

func fetchUserProfile(email string) (map[string]string, error) {
	client := config.DynamoDBClient()

	if email == "" {
		log.Println("fetchUserProfile: email is empty")
		return nil, fmt.Errorf("email is empty")
	}

	email = strings.TrimSpace(email)

	log.Printf("Fetching user profile for UserId: %s from table: %s", email, tableName)

	result, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"UserId": &types.AttributeValueMemberS{Value: email},
		},
	})
	if err != nil {
		log.Printf("DynamoDB GetItem error: %v", err)
		return nil, err
	}

	if result.Item == nil {
		log.Printf("No item found for UserId: %s", email)
		return nil, errorpackage.ErrNoSuchKey
	}

	userDetails := map[string]string{}
	for key, value := range result.Item {
		switch v := value.(type) {
		case *types.AttributeValueMemberS:
			userDetails[key] = v.Value
		}
	}

	log.Printf("Fetched user details: %+v", userDetails)
	return userDetails, nil
}

func main() {
	var err error
	cfg, err = config.LoadConfig(environment)
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	err = config.InitAWSConfig(cfg)
	if err != nil {
		log.Fatalf("failed to initialize AWS config: %v", err)
	}

	lambda.Start(wrapper.HandlerWrapper(MeHandler, "#auth-cognito", "MeHandler"))
}
