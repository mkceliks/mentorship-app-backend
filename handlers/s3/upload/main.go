package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	errorPackage "mentorship-app-backend/handlers/errorpackage"
	"mentorship-app-backend/handlers/s3/config"
	"mentorship-app-backend/handlers/validator"
	"mentorship-app-backend/handlers/wrapper"
)

func UploadHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	config.Init()
	s3Client := config.S3Client()
	bucketName := config.BucketName()

	if request.IsBase64Encoded {
		decodedBody, err := base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			return errorPackage.ServerError(fmt.Sprintf("Failed to decode base64 payload: %v", err))
		}
		request.Body = string(decodedBody)
	}

	fmt.Printf("Content-Type: %s\n", request.Headers["content-type"])
	fmt.Printf("X-File-Content-Type: %s\n", request.Headers["x-file-content-type"])

	contentType := request.Headers["content-type"]
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil || !strings.HasPrefix(mediaType, "multipart/") {
		return errorPackage.ClientError(http.StatusBadRequest, fmt.Sprintf("Invalid content-type for multipart upload: %v", err))
	}

	bodyReader := multipart.NewReader(strings.NewReader(request.Body), params["boundary"])
	part, err := bodyReader.NextPart()
	if err != nil {
		return errorPackage.ClientError(http.StatusBadRequest, fmt.Sprintf("Failed to read multipart content: %v", err))
	}

	key := part.FileName()
	if err = validator.ValidateKey(key); err != nil {
		return errorPackage.ClientError(http.StatusBadRequest, fmt.Sprintf("Invalid file key: %v", err))
	}

	buf := new(bytes.Buffer)
	if _, err = io.Copy(buf, part); err != nil {
		return errorPackage.ServerError(fmt.Sprintf("Failed to read file content: %v", err))
	}
	if buf.Len() == 0 {
		return errorPackage.ClientError(http.StatusBadRequest, "File content is empty")
	}

	fileContentType := request.Headers["x-file-content-type"]
	if fileContentType == "" {
		fileContentType = "application/octet-stream"
	}

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(fileContentType),
	})
	if err != nil {
		return errorPackage.HandleS3Error(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       fmt.Sprintf(`{"message": "File '%s' uploaded successfully"}`, key),
		Headers:    wrapper.SetHeadersPost(),
	}, nil
}

func main() {
	lambda.Start(UploadHandler)
}
