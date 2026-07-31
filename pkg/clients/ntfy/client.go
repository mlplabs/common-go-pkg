package ntfy

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
)

type HttpClient struct {
	client *resty.Client
}

func NewHttpClient(clientName string, cfg *Config) *HttpClient {
	client := resty.New()
	client.SetAuthToken(cfg.Token)
	client.SetTimeout(cfg.Timeout)

	client.SetBaseURL(cfg.BaseURL).
		SetHeader("Content-Type", "application/json")
	return &HttpClient{
		client: client,
	}
}

func (c *HttpClient) SendMessage(ctx context.Context, message *Message) error {
	req := c.client.R().
		SetContext(ctx).
		SetBody(strings.NewReader(message.Text))

	if message.Title != nil {
		req.SetHeader("Title", *message.Title)
	}

	if message.Priority != nil {
		priority := strconv.Itoa(int(*message.Priority))
		req.Header.Set("Priority", priority)
	}
	if message.Tags != nil {
		req.Header.Set("Tags", *message.Tags)
	}
	if message.Click != nil {
		req.Header.Set("Click", *message.Click)
	}
	resp, err := req.Post(message.Topic)

	if err != nil {
		return fmt.Errorf("failed to sent message, err: %w", err)
	}
	if resp.IsError() {
		if resp.StatusCode() == 404 {
			return fmt.Errorf("topic not found, err: %w", err)
		}
		if resp.StatusCode() == 403 {
			return fmt.Errorf("forbidden, err: %w", err)
		}
	}
	return nil
}
