package main

import (
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

func NewMentorshipAppBackendStack(scope constructs.Construct, id string, props *MentorshipAppBackendStackProps, isProduction bool) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)

	// Set bucket name depending on environment
	bucketName := "MentorshipAppBucket-Staging"
	if isProduction {
		bucketName = "MentorshipAppBucket-Production"
	}

	// Create the S3 bucket
	bucket := awss3.NewBucket(stack, jsii.String(bucketName), &awss3.BucketProps{
		Versioned: jsii.Bool(true),
	})

	// Create the Lambda function for handling file uploads
	uploadLambda := awslambda.NewFunction(stack, jsii.String("UploadLambda"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_GO_1_X(),
		Handler: jsii.String("handlers/s3/upload"),
		Code:    awslambda.Code_FromAsset(jsii.String("./handlers"), nil),
		Environment: &map[string]*string{
			"BUCKET_NAME": bucket.BucketName(),
		},
	})

	// Grant S3 access to the Lambda, with '*' allowing access to all objects in the bucket
	bucket.GrantReadWrite(uploadLambda, jsii.String("*"))

	// Create the API Gateway
	api := awsapigateway.NewRestApi(stack, jsii.String("MentorshipAppAPI"), &awsapigateway.RestApiProps{
		RestApiName: jsii.String("MentorshipAppAPI"),
		Description: jsii.String("API Gateway for handling S3 file uploads."),
	})

	// Define the /upload route
	upload := api.Root().AddResource(jsii.String("upload"), nil)
	upload.AddMethod(jsii.String("POST"), awsapigateway.NewLambdaIntegration(uploadLambda, nil), nil)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	// Create a staging environment stack
	NewMentorshipAppBackendStack(app, "MentorshipAppBackendStagingStack", &MentorshipAppBackendStackProps{
		awscdk.StackProps{
			Env: envStaging(),
		},
	}, false)

	// Create a production environment stack
	NewMentorshipAppBackendStack(app, "MentorshipAppBackendProductionStack", &MentorshipAppBackendStackProps{
		awscdk.StackProps{
			Env: envProduction(),
		},
	}, true)

	app.Synth(nil)
}

// Staging environment settings
func envStaging() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String("034362052544"),
		Region:  jsii.String("us-east-1"),
	}
}

// Production environment settings
func envProduction() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String("034362052544"),
		Region:  jsii.String("us-east-1"),
	}
}
