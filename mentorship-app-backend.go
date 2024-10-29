package main

import (
	"fmt"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"log"
	"mentorship-app-backend/config"
)

const (
	bucketName         = "big-bucket"
	apiRoutes          = "api-routes"
	uploadLambdaName   = "upload"
	downloadLambdaName = "download"
	listLambdaName     = "list"
	deleteLambdaName   = "delete"
)

// TODO: refactor stack init ( errors and logs )
// TODO: set alarms and router for API gateway
func NewMentorshipAppBackendStack(
	scope constructs.Construct,
	id string,
	props *awscdk.StackProps,
	environment string,
) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, props)
	bucket := initializeBucket(stack, environment)

	initializeAPI(
		stack,
		initializeLambda(stack, bucket, uploadLambdaName),
		initializeLambda(stack, bucket, downloadLambdaName),
		initializeLambda(stack, bucket, listLambdaName),
		initializeLambda(stack, bucket, deleteLambdaName),
		environment,
	)

	return stack
}

func initializeBucket(stack awscdk.Stack, environment string) awss3.Bucket {
	return awss3.NewBucket(stack, jsii.String(bucketName), &awss3.BucketProps{
		BucketName: jsii.String(fmt.Sprintf(bucketName+"%s", environment)),
		Versioned:  jsii.Bool(true),
	})
}

func initializeLambda(stack awscdk.Stack, bucket awss3.Bucket, functionName string) awslambda.Function {
	lambdaFunction := awslambda.NewFunction(stack, jsii.String(functionName), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2(),
		Handler: jsii.String("bootstrap"),
		Code:    awslambda.Code_FromAsset(jsii.String(fmt.Sprintf("./output/%s_function.zip", functionName)), nil),
		Environment: &map[string]*string{
			"BUCKET_NAME": bucket.BucketName(),
		},
	})

	grantAccessForBucket(lambdaFunction, bucket, functionName)

	return lambdaFunction
}

func initializeAPI(stack awscdk.Stack, uploadLambda, downloadLambda, listLambda, deleteLambda awslambda.Function, environment string) {
	apiName := fmt.Sprintf(apiRoutes+"%s", environment)

	// define api gateway
	api := awsapigateway.NewRestApi(stack, jsii.String(apiName), &awsapigateway.RestApiProps{
		RestApiName: jsii.String(apiName),
		Description: jsii.String(fmt.Sprintf("API Gateway for %s environment", environment)),
		DeployOptions: &awsapigateway.StageOptions{
			StageName: jsii.String(environment),
		},
		DefaultCorsPreflightOptions: &awsapigateway.CorsOptions{
			AllowOrigins: awsapigateway.Cors_ALL_ORIGINS(),
			AllowMethods: jsii.Strings("OPTIONS", "GET", "POST", "DELETE"),
			AllowHeaders: jsii.Strings("Content-Type", "Authorization"),
		},
	})

	// create routes
	addApiResource(api, "POST", uploadLambdaName, uploadLambda)
	addApiResource(api, "GET", downloadLambdaName, downloadLambda)
	addApiResource(api, "GET", listLambdaName, listLambda)
	addApiResource(api, "DELETE", deleteLambdaName, deleteLambda)
}

func main() {
	defer jsii.Close()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	app := awscdk.NewApp(nil)

	awsContext := getAWSEnv(*cfg)

	if awsContext.Account == nil || awsContext.Region == nil {
		panic("aws account and region must be set in the context.")
	}

	NewMentorshipAppBackendStack(app,
		fmt.Sprintf(cfg.Environment.AppName+cfg.Environment.Staging),
		&awscdk.StackProps{
			Env: awsContext,
		}, cfg.Environment.Staging,
	)

	NewMentorshipAppBackendStack(app,
		fmt.Sprintf(cfg.Environment.AppName+cfg.Environment.Production),
		&awscdk.StackProps{
			Env: awsContext,
		}, cfg.Environment.Production)

	app.Synth(nil)
}

func getAWSEnv(cfg config.Config) *awscdk.Environment {
	return &awscdk.Environment{
		Account: &cfg.Context.Account,
		Region:  &cfg.Context.Region,
	}
}

// grant handler for s3bucket if needed
func grantAccessForBucket(
	lambda awslambda.Function,
	bucket awss3.Bucket,
	functionName string,
) {
	switch functionName {
	case uploadLambdaName, deleteLambdaName:
		bucket.GrantReadWrite(lambda, "*")
	case downloadLambdaName, listLambdaName:
		bucket.GrantRead(lambda, "*")
	}
}

func addApiResource(api awsapigateway.RestApi, method, resourceName string, lambdaFunction awslambda.Function) {
	resource := api.Root().AddResource(jsii.String(resourceName), nil)
	resource.AddMethod(jsii.String(method), awsapigateway.NewLambdaIntegration(lambdaFunction, nil), nil)
}
