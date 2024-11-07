package notifier

import (
	"github.com/slack-go/slack"
	"log"
)

func NotifySlack(token, baseChannel, message string, fields []slack.AttachmentField, level string) {

	color := "#36a64f"
	switch level {
	case "warning":
		color = "#FFA500"
	case "error":
		color = "#FF0000"
	}

	api := slack.New(token)

	_, _, _, err := api.JoinConversation(baseChannel)
	if err != nil {
		log.Printf("Unable to join conversation %s: %v", baseChannel, err)
	}

	attachment := slack.Attachment{
		Color:  color,
		Text:   message,
		Fields: fields,
	}

	_, _, postErr := api.PostMessage(baseChannel, slack.MsgOptionAttachments(attachment))
	if err != nil {
		log.Printf("Failed to send message to Slack channel %s: %v", baseChannel, postErr)
	}
}
