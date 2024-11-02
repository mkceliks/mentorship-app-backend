package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"log"
	"mentorship-app-backend/api"
	"mentorship-app-backend/components/bucket"
	"mentorship-app-backend/components/cognito"
	"mentorship-app-backend/config"
	"mentorship-app-backend/handlers"
)

func stackInitializer(
	scope constructs.Construct,
	id string,
	props *awscdk.StackProps,
	environment string,
) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, props)

	// cognito
	userPool := cognito.InitializeUserPool(stack, "MentorshipUserPool")
	userPoolClient := cognito.InitializeUserPoolClient(userPool, "MentorshipUserPoolClient")
	cognitoAuthorizer := cognito.InitializeCognitoAuthorizer(stack, "MentorshipCognitoAuthorizer", userPool)

	// s3
	s3Bucket := bucket.InitializeBucket(stack, config.AppConfig.BucketName)

	// lambdas
	lambdas := map[string]awslambda.Function{
		"register": handlers.InitializeLambda(stack, s3Bucket, "register", *userPoolClient.UserPoolClientId()),
		"login":    handlers.InitializeLambda(stack, s3Bucket, "login", *userPoolClient.UserPoolClientId()),
		"upload":   handlers.InitializeLambda(stack, s3Bucket, "upload", ""),
		"download": handlers.InitializeLambda(stack, s3Bucket, "download", ""),
		"list":     handlers.InitializeLambda(stack, s3Bucket, "list", ""),
		"delete":   handlers.InitializeLambda(stack, s3Bucket, "delete", ""),
	}

	api.InitializeAPI(stack, lambdas, cognitoAuthorizer, environment)

	return stack
}

func main() {
	defer jsii.Close()

	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	app := awscdk.NewApp(nil)
	awsContext := &awscdk.Environment{
		Account: jsii.String(config.AppConfig.Context.Account),
		Region:  jsii.String(config.AppConfig.Context.Region),
	}

	stackInitializer(
		app,
		config.AppConfig.Environment.AppName+"-"+config.AppConfig.Environment.Staging,
		&awscdk.StackProps{Env: awsContext},
		config.AppConfig.Environment.Staging,
	)
	stackInitializer(
		app,
		config.AppConfig.Environment.AppName+"-"+config.AppConfig.Environment.Production,
		&awscdk.StackProps{Env: awsContext},
		config.AppConfig.Environment.Production,
	)

	app.Synth(nil)
}
