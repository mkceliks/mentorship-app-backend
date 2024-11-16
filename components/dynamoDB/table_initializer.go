package dynamoDB

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/jsii-runtime-go"
)

func InitializeProfileTable(stack awscdk.Stack, tableName string, removalPolicy awscdk.RemovalPolicy) awsdynamodb.Table {
	return awsdynamodb.NewTable(stack, jsii.String("UserProfiles"), &awsdynamodb.TableProps{
		TableName:     jsii.String(tableName),
		PartitionKey:  &awsdynamodb.Attribute{Name: jsii.String("UserId"), Type: awsdynamodb.AttributeType_STRING},
		SortKey:       &awsdynamodb.Attribute{Name: jsii.String("ProfileType"), Type: awsdynamodb.AttributeType_STRING},
		BillingMode:   awsdynamodb.BillingMode_PAY_PER_REQUEST,
		RemovalPolicy: removalPolicy,
	})
}
