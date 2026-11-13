// Package slack は Slack の Incoming Webhook でメッセージを送信する。
package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Webhook は Incoming Webhook の送信先。
type Webhook struct {
	URL  string
	HTTP *http.Client
}

func NewWebhook(url string) *Webhook {
	return &Webhook{URL: url, HTTP: &http.Client{Timeout: 30 * time.Second}}
}

// Post は text を送信する。text は Slack の mrkdwn として解釈される。
func (w *Webhook) Post(ctx context.Context, text string) error {
	body, err := json.Marshal(map[string]string{"text": text})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.URL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := w.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(res.Body, 1024))
		return fmt.Errorf("slack webhook: status %d: %s", res.StatusCode, b)
	}
	return nil
}
