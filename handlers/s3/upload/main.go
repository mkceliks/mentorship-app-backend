package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"mentorship-app-backend/handlers/errorpackage"
	"mentorship-app-backend/handlers/s3/config"
	"mentorship-app-backend/handlers/wrapper"
)

type UploadRequest struct {
	Filename    string `json:"filename"`
	FileContent string `json:"fileContent"`
}

func UploadHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	config.Init()
	s3Client := config.S3Client()
	bucketName := config.BucketName()

	var uploadReq UploadRequest
	err := json.Unmarshal([]byte(request.Body), &uploadReq)
	if err != nil {
		return errorpackage.ClientError(http.StatusBadRequest, "Invalid request payload")
	}

	fileData, err := base64.StdEncoding.DecodeString(uploadReq.FileContent)
	if err != nil {
		return errorpackage.ClientError(http.StatusBadRequest, "Invalid file data")
	}

	contentType := request.Headers["x-file-content-type"]
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(uploadReq.Filename),
		Body:        bytes.NewReader(fileData),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return errorpackage.HandleS3Error(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       fmt.Sprintf(`{"message": "File '%s' uploaded successfully"}`, uploadReq.Filename),
		Headers:    wrapper.SetHeadersPost(),
	}, nil
}

func main() {
	lambda.Start(UploadHandler)
}
