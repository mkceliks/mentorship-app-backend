package main

import (
	"context"
	"errors"
	"fmt"
	"mentorship-app-backend/components/errorpackage"
	"mentorship-app-backend/handlers/s3/config"
	"mentorship-app-backend/handlers/validator"
	"mentorship-app-backend/handlers/wrapper"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func DeleteHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	config.Init()
	s3Client := config.S3Client()
	bucketName := config.BucketName()

	key := request.QueryStringParameters["key"]

	if err := validator.ValidateKey(key); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    wrapper.SetHeadersDelete(),
		}, fmt.Errorf("failed to extract key: %w", err)
	}

	_, err := s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	switch {
	case err == nil:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Headers:    wrapper.SetHeadersDelete(),
		}, nil

	case errors.Is(err, errorpackage.ErrNoSuchKey):
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Headers:    wrapper.SetHeadersDelete(),
		}, err
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    wrapper.SetHeadersDelete(),
		}, err
	}
}

func main() {
	lambda.Start(wrapper.HandlerWrapper(DeleteHandler, "#s3-bucket", "DeleteHandler"))
}
