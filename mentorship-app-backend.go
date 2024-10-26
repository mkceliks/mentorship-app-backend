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

	bucket := initializeBucket(stack, environment)
	uploadLambda := initializeUploadLambda(stack, bucket)
	downloadLambda := initializeDownloadLambda(stack, bucket)
	initializeAPI(stack, uploadLambda, downloadLambda, environment)

	return stack
}

func initializeBucket(stack awscdk.Stack, environment string) awss3.Bucket {
	bucketName := fmt.Sprintf("mentorshipappbucket-%s", environment)
	bucket := awss3.NewBucket(stack, jsii.String("MentorshipAppBucket"), &awss3.BucketProps{
		BucketName: jsii.String(bucketName),
		Versioned:  jsii.Bool(true),
	})
	return bucket
}

func initializeUploadLambda(stack awscdk.Stack, bucket awss3.Bucket) awslambda.Function {
	uploadLambda := awslambda.NewFunction(stack, jsii.String("UploadLambda"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2(),
		Handler: jsii.String("bootstrap"),
		Code:    awslambda.Code_FromAsset(jsii.String("./handlers/s3/upload/upload_function.zip"), nil),
		Environment: &map[string]*string{
			"BUCKET_NAME": bucket.BucketName(),
		},
	})

	bucket.GrantReadWrite(uploadLambda, "*")
	return uploadLambda
}

func initializeDownloadLambda(stack awscdk.Stack, bucket awss3.Bucket) awslambda.Function {
	downloadLambda := awslambda.NewFunction(stack, jsii.String("DownloadLambda"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2(),
		Handler: jsii.String("bootstrap"),
		Code:    awslambda.Code_FromAsset(jsii.String("./handlers/s3/download/download_function.zip"), nil),
		Environment: &map[string]*string{
			"BUCKET_NAME": bucket.BucketName(),
		},
	})

	bucket.GrantRead(downloadLambda, "*")
	return downloadLambda
}

func initializeAPI(stack awscdk.Stack, uploadLambda, downloadLambda awslambda.Function, environment string) {
	apiName := fmt.Sprintf("MentorshipAppAPI-%s", environment)
	api := awsapigateway.NewRestApi(stack, jsii.String(apiName), &awsapigateway.RestApiProps{
		RestApiName: jsii.String(apiName),
		Description: jsii.String(fmt.Sprintf("API Gateway for %s environment", environment)),
	})

	upload := api.Root().AddResource(jsii.String("upload"), nil)
	upload.AddMethod(jsii.String("POST"), awsapigateway.NewLambdaIntegration(uploadLambda, nil), nil)

	download := api.Root().AddResource(jsii.String("download"), nil)
	download.AddMethod(jsii.String("GET"), awsapigateway.NewLambdaIntegration(downloadLambda, nil), nil)
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	account := app.Node().TryGetContext(jsii.String("awsAccount")).(string)
	region := app.Node().TryGetContext(jsii.String("awsRegion")).(string)

	NewMentorshipAppBackendStack(app, "MentorshipAppBackendStagingStack", &MentorshipAppBackendStackProps{
		awscdk.StackProps{
			Env: &awscdk.Environment{
				Account: jsii.String(account),
				Region:  jsii.String(region),
			},
		},
	}, "staging")

	NewMentorshipAppBackendStack(app, "MentorshipAppBackendProductionStack", &MentorshipAppBackendStackProps{
		awscdk.StackProps{
			Env: &awscdk.Environment{
				Account: jsii.String(account),
				Region:  jsii.String(region),
			},
		},
	}, "production")

	app.Synth(nil)
}
