package permissions

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/jsii-runtime-go"
	"mentorship-app-backend/api"
)

func GrantAccessForBucket(lambda awslambda.Function, bucket awss3.Bucket, functionName string) {
	switch functionName {
	case api.UploadLambdaName, api.DeleteLambdaName:
		bucket.GrantReadWrite(lambda, "*")
	case api.DownloadLambdaName, api.ListLambdaName:
		bucket.GrantRead(lambda, "*")
	}
}

func GrantLambdaInvokePermission(lambdaFunction, targetLambda awslambda.Function) {
	lambdaFunction.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("lambda:InvokeFunction"),
		Resources: jsii.Strings(*targetLambda.FunctionArn()),
	}))
}

func GrantCognitoRegisterPermissions(lambdaFunction awslambda.Function) {
	lambdaFunction.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("cognito-idp:SignUp", "cognito-idp:AdminCreateUser", "cognito-idp:AdminDeleteUser"),
		Resources: jsii.Strings("*"),
	}))
}

func GrantCognitoLoginPermissions(lambdaFunction awslambda.Function) {
	lambdaFunction.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("cognito-idp:AdminInitiateAuth"),
		Resources: jsii.Strings("*"),
	}))
}

func GrantSecretManagerReadWritePermissions(lambdaFunction awslambda.Function, secretArn string) {
	lambdaFunction.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("secretsmanager:GetSecretValue"),
		Resources: jsii.Strings(secretArn),
	}))
}

func GrantDynamoDBPermissions(lambdaFunction awslambda.Function, table awsdynamodb.Table) {
	table.GrantReadWriteData(lambdaFunction)
}

func GrantCognitoDescribePermissions(lambdaFunction awslambda.Function, userPoolArn string) {
	lambdaFunction.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Effect:    awsiam.Effect_ALLOW,
		Actions:   jsii.Strings("cognito-idp:DescribeUserPool", "cognito-idp:ListUsers"),
		Resources: jsii.Strings(userPoolArn),
	}))
}
