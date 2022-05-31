package emailmanager

import (
	"context"
	"log"
)

// Mocked mail sender, only logs the mail that would have been sent

type MailSender interface {
	SendMail(ctx context.Context, email string, title string, content string) error
}

func NewMockMailSender() MailSender {
	return &mockMailSender{}
}

type mockMailSender struct {
}

func (m *mockMailSender) SendMail(ctx context.Context, email string, title string, content string) error {
	log.Printf("\nSending mail to %s: %s\n", email, title)
	return nil
}
