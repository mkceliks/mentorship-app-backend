package main

import (
	"context"
	"encoding/base64"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"log"
	errorPackage "mentorship-app-backend/handlers/errorpackage"
	"mentorship-app-backend/handlers/s3/config"
	"mentorship-app-backend/handlers/validator"
	"mentorship-app-backend/handlers/wrapper"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

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
		return errorPackage.ClientError(http.StatusBadRequest, "Invalid or missing key parameter")
	}

	resp, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return errorPackage.HandleS3Error(err)
	}

	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read file content: %v", err)
		return errorPackage.ServerError("Failed to read file content")
	}

	contentType := aws.ToString(resp.ContentType)
	if contentType == "" {
		contentType = detectContentType(key)
	}

	isBinary := isBinaryType(contentType)

	responseBody := string(content)

	if isBinary {
		responseBody = base64.StdEncoding.EncodeToString(content)
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      http.StatusOK,
		Body:            responseBody,
		IsBase64Encoded: isBinary,
		Headers:         wrapper.SetHeadersGet(contentType),
	}, nil
}

func detectContentType(key string) string {
	ext := strings.ToLower(filepath.Ext(key))
	if mimeType := mime.TypeByExtension(ext); mimeType != "" {
		return mimeType
	}
	return "application/octet-stream"
}

func isBinaryType(contentType string) bool {
	return strings.HasPrefix(contentType, "image/") ||
		strings.HasPrefix(contentType, "video/") ||
		strings.HasPrefix(contentType, "application/octet-stream")
}

func main() {
	lambda.Start(DownloadHandler)
}
