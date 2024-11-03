package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
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

	log.Printf("Initializing stack for environment: %s", environment)

	userPoolArn := os.Getenv(fmt.Sprintf("%s_POOL_ARN", strings.ToUpper(environment)))
	clientID := os.Getenv(fmt.Sprintf("%s_CLIENT_ID", strings.ToUpper(environment)))

	if userPoolArn == "" || clientID == "" {
		log.Fatalf("Missing required environment variables for %s: %s_POOL_ARN or %s_CLIENT_ID", environment, environment, environment)
	}

	s3Bucket := bucket.InitializeBucket(stack, environment)

	fmt.Printf("Bucket Name: %s\n", *s3Bucket.BucketName())

	lambdas := map[string]awslambda.Function{
		api.RegisterLambdaName: handlers.InitializeLambda(
			stack, s3Bucket, api.RegisterLambdaName, clientID, userPoolArn, environment,
		),
		api.LoginLambdaName: handlers.InitializeLambda(
			stack, s3Bucket, api.LoginLambdaName, clientID, userPoolArn, environment,
		),
		api.UploadLambdaName: handlers.InitializeLambda(
			stack, s3Bucket, api.UploadLambdaName, "", "", environment,
		),
		api.DownloadLambdaName: handlers.InitializeLambda(
			stack, s3Bucket, api.DownloadLambdaName, "", "", environment,
		),
		api.ListLambdaName: handlers.InitializeLambda(
			stack, s3Bucket, api.ListLambdaName, "", "", environment,
		),
		api.DeleteLambdaName: handlers.InitializeLambda(
			stack, s3Bucket, api.DeleteLambdaName, "", "", environment,
		),
	}

	userPool := cognito.InitializeUserPool(stack, "UserPool", userPoolArn)
	cognitoAuthorizer := cognito.InitializeCognitoAuthorizer(stack, "MentorshipCognitoAuthorizer", userPool)

	api.InitializeAPI(stack, lambdas, cognitoAuthorizer, environment)

	return stack
}

func main() {
	defer jsii.Close()

	if err := config.InitCognitoClient(); err != nil {
		log.Fatalf("Failed to initialize Cognito client: %v", err)
	}

	app := awscdk.NewApp(nil)
	awsContext := &awscdk.Environment{
		Account: jsii.String(os.Getenv("ACCOUNT")),
		Region:  jsii.String(os.Getenv("REGION")),
	}

	stackInitializer(app, "mentorship-staging", &awscdk.StackProps{Env: awsContext}, "staging")
	stackInitializer(app, "mentorship-production", &awscdk.StackProps{Env: awsContext}, "production")

	app.Synth(nil)
}
