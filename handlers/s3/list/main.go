package main

import (
	"context"
	"encoding/json"
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

func ListHandler(_ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	config.Init()
	s3Client := config.S3Client()
	bucketName := config.BucketName()

	log.Printf("S3 bucket name: %s", bucketName)

	resp, err := s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		log.Printf("Failed to list files in bucket %s: %v", bucketName, err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    wrapper.SetAccessControl(),
			Body:       "Error listing files",
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
			Body:       "Error processing file list",
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
