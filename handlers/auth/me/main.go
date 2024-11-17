package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
	"mentorship-app-backend/components/errorpackage"
	"mentorship-app-backend/config"
	"mentorship-app-backend/handlers/validator"
	"mentorship-app-backend/handlers/wrapper"
	"net/http"
	"os"
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

	if validator.ValidateEmail(payload.Email) != nil {
		return errorpackage.ClientError(http.StatusBadRequest, "Invalid email format")
	}

	profileType := payload.CustomRole
	if profileType == "" {
		return errorpackage.ClientError(http.StatusBadRequest, "ProfileType (custom:role) is missing in the token")
	}

	userDetails, err := fetchUserProfile(payload.Email, profileType)
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

func fetchUserProfile(email, profileType string) (map[string]string, error) {
	client := config.DynamoDBClient()

	if email == "" || profileType == "" {
		log.Println("fetchUserProfile: email or profileType is empty")
		return nil, fmt.Errorf("email or profileType is empty")
	}

	log.Printf("Fetching user profile for UserId: %s and ProfileType: %s from table: %s", email, profileType, tableName)

	result, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"UserId":      &types.AttributeValueMemberS{Value: email},
			"ProfileType": &types.AttributeValueMemberS{Value: profileType},
		},
	})
	if err != nil {
		log.Printf("DynamoDB GetItem error: %v", err)
		return nil, err
	}

	if result.Item == nil {
		log.Printf("No item found for UserId: %s and ProfileType: %s", email, profileType)
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
