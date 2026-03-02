package service

import (
	"context"
)

type MailProvider interface {
	Send(ctx context.Context, to string, subject string, body string) error
}

type MultiMailProvider struct {
	providers []MailProvider
}

func NewMultiMailProvider(providers ...MailProvider) MailProvider {
	return &MultiMailProvider{providers: providers}
}

func (m *MultiMailProvider) Send(ctx context.Context, to, subject, body string) error {
	errCh := make(chan error, len(m.providers))

	for _, p := range m.providers {
		go func(provider MailProvider) {
			errCh <- provider.Send(ctx, to, subject, body)
		}(p)
	}

	var errs []error
	for i := 0; i < len(m.providers); i++ {
		if err := <-errCh; err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}
