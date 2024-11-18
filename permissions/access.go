package permissions

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
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

func GrantCognitoConfirmationPermissions(lambdaFunction awslambda.Function, cognitoPoolArn string) {
	policy := awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("cognito-idp:ConfirmSignUp", "cognito-idp:DescribeUserPool"),
		Resources: jsii.Strings(cognitoPoolArn),
	})
	lambdaFunction.AddToRolePolicy(policy)
}

func GrantPublicReadAccess(bucket awss3.Bucket) {
	bucket.AddToResourcePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Effect: awsiam.Effect_ALLOW,
		Actions: jsii.Strings(
			"s3:GetObject",
		),
		Resources: jsii.Strings(
			*bucket.BucketArn() + "/*",
		),
		Principals: &[]awsiam.IPrincipal{
			awsiam.NewAnyPrincipal(),
		},
	}))
}

func GrantLambdaInvokePermission(lambdaFunction, targetLambda awslambda.Function) {
	lambdaFunction.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("lambda:InvokeFunction"),
		Resources: jsii.Strings(*targetLambda.FunctionArn()),
	}))
}

func GrantCognitoRegisterPermissions(lambdaFunction awslambda.Function) {
	lambdaFunction.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions: jsii.Strings(
			"cognito-idp:SignUp",
			"cognito-idp:AdminCreateUser",
			"cognito-idp:AdminDeleteUser",
			"cognito-idp:AdminUpdateUserAttributes",
		),
		Resources: jsii.Strings("*"),
	}))
}

func GrantCognitoLoginPermissions(lambdaFunction awslambda.Function, cognitoPoolArn string) {
	lambdaFunction.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("cognito-idp:AdminInitiateAuth", "cognito-idp:AdminGetUser"),
		Resources: jsii.Strings(cognitoPoolArn),
	}))
}

func GrantCognitoResendPermissions(lambdaFunction awslambda.Function, cognitoPoolArn string) {
	lambdaFunction.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Effect: awsiam.Effect_ALLOW,
		Actions: jsii.Strings(
			"cognito-idp:ResendConfirmationCode",
		),
		Resources: jsii.Strings(cognitoPoolArn),
	}))
}

func GrantSecretManagerReadWritePermissions(lambdaFunction awslambda.Function, secretArn string) {
	lambdaFunction.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("secretsmanager:GetSecretValue", "secretsmanager:PutSecretValue"),
		Resources: jsii.Strings(secretArn),
	}))
}

func GrantDynamoDBPermissions(lambdaFunction awslambda.Function, table awsdynamodb.Table) {
	table.GrantReadWriteData(lambdaFunction)
}

func GrantDynamoDBStreamPermissions(lambdaFunction awslambda.Function, table awsdynamodb.Table) {
	table.GrantStreamRead(lambdaFunction)
}

func GrantCognitoDescribePermissions(lambdaFunction awslambda.Function, userPoolArn string) {
	lambdaFunction.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Effect:    awsiam.Effect_ALLOW,
		Actions:   jsii.Strings("cognito-idp:DescribeUserPool", "cognito-idp:ListUsers"),
		Resources: jsii.Strings(userPoolArn),
	}))
}

func GrantCognitoTokenValidationPermissions(lambdaFunction awslambda.Function, userPoolArn string) {
	lambdaFunction.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("cognito-idp:AdminGetUser", "cognito-idp:GetSigningCertificate"),
		Resources: jsii.Strings(userPoolArn),
	}))
}

func GrantSNSPublishPermissions(lambdaFunction awslambda.Function, topicArn string) {
	lambdaFunction.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("sns:Publish"),
		Resources: jsii.Strings(topicArn),
	}))
}

func GrantSQSPermissions(lambdaFunction awslambda.Function, queue awssqs.Queue) {
	queue.GrantConsumeMessages(lambdaFunction)
	queue.GrantSendMessages(lambdaFunction)
}

func GrantCloudWatchLogsPermissions(lambdaFunction awslambda.Function) {
	lambdaFunction.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("logs:CreateLogGroup", "logs:CreateLogStream", "logs:PutLogEvents"),
		Resources: jsii.Strings("*"),
	}))
}
