package main

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"mentorship-app-backend/handlers/s3/config"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
)

func UploadHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	config.Init()
	s3Client := config.S3Client()
	bucketName := config.BucketName()

	contentType := request.Headers["content-type"]
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil || !strings.HasPrefix(mediaType, "multipart/") {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Content-Type must be multipart/form-data"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	bodyReader := multipart.NewReader(strings.NewReader(request.Body), params["boundary"])
	part, err := bodyReader.NextPart()
	if err != nil {
		log.Printf("Failed to parse form data: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Failed to parse form data"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	key := part.FileName()
	if key == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "File must have a name"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(part)
	if err != nil || buf.Len() == 0 {
		log.Printf("Failed to read file content or file is empty: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "File content is empty or unreadable"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(buf.Bytes()),
	})
	if err != nil {
		log.Printf("Failed to upload file to S3: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error": "Failed to upload file to S3"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       `{"message": "Successfully uploaded file", "key": "` + key + `"}`,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func main() {
	lambda.Start(UploadHandler)
}
