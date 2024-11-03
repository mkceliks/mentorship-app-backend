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

func validateEnvVars(requiredVars []string) error {
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			return fmt.Errorf("missing required environment variable: %s", v)
		}
		log.Printf("environment variable %s set to %s", v, os.Getenv(v))
	}
	return nil
}

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

	s3Bucket := bucket.InitializeBucket(stack, environment)

	log.Printf("S3 Bucket Name: %s\n", *s3Bucket.BucketName())

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

	err := config.InitCognitoClient()
	if err != nil {
		log.Fatalf("failed to initialize Cognito client: %v", err)
	}

	requiredVars := []string{
		"ACCOUNT",
		"REGION",
		"STAGING_POOL_ARN",
		"PRODUCTION_POOL_ARN",
		"STAGING_CLIENT_ID",
		"PRODUCTION_CLIENT_ID",
		"BUCKET_NAME",
	}

	if err = validateEnvVars(requiredVars); err != nil {
		log.Fatalf("Environment variable validation failed: %v", err)
	}

	app := awscdk.NewApp(nil)
	awsContext := &awscdk.Environment{
		Account: jsii.String(os.Getenv("ACCOUNT")),
		Region:  jsii.String(os.Getenv("REGION")),
	}

	stackInitializer(
		app,
		"mentorship-staging",
		&awscdk.StackProps{Env: awsContext},
		"staging",
	)

	stackInitializer(
		app,
		"mentorship-production",
		&awscdk.StackProps{Env: awsContext},
		"production",
	)

	app.Synth(nil)
}
