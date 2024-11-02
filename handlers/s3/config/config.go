package config

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	s3Client   *s3.Client
	bucketName string
)

func Init() {
	if s3Client != nil && bucketName != "" {
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load AWS config, %v", err)
	}

	s3Client = s3.NewFromConfig(cfg)

	bucketName = os.Getenv("BUCKET_NAME")
	if bucketName == "" {
		log.Fatal("BUCKET_NAME environment variable is not set")
	}
	log.Printf("Initialized with Bucket Name: %s", bucketName)
}

func S3Client() *s3.Client {
	if s3Client == nil {
		Init()
	}
	return s3Client
}

func BucketName() string {
	if bucketName == "" {
		Init()
	}
	return bucketName
}
