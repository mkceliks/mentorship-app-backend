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
			Body:       "Content-Type must be multipart/form-data",
		}, nil
	}

	bodyReader := multipart.NewReader(strings.NewReader(request.Body), params["boundary"])
	part, err := bodyReader.NextPart()
	if err != nil {
		log.Printf("Failed to parse form data: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Failed to parse form data: " + err.Error(),
		}, nil
	}

	key := part.FileName()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(part)
	if err != nil {
		log.Printf("Failed to read file content: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Failed to read file content: " + err.Error(),
		}, nil
	}

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(buf.Bytes()),
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
