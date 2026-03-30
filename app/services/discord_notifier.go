package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

const discordWebhookTimeout = 10 * time.Second

type DiscordNotifier struct {
	webhookURL string
	httpClient *http.Client
}

type DiscordMessage struct {
	Content string `json:"content"`
}

func NewDiscordNotifier() *DiscordNotifier {
	return &DiscordNotifier{
		webhookURL: strings.TrimSpace(os.Getenv("DISCORD_WEBHOOK_URL")),
		httpClient: &http.Client{Timeout: discordWebhookTimeout},
	}
}

func (n *DiscordNotifier) Send(content string) error {
	if n == nil || n.httpClient == nil {
		return fmt.Errorf("discord notifier is unavailable")
	}
	if n.webhookURL == "" {
		return fmt.Errorf("DISCORD_WEBHOOK_URL is not configured")
	}

	payload, err := json.Marshal(DiscordMessage{Content: content})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, n.webhookURL, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("discord webhook returned status %d", resp.StatusCode)
	}

	return nil
}
