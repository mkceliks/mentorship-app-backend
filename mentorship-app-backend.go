package main

import (
	"flag"
	"fmt"
	"log"
	"mentorship-app-backend/api"
	"mentorship-app-backend/components/bucket"
	"mentorship-app-backend/components/cognito"
	"mentorship-app-backend/config"
	"mentorship-app-backend/handlers"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func getEnvironment() string {
	envPtr := flag.String("environment", "", "Specify the deployment environment (staging or production)")
	flag.Parse()

	if *envPtr != "" {
		return *envPtr
	}
	if env := os.Getenv("TARGET_ENV"); env != "" {
		return env
	}

	log.Fatal("Environment not specified. Please set --environment flag or TARGET_ENV env variable.")
	return ""
}

func stackInitializer(scope constructs.Construct, id string, props *awscdk.StackProps, environment string) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, props)

	log.Printf("Initializing stack for environment: %s", environment)

	userPoolArn := config.AppConfig.CognitoPoolArn
	clientID := config.AppConfig.CognitoClientID
	bucketName := config.AppConfig.BucketName

	s3Bucket := bucket.InitializeBucket(stack, bucketName)
	fmt.Printf("Bucket Name: %s\n", *s3Bucket.BucketName())

	lambdas := map[string]awslambda.Function{
		api.RegisterLambdaName: handlers.InitializeLambda(stack, s3Bucket, api.RegisterLambdaName, clientID, userPoolArn, environment),
		api.LoginLambdaName:    handlers.InitializeLambda(stack, s3Bucket, api.LoginLambdaName, clientID, userPoolArn, environment),
		api.UploadLambdaName:   handlers.InitializeLambda(stack, s3Bucket, api.UploadLambdaName, "", "", environment),
		api.DownloadLambdaName: handlers.InitializeLambda(stack, s3Bucket, api.DownloadLambdaName, "", "", environment),
		api.ListLambdaName:     handlers.InitializeLambda(stack, s3Bucket, api.ListLambdaName, "", "", environment),
		api.DeleteLambdaName:   handlers.InitializeLambda(stack, s3Bucket, api.DeleteLambdaName, "", "", environment),
	}

	userPool := cognito.InitializeUserPool(stack, "UserPool", userPoolArn)
	cognitoAuthorizer := cognito.InitializeCognitoAuthorizer(stack, "MentorshipCognitoAuthorizer", userPool)
	api.InitializeAPI(stack, lambdas, cognitoAuthorizer, environment)

	return stack
}

func main() {
	defer jsii.Close()

	environment := getEnvironment()
	if err := config.LoadConfig(environment); err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	app := awscdk.NewApp(nil)
	awsContext := &awscdk.Environment{
		Account: jsii.String(config.AppConfig.Account),
		Region:  jsii.String(config.AppConfig.Region),
	}

	stackInitializer(
		app,
		fmt.Sprintf("mentorship-%s", environment),
		&awscdk.StackProps{Env: awsContext},
		environment,
	)
	app.Synth(nil)
}
