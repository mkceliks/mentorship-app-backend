package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var s3Client *s3.Client
var bucketName string

// initialize s3 client
func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load config, %v", err)
	}
	s3Client = s3.NewFromConfig(cfg)
	bucketName = os.Getenv("BUCKET_NAME")
}

// UploadHandler handles the S3 file upload
func UploadHandler(_ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	key := "test-file.txt"
	content := "This is the content of the file."

	_, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   strings.NewReader(content),
	})
	if err != nil {
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
