package handlers

import (
	"fmt"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/jsii-runtime-go"
	"mentorship-app-backend/api"
	"mentorship-app-backend/permissions"
)

func InitializeLambda(stack awscdk.Stack, bucket awss3.Bucket, functionName string, cognitoClientID string) awslambda.Function {
	lambdaFunction := awslambda.NewFunction(stack, jsii.String(functionName), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2(),
		Handler: jsii.String("bootstrap"),
		Code:    awslambda.Code_FromAsset(jsii.String(fmt.Sprintf("./output/%s_function.zip", functionName)), nil),
		Environment: &map[string]*string{
			"BUCKET_NAME": bucket.BucketName(),
		},
	})

	switch functionName {
	case api.RegisterLambdaName:
		permissions.GrantCognitoRegisterPermissions(lambdaFunction)
		permissions.ConfigureLambdaEnvironment(lambdaFunction, cognitoClientID)

	case api.LoginLambdaName:
		permissions.GrantCognitoSignInPermissions(lambdaFunction)
		permissions.ConfigureLambdaEnvironment(lambdaFunction, cognitoClientID)

	default:
		permissions.GrantAccessForBucket(lambdaFunction, bucket, functionName)
	}

	return lambdaFunction
}
