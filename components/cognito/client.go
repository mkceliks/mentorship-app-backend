package cognito

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func InitializeUserPoolClient(scope constructs.Construct, id, clientID string) awscognito.IUserPoolClient {
	return awscognito.UserPoolClient_FromUserPoolClientId(scope, jsii.String(id), jsii.String(clientID))
}
