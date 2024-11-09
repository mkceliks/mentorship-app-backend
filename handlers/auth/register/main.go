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
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	lambdaservice "github.com/aws/aws-sdk-go-v2/service/lambda"
)

var (
	cfg          config.Config
	environment  = os.Getenv("ENVIRONMENT")
	tableName    = os.Getenv("DDB_TABLE_NAME")
	lambdaClient *lambdaservice.Client
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
	_, err := client.SignUp(context.TODO(), &cognitoidentityprovider.SignUpInput{
		ClientId: &cfg.CognitoClientID,
		Username: &req.Email,
		Password: &req.Password,
		UserAttributes: []types.AttributeType{
			{Name: aws.String("email"), Value: &req.Email},
		},
	})
	if err != nil {
		return errorpackage.ServerError(fmt.Sprintf("Failed to register user: %s", err.Error()))
	}

	uploadResponse, err := invokeUploadLambda(req.Email, req.ProfilePicture)
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

	err = saveUserProfile(req.Email, uploadResponse.FileURL)
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

func invokeUploadLambda(email, base64Image string) (*entity.UploadResponse, error) {
	uploadReq := entity.UploadRequest{
		FileContent: base64Image,
		Filename:    fmt.Sprintf("profile_pictures/%s.jpg", email),
	}
	payload, err := json.Marshal(uploadReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal upload request: %v", err)
	}

	log.Printf("Payload to UploadHandler: %s", string(payload))

	resp, err := lambdaClient.Invoke(context.TODO(), &lambdaservice.InvokeInput{
		FunctionName: aws.String("upload"),
		Payload:      payload,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to invoke upload Lambda: %v", err)
	}

	var uploadResp entity.UploadResponse
	if err := json.Unmarshal(resp.Payload, &uploadResp); err != nil {
		return nil, fmt.Errorf("failed to parse upload response: %v", err)
	}

	return &uploadResp, nil
}

func saveUserProfile(email, profilePicURL string) error {
	profile := map[string]dynamodbTypes.AttributeValue{
		"UserId":        &dynamodbTypes.AttributeValueMemberS{Value: email},
		"ProfileType":   &dynamodbTypes.AttributeValueMemberS{Value: "EndUser"},
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

	lambdaClient = lambdaservice.NewFromConfig(config.AWSConfig())
	lambda.Start(wrapper.HandlerWrapper(RegisterHandler, "#auth-cognito", "RegisterHandler"))
}

func extractUserPoolID(cognitoPoolArn string) string {
	parts := strings.Split(cognitoPoolArn, "/")
	return parts[len(parts)-1]
}
