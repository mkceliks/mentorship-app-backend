package wrapper

import (
	"fmt"
	"log"
	"os"

	"mentorship-app-backend/components/errorpackage"
	"mentorship-app-backend/components/notifier"
	"mentorship-app-backend/components/secrets"

	"github.com/aws/aws-lambda-go/events"
	"github.com/slack-go/slack"
)

var (
	slackWebhookARN = os.Getenv("SLACK_WEBHOOK_SECRET_ARN")
	environment     = os.Getenv("ENVIRONMENT")
)

func HandlerWrapper(handler func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error), baseChannel, handlerName string) func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		response, err := handler(request)

		slackToken, slackErr := secrets.GetSecretValue(slackWebhookARN)
		if slackErr != nil {
			log.Printf("Failed to retrieve Slack webhook token: %v", slackErr)
			return errorpackage.ServerError(fmt.Sprintf("Internal server error: %v", slackErr))
		}

		var message, level, channel string
		var fields []slack.AttachmentField

		switch {
		case response.StatusCode >= 200 && response.StatusCode < 300:
			channel = getEnvironmentChannel(baseChannel)
			message = fmt.Sprintf("%s executed successfully", handlerName)
			level = "info"
			fields = []slack.AttachmentField{
				{Title: "Handler", Value: handlerName, Short: true},
				{Title: "Status", Value: "Success", Short: true},
				{Title: "Environment", Value: environment, Short: true},
			}
		default:
			channel = getEnvironmentChannel(baseChannel + "-alerts")
			message = fmt.Sprintf("%s execution failed", handlerName)
			level = "error"
			fields = []slack.AttachmentField{
				{Title: "Handler", Value: handlerName, Short: true},
				{Title: "Status", Value: "Failure", Short: true},
				{Title: "Environment", Value: environment, Short: true},
			}
			if err != nil {
				fields = append(fields, slack.AttachmentField{Title: "Error", Value: err.Error(), Short: false})
			}
		}

		notifier.NotifySlack(slackToken, channel, message, fields, level)

		return response, err
	}
}

func getEnvironmentChannel(baseChannel string) string {
	if environment == "staging" {
		return baseChannel + "-staging"
	}
	return baseChannel
}
