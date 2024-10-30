package permissions

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"mentorship-app-backend/api"
	"mentorship-app-backend/config"
)

func GetAWSEnv(cfg config.Config) *awscdk.Environment {
	return &awscdk.Environment{
		Account: &cfg.Context.Account,
		Region:  &cfg.Context.Region,
	}
}

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
