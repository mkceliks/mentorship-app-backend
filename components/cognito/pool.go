package cognito

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func InitializeUserPool(scope constructs.Construct, id, userPoolArn string) awscognito.IUserPool {
	return awscognito.UserPool_FromUserPoolArn(scope, jsii.String(id), jsii.String(userPoolArn))
}
