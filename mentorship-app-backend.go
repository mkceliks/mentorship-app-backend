package main

import (
	"fmt"
	"log"
	"mentorship-app-backend/api"
	"mentorship-app-backend/components/bucket"
	"mentorship-app-backend/components/cognito"
	"mentorship-app-backend/config"
	"mentorship-app-backend/handlers"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"os"
)

func getEnvironment() string {
	env := os.Getenv("TARGET_ENV")
	if env == "" {
		log.Fatal("Environment not specified. Please set TARGET_ENV environment variable.")
	}
	return env
}

func stackInitializer(scope constructs.Construct, id string, props *awscdk.StackProps, cfg config.Config) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, props)

	log.Printf("Initializing stack for environment: %s", cfg.AppName)

	s3Bucket := bucket.InitializeBucket(stack, cfg.BucketName)
	fmt.Printf("Bucket Name: %s\n", *s3Bucket.BucketName())

	lambdas := map[string]awslambda.Function{
		api.RegisterLambdaName: handlers.InitializeLambda(stack, s3Bucket, api.RegisterLambdaName, cfg.CognitoClientID, cfg.CognitoPoolArn, cfg.AppName),
		api.LoginLambdaName:    handlers.InitializeLambda(stack, s3Bucket, api.LoginLambdaName, cfg.CognitoClientID, cfg.CognitoPoolArn, cfg.AppName),
		api.UploadLambdaName:   handlers.InitializeLambda(stack, s3Bucket, api.UploadLambdaName, "", "", cfg.AppName),
		api.DownloadLambdaName: handlers.InitializeLambda(stack, s3Bucket, api.DownloadLambdaName, "", "", cfg.AppName),
		api.ListLambdaName:     handlers.InitializeLambda(stack, s3Bucket, api.ListLambdaName, "", "", cfg.AppName),
		api.DeleteLambdaName:   handlers.InitializeLambda(stack, s3Bucket, api.DeleteLambdaName, "", "", cfg.AppName),
	}

	userPool := cognito.InitializeUserPool(stack, "UserPool", cfg.CognitoPoolArn)
	cognitoAuthorizer := cognito.InitializeCognitoAuthorizer(stack, "MentorshipCognitoAuthorizer", userPool)

	api.InitializeAPI(stack, lambdas, cognitoAuthorizer, cfg.AppName)

	return stack
}

func main() {
	defer jsii.Close()

	environment := getEnvironment()
	cfg, err := config.LoadConfig(environment)
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	if err = config.InitCognitoClient(cfg); err != nil {
		log.Fatalf("failed to initialize Cognito client: %v", err)
	}

	app := awscdk.NewApp(nil)
	awsContext := &awscdk.Environment{
		Account: jsii.String(cfg.Account),
		Region:  jsii.String(cfg.Region),
	}

	stackInitializer(app, fmt.Sprintf("mentorship-%s", environment), &awscdk.StackProps{Env: awsContext}, cfg)
	app.Synth(nil)
}
