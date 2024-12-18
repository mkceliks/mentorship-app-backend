package api

import (
	"fmt"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/jsii-runtime-go"
)

const (
	UploadLambdaName   = "upload"
	DownloadLambdaName = "download"
	ListLambdaName     = "list"
	DeleteLambdaName   = "delete"
	LoginLambdaName    = "login"
	RegisterLambdaName = "register"
	MeLambdaName       = "me"
	ConfirmLambdaName  = "confirm"
	ResendLambdaName   = "resend"
)

func InitializeAPI(stack awscdk.Stack, lambdas map[string]awslambda.Function, cognitoAuthorizer awsapigateway.IAuthorizer, environment string) awsapigateway.RestApi {
	api := awsapigateway.NewRestApi(stack, jsii.String(fmt.Sprintf("api-gateway-%s", environment)), &awsapigateway.RestApiProps{
		RestApiName: jsii.String(fmt.Sprintf("api-gateway-%s", environment)),
		DefaultCorsPreflightOptions: &awsapigateway.CorsOptions{
			AllowOrigins: awsapigateway.Cors_ALL_ORIGINS(),
			AllowMethods: awsapigateway.Cors_ALL_METHODS(),
			AllowHeaders: jsii.Strings("Content-Type", "Authorization", "x-file-content-type"),
		},
		DeployOptions: &awsapigateway.StageOptions{
			StageName: jsii.String(environment),
		},
	})

	SetupPublicEndpoints(api, lambdas)
	SetupProtectedEndpoints(api, lambdas, cognitoAuthorizer)

	return api
}

func SetupPublicEndpoints(api awsapigateway.RestApi, lambdas map[string]awslambda.Function) {
	addApiResource(api, "POST", RegisterLambdaName, lambdas[RegisterLambdaName], nil)
	addApiResource(api, "POST", LoginLambdaName, lambdas[LoginLambdaName], nil)
	addApiResource(api, "POST", UploadLambdaName, lambdas[UploadLambdaName], nil)
	addApiResource(api, "POST", ConfirmLambdaName, lambdas[ConfirmLambdaName], nil)
	addApiResource(api, "GET", ResendLambdaName, lambdas[ResendLambdaName], nil)
}

func SetupProtectedEndpoints(api awsapigateway.RestApi, lambdas map[string]awslambda.Function, cognitoAuthorizer awsapigateway.IAuthorizer) {
	addApiResource(api, "GET", DownloadLambdaName, lambdas[DownloadLambdaName], cognitoAuthorizer)
	addApiResource(api, "GET", ListLambdaName, lambdas[ListLambdaName], cognitoAuthorizer)
	addApiResource(api, "DELETE", DeleteLambdaName, lambdas[DeleteLambdaName], cognitoAuthorizer)
	addApiResource(api, "GET", MeLambdaName, lambdas[MeLambdaName], cognitoAuthorizer)
}

func addApiResource(api awsapigateway.RestApi, method, resourceName string, lambdaFunction awslambda.Function, cognitoAuthorizer awsapigateway.IAuthorizer) {
	resource := api.Root().AddResource(jsii.String(resourceName), nil)
	methodOptions := &awsapigateway.MethodOptions{}
	if cognitoAuthorizer != nil {
		methodOptions = &awsapigateway.MethodOptions{
			AuthorizationType: awsapigateway.AuthorizationType_COGNITO,
			Authorizer:        cognitoAuthorizer,
		}
	}
	resource.AddMethod(jsii.String(method), awsapigateway.NewLambdaIntegration(lambdaFunction, nil), methodOptions)
}
