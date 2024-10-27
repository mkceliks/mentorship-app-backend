package main

import (
	"context"
	"encoding/base64"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"log"
	"net/http"

	"mentorship-app-backend/handlers/s3/config"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
)

func DownloadHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	config.Init()
	s3Client := config.S3Client()
	bucketName := config.BucketName()

	key := request.QueryStringParameters["key"]
	if key == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Missing 'key' query parameter",
		}, nil
	}

	resp, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Printf("Failed to download file: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Failed to download file: " + err.Error(),
		}, nil
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read file content: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Failed to read file content: " + err.Error(),
		}, nil
	}

	encodedContent := base64.StdEncoding.EncodeToString(content)
	return events.APIGatewayProxyResponse{
		StatusCode:      http.StatusOK,
		Body:            encodedContent,
		IsBase64Encoded: true,
		Headers: map[string]string{
			"Content-Type":        aws.ToString(resp.ContentType),
			"Content-Disposition": "attachment; filename=\"" + key + "\"",
		},
	}, nil
}

func main() {
	lambda.Start(DownloadHandler)
}
