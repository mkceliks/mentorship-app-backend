package cognito

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func InitializeUserPool(scope constructs.Construct, id string) awscognito.UserPool {
	return awscognito.NewUserPool(scope, jsii.String(id), &awscognito.UserPoolProps{
		SignInAliases: &awscognito.SignInAliases{
			Email: jsii.Bool(true),
		},
	})
}
