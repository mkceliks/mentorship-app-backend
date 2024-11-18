package handlers

import (
	"fmt"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/jsii-runtime-go"
	"log"
	"mentorship-app-backend/api"
	"mentorship-app-backend/config"
	"mentorship-app-backend/permissions"
)

func InitializeLambda(stack awscdk.Stack, bucket awss3.Bucket, table awsdynamodb.Table, functionName string, dependentLambdas map[string]awslambda.Function, cfg config.Config) awslambda.Function {
	fullFunctionName := fmt.Sprintf("%s-%s", functionName, cfg.Environment)

	envVars := getLambdaEnvironmentVars(cfg.CognitoClientID, cfg.CognitoPoolArn, cfg.Environment, *bucket.BucketName(), *table.TableName())

	log.Printf("env vars: %v", envVars)

	lambdaFunction := awslambda.NewFunction(stack, jsii.String(fullFunctionName), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Handler:      jsii.String("bootstrap"),
		FunctionName: jsii.String(fullFunctionName),
		Code:         awslambda.Code_FromAsset(jsii.String(fmt.Sprintf("./output/%s_function.zip", functionName)), nil),
		Environment:  &envVars,
		Timeout:      awscdk.Duration_Seconds(jsii.Number(15)),
	})

	grantPermissions(lambdaFunction, dependentLambdas, functionName, bucket, table, cfg)

	return lambdaFunction
}

func getLambdaEnvironmentVars(cognitoClientID, arn, environment, bucketName, tableName string) map[string]*string {
	return map[string]*string{
		"BUCKET_NAME":              jsii.String(bucketName),
		"ENVIRONMENT":              jsii.String(environment),
		"COGNITO_CLIENT_ID":        jsii.String(cognitoClientID),
		"COGNITO_POOL_ARN":         jsii.String(arn),
		"ACCOUNT":                  jsii.String(config.AppConfig.Account),
		"REGION":                   jsii.String(config.AppConfig.Region),
		"SLACK_WEBHOOK_SECRET_ARN": jsii.String(config.AppConfig.SlackWebhookSecretARN),
		"DDB_TABLE_NAME":           jsii.String(tableName),
	}
}

func grantPermissions(lambdaFunction awslambda.Function, dependentLambdas map[string]awslambda.Function, functionName string, bucket awss3.Bucket, table awsdynamodb.Table, cfg config.Config) {
	switch functionName {
	case api.RegisterLambdaName:
		permissions.GrantCognitoRegisterPermissions(lambdaFunction)
		if uploadLambda, exists := dependentLambdas[api.UploadLambdaName]; exists {
			permissions.GrantLambdaInvokePermission(lambdaFunction, uploadLambda)
		}
	case api.LoginLambdaName:
		permissions.GrantCognitoLoginPermissions(lambdaFunction, cfg.CognitoPoolArn)
	case api.ConfirmLambdaName:
		permissions.GrantCognitoConfirmationPermissions(lambdaFunction, cfg.CognitoPoolArn)
	default:
		permissions.GrantAccessForBucket(lambdaFunction, bucket, functionName)
		permissions.GrantCognitoDescribePermissions(lambdaFunction, cfg.CognitoPoolArn)
		permissions.GrantCognitoTokenValidationPermissions(lambdaFunction, cfg.CognitoPoolArn)
	}

	permissions.GrantDynamoDBPermissions(lambdaFunction, table)
	permissions.GrantSecretManagerReadWritePermissions(lambdaFunction, cfg.SlackWebhookSecretARN)
}
