package s3

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var s3Client *s3.Client
var bucketName string

func init() {
	// Load AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
	s3Client = s3.NewFromConfig(cfg)
	bucketName = os.Getenv("BUCKET_NAME")
}

// UploadHandler handles the S3 file upload
func UploadHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	key := "test-file.txt"
	content := "This is the content of the file."

	// Upload the file content to S3
	_, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   strings.NewReader(content), // No need for ReadSeekCloser
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Failed to upload file: %v", err),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       fmt.Sprintf("Successfully uploaded file to S3 with key: %s", key),
	}, nil
}

func main() {
	lambda.Start(UploadHandler)
}
