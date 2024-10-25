package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"net/http"
	"strings"

	"mentorship-app-backend/handlers/s3/config"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
)

func UploadHandler(_ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Ensure s3config is initialized
	config.Init()
	s3Client := config.S3Client()
	bucketName := config.BucketName()

	key := "test-file.txt"
	content := "This is the content of the file."

	_, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   strings.NewReader(content),
	})
	if err != nil {
		log.Printf("Failed to upload file: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Failed to upload file: " + err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Successfully uploaded file to S3 with key: " + key,
	}, nil
}

func main() {
	lambda.Start(UploadHandler)
}
