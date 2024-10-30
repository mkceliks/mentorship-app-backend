package bucket

import (
	"fmt"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/jsii-runtime-go"
)

const (
	bucketName = "big-bucket"
)

func InitializeBucket(stack awscdk.Stack, environment string) awss3.Bucket {
	return awss3.NewBucket(stack, jsii.String(bucketName), &awss3.BucketProps{
		BucketName: jsii.String(fmt.Sprintf(bucketName+"%s", environment)),
		Versioned:  jsii.Bool(true),
	})
}
