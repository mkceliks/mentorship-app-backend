package bucket

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/jsii-runtime-go"
)

func InitializeBucket(stack awscdk.Stack, bucketName string) awss3.Bucket {
	return awss3.NewBucket(stack, jsii.String(bucketName), &awss3.BucketProps{
		BucketName:       jsii.String(bucketName),
		Versioned:        jsii.Bool(true),
		PublicReadAccess: jsii.Bool(true),
	})
}
