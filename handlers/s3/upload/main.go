package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mentorship-app-backend/components/errorpackage"
	"mentorship-app-backend/entity"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"mentorship-app-backend/handlers/s3/config"
	"mentorship-app-backend/handlers/wrapper"
)

func UploadHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	config.Init()
	s3Client := config.S3Client()
	bucketName := config.BucketName()

	var uploadReq entity.UploadRequest
	err := json.Unmarshal([]byte(request.Body), &uploadReq)
	if err != nil {
		return errorpackage.ClientError(http.StatusBadRequest, fmt.Sprintf("Invalid request payload: %v err: %v ", request.Body, err.Error()))
	}

	fileData, err := base64.StdEncoding.DecodeString(uploadReq.FileContent)
	if err != nil {
		return errorpackage.ClientError(http.StatusBadRequest, fmt.Sprintf("Invalid file data: %v", fileData))
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

	fileURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, uploadReq.Filename)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       fmt.Sprintf(`{"FileURL": "%s"}`, fileURL),
		Headers:    wrapper.SetHeadersPost(),
	}, nil
}

func main() {
	lambda.Start(wrapper.HandlerWrapper(UploadHandler, "#s3-bucket", "UploadHandler"))
}
