package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func SendSlackNotification(webhookURL, message string) error {
	payload, _ := json.Marshal(map[string]string{"text": message})

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to Slack: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Slack API error: %s", resp.Status)
	}
	return nil
}
