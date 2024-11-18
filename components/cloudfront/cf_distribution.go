package cloudfront

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfront"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfrontorigins"
	"github.com/aws/jsii-runtime-go"
)

func CreateCloudFrontDistribution(stack awscdk.Stack, api awsapigateway.RestApi, environment string) awscloudfront.Distribution {
	apiDomain := getDomainName(api)

	distribution := awscloudfront.NewDistribution(stack, jsii.String(fmt.Sprintf("CloudFrontDistribution-%s", environment)), &awscloudfront.DistributionProps{
		DefaultBehavior: &awscloudfront.BehaviorOptions{
			Origin:         awscloudfrontorigins.NewHttpOrigin(jsii.String(apiDomain), nil),
			AllowedMethods: awscloudfront.AllowedMethods_ALLOW_ALL(),
			CachePolicy:    awscloudfront.CachePolicy_CACHING_DISABLED(),
		},
		AdditionalBehaviors: &map[string]*awscloudfront.BehaviorOptions{
			"/public/*": {
				Origin:         awscloudfrontorigins.NewHttpOrigin(jsii.String(apiDomain), nil),
				AllowedMethods: awscloudfront.AllowedMethods_ALLOW_ALL(),
				CachePolicy:    awscloudfront.CachePolicy_CACHING_OPTIMIZED(),
			},
			"/protected/*": {
				Origin:         awscloudfrontorigins.NewHttpOrigin(jsii.String(apiDomain), nil),
				AllowedMethods: awscloudfront.AllowedMethods_ALLOW_ALL(),
				CachePolicy:    awscloudfront.CachePolicy_CACHING_DISABLED(),
			},
		},
	})

	awscdk.NewCfnOutput(stack, jsii.String(fmt.Sprintf("CloudFrontDistributionUrl-%s", environment)), &awscdk.CfnOutputProps{
		Value:       distribution.DistributionDomainName(),
		Description: jsii.String("CloudFront Distribution URL"),
		ExportName:  jsii.String(fmt.Sprintf("CloudFrontDistributionUrl-%s", environment)),
	})

	return distribution
}

func getDomainName(api awsapigateway.RestApi) string {
	return fmt.Sprintf("%s.execute-api.%s.amazonaws.com", *api.RestApiId(), awscdk.Stack_Of(api).Region())
}
