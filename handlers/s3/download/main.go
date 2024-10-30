package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"log"
	"mentorship-app-backend/handlers/errorpackage"
	"mentorship-app-backend/handlers/validator"
	"mentorship-app-backend/handlers/wrapper"
	"net/http"
	"strings"

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

	if err := validator.ValidateKey(key); err != nil {
		return events.APIGatewayProxyResponse{}, fmt.Errorf("failed to extract key : %w", err)
	}

	resp, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		switch {
		case errors.Is(err, errorPackage.ErrNoSuchKey):
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusNotFound,
				Headers:    wrapper.SetHeadersGet(),
			}, err
		default:
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Headers:    wrapper.SetHeadersGet(),
			}, err
		}
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read file content: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    wrapper.SetHeadersGet(),
		}, err
	}

	contentType := aws.ToString(resp.ContentType)
	if contentType == "" || strings.HasSuffix(strings.ToLower(key), ".png") {
		contentType = "image/png"
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      http.StatusOK,
		Body:            string(content),
		IsBase64Encoded: false,
		Headers:         wrapper.SetHeadersGet(),
	}, nil
}

func main() {
	lambda.Start(DownloadHandler)
}
