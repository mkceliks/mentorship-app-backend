package permissions

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/jsii-runtime-go"
	"mentorship-app-backend/api"
)

func GrantAccessForBucket(
	lambda awslambda.Function,
	bucket awss3.Bucket,
	functionName string,
) {
	switch functionName {
	case api.UploadLambdaName, api.DeleteLambdaName:
		bucket.GrantReadWrite(lambda, "*")
	case api.DownloadLambdaName, api.ListLambdaName:
		bucket.GrantRead(lambda, "*")
	}
}

func ConfigureLambdaEnvironment(lambdaFunction awslambda.Function, cognitoClientID string) {
	lambdaFunction.AddEnvironment(jsii.String("COGNITO_CLIENT_ID"), jsii.String(cognitoClientID), nil)
}

func GrantCognitoRegisterPermissions(lambdaFunction awslambda.Function) {
	lambdaFunction.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("cognito-idp:SignUp", "cognito-idp:AdminCreateUser"),
		Resources: jsii.Strings("*"),
	}))
}

func GrantCognitoSignInPermissions(lambdaFunction awslambda.Function) {
	lambdaFunction.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("cognito-idp:AdminInitiateAuth"),
		Resources: jsii.Strings("*"),
	}))
}
