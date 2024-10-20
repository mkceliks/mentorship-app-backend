package main

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type MentorshipAppBackendStackProps struct {
	awscdk.StackProps
}

func NewMentorshipAppBackendStack(scope constructs.Construct, id string, props *MentorshipAppBackendStackProps, environment string) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)

	// Set bucket name based on environment
	bucketName := fmt.Sprintf("mentorshipappbucket-%s", environment)

	// Check if the bucket exists, otherwise, use fromBucketName to reference the existing bucket
	bucket := awss3.Bucket_FromBucketName(stack, jsii.String("ExistingBucket"), jsii.String(bucketName))

	// Create Lambda function for S3 upload
	uploadLambda := awslambda.NewFunction(stack, jsii.String("UploadLambda"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2(),
		Handler: jsii.String("bootstrap"),
		Code:    awslambda.Code_FromAsset(jsii.String("./handlers/s3"), nil),
		Environment: &map[string]*string{
			"BUCKET_NAME": jsii.String(bucketName),
		},
	})

	// Grant Lambda permission to read/write from the S3 bucket
	bucket.GrantReadWrite(uploadLambda, "*")

	// Create API Gateway for environment
	apiName := fmt.Sprintf("MentorshipAppAPI-%s", environment)
	api := awsapigateway.NewRestApi(stack, jsii.String(apiName), &awsapigateway.RestApiProps{
		RestApiName: jsii.String(apiName),
		Description: jsii.String(fmt.Sprintf("API Gateway for %s environment", environment)),
	})

	upload := api.Root().AddResource(jsii.String("upload"), nil)
	upload.AddMethod(jsii.String("POST"), awsapigateway.NewLambdaIntegration(uploadLambda, nil), nil)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	// Staging environment
	NewMentorshipAppBackendStack(app, "MentorshipAppBackendStagingStack", &MentorshipAppBackendStackProps{
		awscdk.StackProps{
			Env: envStaging(),
		},
	}, "staging")

	// Production environment
	NewMentorshipAppBackendStack(app, "MentorshipAppBackendProductionStack", &MentorshipAppBackendStackProps{
		awscdk.StackProps{
			Env: envProduction(),
		},
	}, "production")

	app.Synth(nil)
}

func envStaging() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String("034362052544"),
		Region:  jsii.String("us-east-1"),
	}
}

func envProduction() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String("034362052544"),
		Region:  jsii.String("us-east-1"),
	}
}
