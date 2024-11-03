package handlers

import (
	"fmt"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/jsii-runtime-go"
	"log"
	"mentorship-app-backend/api"
	"mentorship-app-backend/permissions"
	"os"
)

func InitializeLambda(stack awscdk.Stack, bucket awss3.Bucket, functionName, cognitoClientID, arn, environment string) awslambda.Function {
	envVars := getLambdaEnvironmentVars(cognitoClientID, arn, environment, *bucket.BucketName())

	log.Printf("env vars: %v", envVars)

	lambdaFunction := awslambda.NewFunction(stack, jsii.String(functionName), &awslambda.FunctionProps{
		Runtime:     awslambda.Runtime_PROVIDED_AL2(),
		Handler:     jsii.String("bootstrap"),
		Code:        awslambda.Code_FromAsset(jsii.String(fmt.Sprintf("./output/%s_function.zip", functionName)), nil),
		Environment: &envVars,
	})

	switch functionName {
	case api.RegisterLambdaName:
		permissions.GrantCognitoRegisterPermissions(lambdaFunction)
	case api.LoginLambdaName:
		permissions.GrantCognitoLoginPermissions(lambdaFunction)
	default:
		permissions.GrantAccessForBucket(lambdaFunction, bucket, functionName)
	}

	return lambdaFunction
}

func getLambdaEnvironmentVars(cognitoClientID, arn, environment, bucketName string) map[string]*string {
	return map[string]*string{
		"BUCKET_NAME":       jsii.String(bucketName),
		"ENVIRONMENT":       jsii.String(environment),
		"COGNITO_CLIENT_ID": jsii.String(cognitoClientID),
		"COGNITO_POOL_ARN":  jsii.String(arn),
		"ACCOUNT":           jsii.String(os.Getenv("ACCOUNT")),
		"REGION":            jsii.String(os.Getenv("REGION")),
	}
}
