package permissions

import (
	"fmt"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/jsii-runtime-go"
	"log"
	"mentorship-app-backend/api"
	"mentorship-app-backend/config"
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

func GrantCognitoLoginPermissions(lambdaFunction awslambda.Function) {
	lambdaFunction.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("cognito-idp:AdminInitiateAuth"),
		Resources: jsii.Strings("*"),
	}))
}

func GetCognitoSettings(environment string) (userPoolArn, clientID string, err error) {
	err = config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	switch environment {
	case config.AppConfig.Environment.Staging:
		return config.AppConfig.Environment.Cognito.StagingPoolArn, config.AppConfig.Environment.Cognito.StagingClientID, nil
	case config.AppConfig.Environment.Production:
		return config.AppConfig.Environment.Cognito.ProductionPoolArn, config.AppConfig.Environment.Cognito.ProductionClientID, nil
	default:
		return "", "", fmt.Errorf("unknown environment: %s", environment)
	}
}
