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

	_, _, _, err := api.JoinConversation(channel)
	if err != nil {
		log.Printf("Unable to join conversation %s: %v", channel, err)
	}

	attachment := slack.Attachment{
		Color:  color,
		Text:   message,
		Fields: fields,
	}

	_, _, postErr := api.PostMessage(channel, slack.MsgOptionAttachments(attachment))
	if err != nil {
		log.Printf("Failed to send message to Slack channel %s: %v", channel, postErr)
	}
}
