package main

import (
	"context"
	"encoding/base64"
	"io"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"mentorship-app-backend/handlers/errorpackage"
	"mentorship-app-backend/handlers/s3/config"
	"mentorship-app-backend/handlers/wrapper"
)

func DownloadHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	config.Init()
	s3Client := config.S3Client()
	bucketName := config.BucketName()

	fileName := request.QueryStringParameters["filename"]
	if fileName == "" {
		return errorpackage.ClientError(http.StatusBadRequest, "Invalid or missing key parameter")
	}

	output, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return errorpackage.HandleS3Error(err)
	}
	defer output.Body.Close()

	fileContent, err := io.ReadAll(output.Body)
	if err != nil {
		log.Printf("Failed to read file content: %v", err)
		return errorpackage.ServerError("Failed to read file content")
	}

	base64File := base64.StdEncoding.EncodeToString(fileContent)

	return events.APIGatewayProxyResponse{
		StatusCode:      http.StatusOK,
		Body:            base64File,
		IsBase64Encoded: true,
		Headers:         wrapper.SetHeadersGet(aws.ToString(output.ContentType)),
	}, nil
}

func main() {
	lambda.Start(DownloadHandler)
}
