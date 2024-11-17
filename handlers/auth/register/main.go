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
	"mentorship-app-backend/pkg"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var (
	cfg         config.Config
	environment = os.Getenv("ENVIRONMENT")
	tableName   = os.Getenv("DDB_TABLE_NAME")
	apiClient   *pkg.Client
)

func RegisterHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req entity.AuthRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return errorpackage.ClientError(http.StatusBadRequest, "Invalid request body")
	}

	log.Printf("Unmarshaled payload: %v", req)

	if err := validator.ValidateFields(req.Name, req.Email, req.Password, req.Role); err != nil {
		return errorpackage.ClientError(http.StatusBadRequest, err.Error())
	}

	client := config.CognitoClient()
	_, err := client.SignUp(context.TODO(), &cognitoidentityprovider.SignUpInput{
		ClientId: &cfg.CognitoClientID,
		Username: &req.Email,
		Password: &req.Password,
		UserAttributes: []types.AttributeType{
			{Name: aws.String("email"), Value: &req.Email},
			{Name: aws.String("name"), Value: &req.Name},
		},
	})
	if err != nil {
		return errorpackage.ServerError(fmt.Sprintf("Failed to register user: %s", err.Error()))
	}

	_, err = client.AdminUpdateUserAttributes(context.TODO(), &cognitoidentityprovider.AdminUpdateUserAttributesInput{
		UserPoolId: aws.String(extractUserPoolID(cfg.CognitoPoolArn)),
		Username:   &req.Email,
		UserAttributes: []types.AttributeType{
			{Name: aws.String("custom:role"), Value: &req.Role},
		},
	})
	if err != nil {
		return errorpackage.ServerError(fmt.Sprintf("Failed to update user role: %s", err.Error()))
	}

	uploadResponse, err := apiClient.UploadProfilePicture(req.FileName, req.ProfilePicture, request.Headers["x-file-content-type"])
	if err != nil {
		user, delErr := client.AdminDeleteUser(context.TODO(), &cognitoidentityprovider.AdminDeleteUserInput{
			UserPoolId: aws.String(extractUserPoolID(cfg.CognitoPoolArn)),
			Username:   &req.Email,
		})
		if delErr != nil {
			return events.APIGatewayProxyResponse{}, fmt.Errorf("user delete process failed : %v deleteErr: %v invokeErr: %v", user, delErr, err)
		}
		return errorpackage.ServerError(fmt.Sprintf("Failed to upload profile picture: %s", err.Error()))
	}

	err = saveUserProfile(req.Email, req.Name, req.Role, uploadResponse.FileURL)
	if err != nil {
		user, delErr := client.AdminDeleteUser(context.TODO(), &cognitoidentityprovider.AdminDeleteUserInput{
			UserPoolId: aws.String(extractUserPoolID(cfg.CognitoPoolArn)),
			Username:   &req.Email,
		})
		if delErr != nil {
			return events.APIGatewayProxyResponse{}, fmt.Errorf("user delete process failed : %v deleteErr: %v invokeErr: %v", user, delErr, err)
		}
		return errorpackage.ServerError(fmt.Sprintf("Failed to save user profile: %s", err.Error()))
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Headers:    wrapper.SetHeadersPost(),
		Body:       `{"message":"User registered and profile created successfully"}`,
	}, nil
}

func saveUserProfile(email, name, role, profilePicURL string) error {
	profile := map[string]dynamodbTypes.AttributeValue{
		"UserId":        &dynamodbTypes.AttributeValueMemberS{Value: email},
		"Name":          &dynamodbTypes.AttributeValueMemberS{Value: name},
		"ProfileType":   &dynamodbTypes.AttributeValueMemberS{Value: role},
		"Email":         &dynamodbTypes.AttributeValueMemberS{Value: email},
		"ProfilePicURL": &dynamodbTypes.AttributeValueMemberS{Value: profilePicURL},
	}

	_, err := config.DynamoDBClient().PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName:           aws.String(tableName),
		Item:                profile,
		ConditionExpression: aws.String("attribute_not_exists(UserId)"),
	})
	return err
}

func extractUserPoolID(cognitoPoolArn string) string {
	parts := strings.Split(cognitoPoolArn, "/")
	return parts[len(parts)-1]
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

	apiClient = pkg.NewClient()

	lambda.Start(wrapper.HandlerWrapper(RegisterHandler, "#auth-cognito", "RegisterHandler"))
}
