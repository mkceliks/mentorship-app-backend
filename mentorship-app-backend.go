package main

import (
	"fmt"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type MentorshipAppBackendStackProps struct {
	awscdk.StackProps
}

func NewMentorshipAppBackendStack(scope constructs.Construct, id string, props *MentorshipAppBackendStackProps, isProduction bool) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)

	bucketName := "MentorshipAppBucket-Staging"
	if isProduction {
		bucketName = "MentorshipAppBucket-Production"
	}

	bucket := awss3.NewBucket(stack, jsii.String(bucketName), &awss3.BucketProps{
		Versioned: jsii.Bool(true),
	})

	uploadLambda := awslambda.NewFunction(stack, jsii.String("UploadLambda"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2(),
		Handler: jsii.String("bootstrap"),
		Code:    awslambda.Code_FromAsset(jsii.String("./function.zip"), nil),
		Environment: &map[string]*string{
			"BUCKET_NAME": bucket.BucketName(),
		},
	})

	bucket.GrantReadWrite(uploadLambda, jsii.String("*"))

	bucket.AddToResourcePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("s3:PutObject", "s3:GetObject"),
		Resources: jsii.Strings(fmt.Sprintf("%s/*", *bucket.BucketArn())),
		Principals: &[]awsiam.IPrincipal{
			awsiam.NewArnPrincipal(uploadLambda.Role().RoleArn()),
		},
	}))

	api := awsapigateway.NewRestApi(stack, jsii.String("MentorshipAppAPI"), &awsapigateway.RestApiProps{
		RestApiName: jsii.String("MentorshipAppAPI"),
		Description: jsii.String("API Gateway for handling S3 file uploads."),
	})

	upload := api.Root().AddResource(jsii.String("upload"), nil)
	upload.AddMethod(jsii.String("POST"), awsapigateway.NewLambdaIntegration(uploadLambda, nil), nil)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewMentorshipAppBackendStack(app, "MentorshipAppBackendStagingStack", &MentorshipAppBackendStackProps{
		awscdk.StackProps{
			Env: envStaging(),
		},
	}, false)

	NewMentorshipAppBackendStack(app, "MentorshipAppBackendProductionStack", &MentorshipAppBackendStackProps{
		awscdk.StackProps{
			Env: envProduction(),
		},
	}, true)

	app.Synth(nil)
}

func envStaging() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String("034362052544"),
		Region:  jsii.String("us-east-1"),
	}
}

func envProduction() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String("034362052544"),
		Region:  jsii.String("us-east-1"),
	}
}
