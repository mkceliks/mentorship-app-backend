package bucket

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/jsii-runtime-go"
	"mentorship-app-backend/permissions"
)

func InitializeBucket(stack awscdk.Stack, bucketName string) awss3.Bucket {
	bucket := awss3.NewBucket(stack, jsii.String(bucketName), &awss3.BucketProps{
		BucketName: jsii.String(bucketName),
		Versioned:  jsii.Bool(false),
		BlockPublicAccess: awss3.NewBlockPublicAccess(&awss3.BlockPublicAccessOptions{
			BlockPublicAcls:       jsii.Bool(false),
			BlockPublicPolicy:     jsii.Bool(false),
			IgnorePublicAcls:      jsii.Bool(false),
			RestrictPublicBuckets: jsii.Bool(false),
		}),
		PublicReadAccess: jsii.Bool(true),
	})

	permissions.GrantPublicReadAccess(bucket)

	return bucket
}
