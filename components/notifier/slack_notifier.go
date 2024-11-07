package notifier

import (
	"log"
	"os"

	"github.com/slack-go/slack"
)

var environment = os.Getenv("ENVIRONMENT")

func NotifySlack(token, baseChannel, message string, fields []slack.AttachmentField, level string) {
	channel := baseChannel
	if environment == "staging" {
		channel = baseChannel + "-staging"
	}

	color := "#36a64f"
	switch level {
	case "warning":
		color = "#FFA500"
	case "error":
		color = "#FF0000"
	}

	api := slack.New(token)

	attachment := slack.Attachment{
		Color:  color,
		Text:   message,
		Fields: fields,
	}

	_, _, err := api.PostMessage(channel, slack.MsgOptionAttachments(attachment))
	if err != nil {
		log.Printf("Failed to send message to Slack channel %s: %v", channel, err)
	}
}
