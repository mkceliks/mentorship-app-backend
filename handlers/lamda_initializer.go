package handlers

import (
	"fmt"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/jsii-runtime-go"
	"mentorship-app-backend/api"
	"mentorship-app-backend/config"
	"mentorship-app-backend/permissions"
)

func InitializeLambda(stack awscdk.Stack, bucket awss3.Bucket, functionName, cognitoClientID, environment string) awslambda.Function {
	envVars := getLambdaEnvironmentVars(functionName, cognitoClientID, environment, *bucket.BucketName())

	lambdaFunction := awslambda.NewFunction(stack, jsii.String(functionName), &awslambda.FunctionProps{
		Runtime:     awslambda.Runtime_PROVIDED_AL2(),
		Handler:     jsii.String("bootstrap"),
		Code:        awslambda.Code_FromAsset(jsii.String(fmt.Sprintf("./output/%s_function.zip", functionName)), nil),
		Environment: &envVars,
	})

	switch functionName {
	case api.RegisterLambdaName:
		permissions.GrantCognitoRegisterPermissions(lambdaFunction)
		permissions.ConfigureLambdaEnvironment(lambdaFunction, cognitoClientID)

	case api.LoginLambdaName:
		permissions.GrantCognitoLoginPermissions(lambdaFunction)
		permissions.ConfigureLambdaEnvironment(lambdaFunction, cognitoClientID)

	default:
		permissions.GrantAccessForBucket(lambdaFunction, bucket, functionName)
	}

	return lambdaFunction
}

func getLambdaEnvironmentVars(functionName, cognitoClientID, environment, bucketName string) map[string]*string {
	envVars := map[string]*string{
		"BUCKET_NAME": jsii.String(bucketName),
		"ENVIRONMENT": jsii.String(environment),
		"APP_NAME":    jsii.String(config.AppConfig.Environment.AppName),
		"ACCOUNT":     jsii.String(config.AppConfig.Context.Account),
		"REGION":      jsii.String(config.AppConfig.Context.Region),
		"ROUTE_NAME":  jsii.String(config.AppConfig.RouteName),
	}

	switch environment {
	case config.AppConfig.Environment.Staging:
		envVars["COGNITO_POOL_ARN"] = jsii.String(config.AppConfig.Environment.Cognito.StagingPoolArn)
		envVars["COGNITO_CLIENT_ID"] = jsii.String(config.AppConfig.Environment.Cognito.StagingClientID)
	case config.AppConfig.Environment.Production:
		envVars["COGNITO_POOL_ARN"] = jsii.String(config.AppConfig.Environment.Cognito.ProductionPoolArn)
		envVars["COGNITO_CLIENT_ID"] = jsii.String(config.AppConfig.Environment.Cognito.ProductionClientID)
	}

	if functionName == api.RegisterLambdaName || functionName == api.LoginLambdaName {
		envVars["COGNITO_CLIENT_ID"] = jsii.String(cognitoClientID)
	}

	return envVars
}
