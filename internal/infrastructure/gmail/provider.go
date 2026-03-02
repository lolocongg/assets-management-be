package gmail

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/davidcm146/assets-management-be.git/internal/config"
	"github.com/davidcm146/assets-management-be.git/internal/service"
)

type Provider struct {
	from     string
	password string
	host     string
	port     string
}

func NewProvider(cfg *config.GmailConfig) service.MailProvider {
	return &Provider{
		from:     cfg.Email,
		password: cfg.Password,
		host:     cfg.Host,
		port:     cfg.Port,
	}
}

func (p *Provider) Send(ctx context.Context, to, subject, body string) error {
	auth := smtp.PlainAuth(
		"",
		p.from,
		p.password,
		p.host,
	)

	msg := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
			"\r\n"+
			"%s",
		p.from,
		to,
		subject,
		body,
	))

	addr := p.host + ":" + p.port

	return smtp.SendMail(
		addr,
		auth,
		p.from,
		[]string{to},
		msg,
	)
}
