package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func UploadHandler(_ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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
