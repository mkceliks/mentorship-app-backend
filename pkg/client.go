package pkg

import (
	"github.com/go-resty/resty/v2"
	"mentorship-app-backend/config"
)

type Client struct {
	client *resty.Client
}

func NewClient() *Client {
	client := resty.New()

	if config.AppConfig.EndpointBaseURL != "" {
		client.SetBaseURL(config.AppConfig.EndpointBaseURL)
	}

	client.SetHeader("Content-Type", "application/json")

	return &Client{client: client}
}

func (c *Client) SetHeader(key, value string) {
	c.client.SetHeader(key, value)
}

func (c *Client) SetHeaders(headers map[string]string) {
	for key, value := range headers {
		c.client.SetHeader(key, value)
	}
}
