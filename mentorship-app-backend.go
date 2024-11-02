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
	"mentorship-app-backend/permissions"
)

func stackInitializer(
	scope constructs.Construct,
	id string,
	props *awscdk.StackProps,
	environment string,
) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, props)

	userPoolArn, clientID, err := permissions.GetCognitoSettings(environment)
	if err != nil {
		log.Fatalf("Failed to get Cognito settings: %v", err)
	}

	userPool := cognito.InitializeUserPool(stack, "ExistingUserPool", userPoolArn)
	userPoolClient := cognito.InitializeUserPoolClient(stack, "ExistingUserPoolClient", clientID)

	cognitoAuthorizer := cognito.InitializeCognitoAuthorizer(stack, "MentorshipCognitoAuthorizer", userPool)

	s3Bucket := bucket.InitializeBucket(stack, environment)

	// lambdas
	lambdas := map[string]awslambda.Function{
		api.RegisterLambdaName: handlers.InitializeLambda(
			stack, s3Bucket, api.RegisterLambdaName, *userPoolClient.UserPoolClientId(), environment,
		),
		api.LoginLambdaName: handlers.InitializeLambda(
			stack, s3Bucket, api.LoginLambdaName, *userPoolClient.UserPoolClientId(), environment,
		),
		api.UploadLambdaName: handlers.InitializeLambda(
			stack, s3Bucket, api.UploadLambdaName, "", environment,
		),
		api.DownloadLambdaName: handlers.InitializeLambda(
			stack, s3Bucket, api.DownloadLambdaName, "", environment,
		),
		api.ListLambdaName: handlers.InitializeLambda(
			stack, s3Bucket, api.ListLambdaName, "", environment,
		),
		api.DeleteLambdaName: handlers.InitializeLambda(
			stack, s3Bucket, api.DeleteLambdaName, "", environment,
		),
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
