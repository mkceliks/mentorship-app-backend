package cognito

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	"github.com/aws/jsii-runtime-go"
)

func InitializeUserPoolClient(userPool awscognito.UserPool, clientID string) awscognito.UserPoolClient {
	return awscognito.NewUserPoolClient(userPool, jsii.String(clientID), &awscognito.UserPoolClientProps{
		UserPool: userPool,
	})
}
