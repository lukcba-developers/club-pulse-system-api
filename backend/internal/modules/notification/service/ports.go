package service

import "context"

type DeliveryResult struct {
	Success      bool
	ProviderID   string
	ErrorMessage string
}

type EmailProvider interface {
	SendEmail(ctx context.Context, to string, subject string, body string) (*DeliveryResult, error)
}

type SMSProvider interface {
	SendSMS(ctx context.Context, to string, message string) (*DeliveryResult, error)
}
