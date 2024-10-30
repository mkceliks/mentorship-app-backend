package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"log"
	"mentorship-app-backend/api"
	"mentorship-app-backend/components/bucket"
	"mentorship-app-backend/config"
	"mentorship-app-backend/handlers"
	"mentorship-app-backend/permissions"
)

// TODO: refactor stack init ( errorpackage and logs )
// TODO: set alarms and router for API gateway
func stackInitializer(
	scope constructs.Construct,
	id string,
	props *awscdk.StackProps,
	environment string,
) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, props)
	s3Bucket := bucket.InitializeBucket(stack, environment)

	api.InitializeAPI(
		stack,
		handlers.InitializeLambda(stack, s3Bucket, handlers.UploadLambdaName),
		handlers.InitializeLambda(stack, s3Bucket, handlers.DownloadLambdaName),
		handlers.InitializeLambda(stack, s3Bucket, handlers.ListLambdaName),
		handlers.InitializeLambda(stack, s3Bucket, handlers.DeleteLambdaName),
		environment,
	)

	return stack
}

func main() {
	defer jsii.Close()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	app := awscdk.NewApp(nil)

	awsContext := permissions.GetAWSEnv(*cfg)

	if awsContext.Account == nil || awsContext.Region == nil {
		panic("aws account and region must be set in the context.")
	}

	stackInitializer(app,
		cfg.Environment.AppName+cfg.Environment.Staging,
		&awscdk.StackProps{
			Env: awsContext,
		}, cfg.Environment.Staging,
	)

	stackInitializer(app,
		cfg.Environment.AppName+cfg.Environment.Production,
		&awscdk.StackProps{
			Env: awsContext,
		}, cfg.Environment.Production)

	app.Synth(nil)
}
