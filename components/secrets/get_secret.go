package secrets

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func GetSecretValue(secretARN string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", fmt.Errorf("unable to load SDK config, %v", err)
	}

	client := secretsmanager.NewFromConfig(cfg)

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretARN),
	}
	result, err := client.GetSecretValue(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve secret: %v", err)
	}

	var secretData map[string]string
	if err = json.Unmarshal([]byte(*result.SecretString), &secretData); err != nil {
		return "", fmt.Errorf("failed to parse secret JSON: %v", err)
	}

	slackWebhookURL, ok := secretData["slack_webhook_url"]
	if !ok {
		return "", fmt.Errorf("slack_webhook_url key not found in secret")
	}

	return slackWebhookURL, nil
}
