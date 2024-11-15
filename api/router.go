package api

import (
	"fmt"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/jsii-runtime-go"
)

const (
	UploadLambdaName   = "upload"
	DownloadLambdaName = "download"
	ListLambdaName     = "list"
	DeleteLambdaName   = "delete"
	LoginLambdaName    = "login"
	RegisterLambdaName = "register"
)

func InitializeAPI(stack awscdk.Stack, lambdas map[string]awslambda.Function, cognitoAuthorizer awsapigateway.IAuthorizer, environment string) {
	logGroup := awslogs.NewLogGroup(stack, jsii.String(fmt.Sprintf("APIGatewayLogGroup-%s", environment)), &awslogs.LogGroupProps{
		Retention: awslogs.RetentionDays_ONE_WEEK,
	})

	api := awsapigateway.NewRestApi(stack, jsii.String(fmt.Sprintf("api-gateway-%s", environment)), &awsapigateway.RestApiProps{
		RestApiName: jsii.String(fmt.Sprintf("api-gateway-%s", environment)),
		DefaultCorsPreflightOptions: &awsapigateway.CorsOptions{
			AllowOrigins: awsapigateway.Cors_ALL_ORIGINS(),
			AllowMethods: awsapigateway.Cors_ALL_METHODS(),
			AllowHeaders: jsii.Strings("Content-Type", "Authorization", "x-file-content-type"),
		},
		DeployOptions: &awsapigateway.StageOptions{
			StageName:            jsii.String(environment),
			LoggingLevel:         awsapigateway.MethodLoggingLevel_INFO,
			DataTraceEnabled:     jsii.Bool(true),
			AccessLogDestination: awsapigateway.NewLogGroupLogDestination(logGroup),
			AccessLogFormat: awsapigateway.AccessLogFormat_JsonWithStandardFields(&awsapigateway.JsonWithStandardFieldProps{
				Caller:         jsii.Bool(true),
				HttpMethod:     jsii.Bool(true),
				RequestTime:    jsii.Bool(true),
				ResponseLength: jsii.Bool(true),
				Status:         jsii.Bool(true),
				User:           jsii.Bool(true),
				Protocol:       jsii.Bool(true),
				ResourcePath:   jsii.Bool(true),
				Ip:             jsii.Bool(true),
			}),
		},
	})

	SetupPublicEndpoints(api, lambdas)
	SetupProtectedEndpoints(api, lambdas, cognitoAuthorizer)
}

func SetupPublicEndpoints(api awsapigateway.RestApi, lambdas map[string]awslambda.Function) {
	addApiResource(api, "POST", RegisterLambdaName, lambdas[RegisterLambdaName], nil)
	addApiResource(api, "POST", LoginLambdaName, lambdas[LoginLambdaName], nil)
	addApiResource(api, "POST", UploadLambdaName, lambdas[UploadLambdaName], nil)
}

func SetupProtectedEndpoints(api awsapigateway.RestApi, lambdas map[string]awslambda.Function, cognitoAuthorizer awsapigateway.IAuthorizer) {
	addApiResource(api, "GET", DownloadLambdaName, lambdas[DownloadLambdaName], cognitoAuthorizer)
	addApiResource(api, "GET", ListLambdaName, lambdas[ListLambdaName], cognitoAuthorizer)
	addApiResource(api, "DELETE", DeleteLambdaName, lambdas[DeleteLambdaName], cognitoAuthorizer)
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
