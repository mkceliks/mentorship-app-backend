package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"mentorship-app-backend/handlers/validator"
	"mentorship-app-backend/handlers/wrapper"
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
			Headers:    wrapper.SetHeadersPost(),
		}, err
	}

	bodyReader := multipart.NewReader(strings.NewReader(request.Body), params["boundary"])
	part, err := bodyReader.NextPart()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    wrapper.SetHeadersPost(),
		}, err
	}

	key := part.FileName()
	if err = validator.ValidateKey(key); err != nil {
		return events.APIGatewayProxyResponse{}, fmt.Errorf("failed to extract key : %w", err)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(part)
	if err != nil || buf.Len() == 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "File content is empty or unreadable"}`,
			Headers:    wrapper.SetHeadersPost(),
		}, err
	}

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(buf.Bytes()),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    wrapper.SetHeadersPost(),
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    wrapper.SetHeadersPost(),
	}, err
}

func main() {
	lambda.Start(UploadHandler)
}
