package postmark

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/davidcm146/assets-management-be.git/internal/config"
	"github.com/davidcm146/assets-management-be.git/internal/error_middleware"
	"github.com/davidcm146/assets-management-be.git/internal/service"
)

type Provider struct {
	token  string
	from   string
	apiURL string
}

func NewProvider(cfg *config.PostmarkConfig) service.MailProvider {
	return &Provider{
		apiURL: cfg.APIURL,
		token:  cfg.Token,
		from:   cfg.From,
	}
}

func (p *Provider) Send(ctx context.Context, to, subject, body string) error {
	payload := map[string]any{
		"From":          p.from,
		"To":            to,
		"Subject":       subject,
		"HtmlBody":      body,
		"MessageStream": "outbound",
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		p.apiURL,
		bytes.NewReader(b),
	)
	if err != nil {
		return error_middleware.NewInternal("Lỗi gửi email")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Postmark-Server-Token", p.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return error_middleware.NewInternal("Lỗi gửi email")
	}

	return nil
}
