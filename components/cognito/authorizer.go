package cognito

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func InitializeCognitoAuthorizer(scope constructs.Construct, id string, userPool awscognito.IUserPool) awsapigateway.CognitoUserPoolsAuthorizer {
	return awsapigateway.NewCognitoUserPoolsAuthorizer(scope, jsii.String(id), &awsapigateway.CognitoUserPoolsAuthorizerProps{
		CognitoUserPools: &[]awscognito.IUserPool{userPool},
	})
}
