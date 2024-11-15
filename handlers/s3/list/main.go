package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mentorship-app-backend/entity"
	"mentorship-app-backend/handlers/wrapper"
	"net/http"

	"mentorship-app-backend/handlers/s3/config"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func ListHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	claims, ok := request.RequestContext.Authorizer["claims"].(map[string]interface{})
	if !ok {
		log.Println("No claims found in the request context. User is not authenticated.")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       "Unauthorized: User is not logged in",
		}, nil
	}

	username := claims["username"]
	email := claims["email"]
	log.Printf("Authenticated user - Username: %v, Email: %v", username, email)

	config.Init()
	s3Client := config.S3Client()
	bucketName := config.BucketName()

	log.Printf("S3 bucket name: %s", bucketName)
	log.Printf("S3 client: %v", s3Client)

	resp, err := s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		log.Printf("Failed to list files: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    wrapper.SetAccessControl(),
			Body:       fmt.Sprintf("Error listing files: %v", err),
		}, nil
	}

	var files []entity.File
	for _, item := range resp.Contents {
		files = append(files, entity.File{
			Key:  *item.Key,
			Size: *item.Size,
		})
	}

	filesJSON, err := json.Marshal(files)
	if err != nil {
		log.Printf("Failed to marshal file list: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    wrapper.SetAccessControl(),
			Body:       fmt.Sprintf("Error marshaling file list: %v", err),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(filesJSON),
		Headers:    wrapper.SetHeadersGet(""),
	}, nil
}

func main() {
	lambda.Start(wrapper.HandlerWrapper(ListHandler, "#s3-bucket", "ListHandler"))
}
