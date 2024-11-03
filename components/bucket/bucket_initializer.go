package bucket

import (
	"fmt"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/jsii-runtime-go"
	"log"
)

const (
	bucketName = "big-bucket"
)

func InitializeBucket(stack awscdk.Stack, environment string) awss3.Bucket {
	finalBucketName := fmt.Sprintf("%s-%s", bucketName, environment)

	// Log the bucket name to verify it is as expected
	log.Printf("Creating bucket with name: %s", finalBucketName)

	return awss3.NewBucket(stack, jsii.String(finalBucketName), &awss3.BucketProps{
		BucketName: jsii.String(finalBucketName),
		Versioned:  jsii.Bool(true),
	})
}
