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

const (
	EnvironmentStaging    = "staging"
	EnvironmentProduction = "production"
)

type MentorshipAppBackendStackProps struct {
	awscdk.StackProps
}

func NewMentorshipAppBackendStack(scope constructs.Construct, id string, props *MentorshipAppBackendStackProps, environment string) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)

	bucket := initializeBucket(stack, environment)
	uploadLambda := initializeLambda(stack, bucket, "upload")
	downloadLambda := initializeLambda(stack, bucket, "download")
	listLambda := initializeLambda(stack, bucket, "list")

	initializeAPI(stack, uploadLambda, downloadLambda, listLambda, environment)

	awscdk.NewCfnOutput(stack, jsii.String("BucketNameOutput"), &awscdk.CfnOutputProps{
		Value:       bucket.BucketName(),
		Description: jsii.String("The name of the S3 bucket used for the mentorship app."),
	})

	return stack
}

func initializeBucket(stack awscdk.Stack, environment string) awss3.Bucket {
	bucketName := fmt.Sprintf("mentorshipappbucket-%s-%s", environment, *stack.Account())
	return awss3.NewBucket(stack, jsii.String("MentorshipAppBucket"), &awss3.BucketProps{
		BucketName: jsii.String(bucketName),
		Versioned:  jsii.Bool(true),
	})
}

func initializeLambda(stack awscdk.Stack, bucket awss3.Bucket, functionName string) awslambda.Function {
	lambdaFunction := awslambda.NewFunction(stack, jsii.String(fmt.Sprintf("%sLambda", functionName)), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2(),
		Handler: jsii.String("bootstrap"),
		Code:    awslambda.Code_FromAsset(jsii.String(fmt.Sprintf("./output/%s_function.zip", functionName)), nil),
		Environment: &map[string]*string{
			"BUCKET_NAME": bucket.BucketName(),
		},
	})

	if functionName == "upload" {
		bucket.GrantReadWrite(lambdaFunction, "*")
	} else if functionName == "download" || functionName == "list" {
		bucket.GrantRead(lambdaFunction, "*")
	}
	return lambdaFunction
}

func initializeAPI(stack awscdk.Stack, uploadLambda, downloadLambda, listLambda awslambda.Function, environment string) {
	apiName := fmt.Sprintf("MentorshipAppAPI-%s", environment)
	stageName := environment

	api := awsapigateway.NewRestApi(stack, jsii.String(apiName), &awsapigateway.RestApiProps{
		RestApiName: jsii.String(apiName),
		Description: jsii.String(fmt.Sprintf("API Gateway for %s environment", environment)),
		DeployOptions: &awsapigateway.StageOptions{
			StageName: jsii.String(stageName),
		},
		DefaultCorsPreflightOptions: &awsapigateway.CorsOptions{
			AllowOrigins: awsapigateway.Cors_ALL_ORIGINS(),
			AllowMethods: jsii.Strings("OPTIONS", "GET", "POST"),
			AllowHeaders: jsii.Strings("Content-Type", "Authorization"),
		},
	})

	upload := api.Root().AddResource(jsii.String("upload"), nil)
	upload.AddMethod(jsii.String("POST"), awsapigateway.NewLambdaIntegration(uploadLambda, nil), nil)

	download := api.Root().AddResource(jsii.String("download"), nil)
	download.AddMethod(jsii.String("GET"), awsapigateway.NewLambdaIntegration(downloadLambda, nil), nil)

	list := api.Root().AddResource(jsii.String("list"), nil)
	list.AddMethod(jsii.String("GET"), awsapigateway.NewLambdaIntegration(listLambda, nil), nil)

	awscdk.NewCfnOutput(stack, jsii.String("ApiUrlOutput"), &awscdk.CfnOutputProps{
		Value:       api.Url(),
		Description: jsii.String(fmt.Sprintf("The endpoint URL for the %s API", environment)),
	})
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	account := app.Node().TryGetContext(jsii.String("awsAccount"))
	region := app.Node().TryGetContext(jsii.String("awsRegion"))

	if account == nil || region == nil {
		panic("AWS Account and Region must be set in the context.")
	}

	NewMentorshipAppBackendStack(app, "MentorshipAppBackendStagingStack", &MentorshipAppBackendStackProps{
		awscdk.StackProps{
			Env: &awscdk.Environment{
				Account: jsii.String(account.(string)),
				Region:  jsii.String(region.(string)),
			},
		},
	}, EnvironmentStaging)

	NewMentorshipAppBackendStack(app, "MentorshipAppBackendProductionStack", &MentorshipAppBackendStackProps{
		awscdk.StackProps{
			Env: &awscdk.Environment{
				Account: jsii.String(account.(string)),
				Region:  jsii.String(region.(string)),
			},
		},
	}, EnvironmentProduction)

	app.Synth(nil)
}
