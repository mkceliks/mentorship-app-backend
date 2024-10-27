package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"mentorship-app-backend/handlers/s3/config"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type FileInfo struct {
	Key          string    `json:"key"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"lastModified"`
	StorageClass string    `json:"storageClass"`
}

func ListHandler(_ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	config.Init()
	s3Client := config.S3Client()
	bucketName := config.BucketName()

	var files []FileInfo
	var continuationToken *string = nil

	for {
		resp, err := s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
			Bucket:            aws.String(bucketName),
			ContinuationToken: continuationToken,
		})
		if err != nil {
			log.Printf("Error listing files: %v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       "Error retrieving files: " + err.Error(),
			}, nil
		}

		for _, item := range resp.Contents {
			files = append(files, FileInfo{
				Key:          aws.ToString(item.Key),
				Size:         *item.Size,
				LastModified: aws.ToTime(item.LastModified),
				StorageClass: string(item.StorageClass),
			})
		}

		if *resp.IsTruncated {
			continuationToken = resp.NextContinuationToken
		} else {
			break
		}
	}

	filesJSON, err := json.Marshal(files)
	if err != nil {
		log.Printf("Error marshalling file list: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error processing file list",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(filesJSON),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(ListHandler)
}
