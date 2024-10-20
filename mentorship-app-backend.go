package main

import (
	"fmt"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
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

	// Define bucket and API name based on the environment (staging vs production)
	bucketName := "MentorshipAppBucket-Staging"
	apiName := "MentorshipAppAPI-Staging"
	if isProduction {
		bucketName = "MentorshipAppBucket-Production"
		apiName = "MentorshipAppAPI-Production"
	}

	// Create the S3 bucket
	bucket := awss3.NewBucket(stack, jsii.String(bucketName), &awss3.BucketProps{
		Versioned: jsii.Bool(true),
	})

	// Create the Lambda function using Go 1.x runtime
	uploadLambda := awslambda.NewFunction(stack, jsii.String("UploadLambda"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_GO_1_X(),                                               // Specify Go 1.x runtime
		Handler: jsii.String("bootstrap"),                                                 // Lambda handler is the Go binary
		Code:    awslambda.Code_FromAsset(jsii.String("./handlers/s3/function.zip"), nil), // Zip file containing Go binary
		Environment: &map[string]*string{
			"BUCKET_NAME": jsii.String(bucketName),
		},
	})

	// Grant the Lambda permissions to read/write to the S3 bucket
	bucket.GrantReadWrite(uploadLambda, "*") // Grants read/write to all objects in the bucket

	// Add S3 specific permissions to the Lambda's execution role (explicitly listing permissions)
	uploadLambda.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("s3:PutObject", "s3:GetObject"),
		Resources: jsii.Strings(fmt.Sprintf("%s/*", *bucket.BucketArn())),
	}))

	// Set up API Gateway and integrate it with Lambda
	api := awsapigateway.NewRestApi(stack, jsii.String(apiName), &awsapigateway.RestApiProps{
		RestApiName: jsii.String(apiName),
		Description: jsii.String(fmt.Sprintf("API Gateway for %s environment", apiName)),
	})

	upload := api.Root().AddResource(jsii.String("upload"), nil)
	upload.AddMethod(jsii.String("POST"), awsapigateway.NewLambdaIntegration(uploadLambda, nil), nil)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	// Create Staging Stack
	NewMentorshipAppBackendStack(app, "MentorshipAppBackendStagingStack", &MentorshipAppBackendStackProps{
		awscdk.StackProps{
			Env: envStaging(),
		},
	}, false)

	// Create Production Stack
	NewMentorshipAppBackendStack(app, "MentorshipAppBackendProductionStack", &MentorshipAppBackendStackProps{
		awscdk.StackProps{
			Env: envProduction(),
		},
	}, true)

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
