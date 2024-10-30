package api

import (
	"fmt"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/jsii-runtime-go"
	"mentorship-app-backend/handlers"
)

const (
	apiRoutes = "api-routes"
)

func InitializeAPI(stack awscdk.Stack, uploadLambda, downloadLambda, listLambda, deleteLambda awslambda.Function, environment string) {
	apiName := fmt.Sprintf(apiRoutes+"%s", environment)

	// define api gateway
	api := awsapigateway.NewRestApi(stack, jsii.String(apiName), &awsapigateway.RestApiProps{
		RestApiName: jsii.String(apiName),
		Description: jsii.String(fmt.Sprintf("API Gateway for %s environment", environment)),
		DeployOptions: &awsapigateway.StageOptions{
			StageName: jsii.String(environment),
		},
		DefaultCorsPreflightOptions: &awsapigateway.CorsOptions{
			AllowOrigins: awsapigateway.Cors_ALL_ORIGINS(),
			AllowMethods: jsii.Strings("OPTIONS", "GET", "POST", "DELETE"),
			AllowHeaders: jsii.Strings("Content-Type", "Authorization"),
		},
	})

	// create routes
	addApiResource(api, "POST", handlers.UploadLambdaName, uploadLambda)
	addApiResource(api, "GET", handlers.DownloadLambdaName, downloadLambda)
	addApiResource(api, "GET", handlers.ListLambdaName, listLambda)
	addApiResource(api, "DELETE", handlers.DeleteLambdaName, deleteLambda)
}

func addApiResource(api awsapigateway.RestApi, method, resourceName string, lambdaFunction awslambda.Function) {
	resource := api.Root().AddResource(jsii.String(resourceName), nil)
	resource.AddMethod(jsii.String(method), awsapigateway.NewLambdaIntegration(lambdaFunction, nil), nil)
}
