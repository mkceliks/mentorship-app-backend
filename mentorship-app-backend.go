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
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
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

	log.Printf("Initializing stack for environment: %s", cfg.Environment)

	s3Bucket := bucket.InitializeBucket(stack, cfg.BucketName)
	fmt.Printf("Bucket Name: %s\n", *s3Bucket.BucketName())

	removalPolicy := awscdk.RemovalPolicy_RETAIN
	if cfg.Environment == "staging" {
		removalPolicy = awscdk.RemovalPolicy_DESTROY
	}

	profileTable := awsdynamodb.NewTable(stack, jsii.String("UserProfiles"), &awsdynamodb.TableProps{
		TableName:     jsii.String("UserProfiles"),
		PartitionKey:  &awsdynamodb.Attribute{Name: jsii.String("UserId"), Type: awsdynamodb.AttributeType_STRING},
		SortKey:       &awsdynamodb.Attribute{Name: jsii.String("ProfileType"), Type: awsdynamodb.AttributeType_STRING},
		BillingMode:   awsdynamodb.BillingMode_PAY_PER_REQUEST,
		RemovalPolicy: removalPolicy,
	})

	lambdas := map[string]awslambda.Function{
		api.RegisterLambdaName: handlers.InitializeLambda(stack, s3Bucket, profileTable, api.RegisterLambdaName, cfg),
		api.LoginLambdaName:    handlers.InitializeLambda(stack, s3Bucket, profileTable, api.LoginLambdaName, cfg),
		api.UploadLambdaName:   handlers.InitializeLambda(stack, s3Bucket, profileTable, api.UploadLambdaName, cfg),
		api.DownloadLambdaName: handlers.InitializeLambda(stack, s3Bucket, profileTable, api.DownloadLambdaName, cfg),
		api.ListLambdaName:     handlers.InitializeLambda(stack, s3Bucket, profileTable, api.ListLambdaName, cfg),
		api.DeleteLambdaName:   handlers.InitializeLambda(stack, s3Bucket, profileTable, api.DeleteLambdaName, cfg),
	}

	userPool := cognito.InitializeUserPool(stack, "UserPool", cfg.CognitoPoolArn)
	cognitoAuthorizer := cognito.InitializeCognitoAuthorizer(stack, "MentorshipCognitoAuthorizer", userPool)

	api.InitializeAPI(stack, lambdas, cognitoAuthorizer, cfg.Environment)

	return stack
}

func main() {
	defer jsii.Close()

	environment := getEnvironment()
	cfg, err := config.LoadConfig(environment)
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	app := awscdk.NewApp(nil)
	awsContext := &awscdk.Environment{
		Account: jsii.String(cfg.Account),
		Region:  jsii.String(cfg.Region),
	}

	stackInitializer(app, fmt.Sprintf("mentorship-%s", environment), &awscdk.StackProps{Env: awsContext}, cfg)
	app.Synth(nil)
}
