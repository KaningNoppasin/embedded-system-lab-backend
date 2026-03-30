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

const (
	linePushMessageURL = "https://api.line.me/v2/bot/message/push"
	lineAPITimeout     = 10 * time.Second
)

type LineNotifier struct {
	channelAccessToken string
	userID             string
	httpClient         *http.Client
}

type linePushMessageRequest struct {
	To       string            `json:"to"`
	Messages []lineTextMessage `json:"messages"`
}

type lineTextMessage struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func NewLineNotifier() *LineNotifier {
	return &LineNotifier{
		channelAccessToken: strings.TrimSpace(os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")),
		userID:             strings.TrimSpace(os.Getenv("LINE_USER_ID")),
		httpClient:         &http.Client{Timeout: lineAPITimeout},
	}
}

func (n *LineNotifier) Send(content string) error {
	if n == nil || n.httpClient == nil {
		return fmt.Errorf("line notifier is unavailable")
	}
	if n.channelAccessToken == "" {
		return fmt.Errorf("LINE_CHANNEL_ACCESS_TOKEN is not configured")
	}
	if n.userID == "" {
		return fmt.Errorf("LINE_USER_ID is not configured")
	}

	payload, err := json.Marshal(linePushMessageRequest{
		To: n.userID,
		Messages: []lineTextMessage{
			{
				Type: "text",
				Text: content,
			},
		},
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, linePushMessageURL, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+n.channelAccessToken)

	resp, err := n.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("line messaging api returned status %d", resp.StatusCode)
	}

	return nil
}
