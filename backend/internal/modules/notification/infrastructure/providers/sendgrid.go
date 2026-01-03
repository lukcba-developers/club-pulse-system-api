package providers

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/notification/service"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridProvider struct {
	client    *sendgrid.Client
	fromName  string
	fromEmail string
}

func NewSendGridProvider(apiKey, fromName, fromEmail string) *SendGridProvider {
	return &SendGridProvider{
		client:    sendgrid.NewSendClient(apiKey),
		fromName:  fromName,
		fromEmail: fromEmail,
	}
}

func (p *SendGridProvider) SendEmail(ctx context.Context, to string, subject string, body string) (*service.DeliveryResult, error) {
	from := mail.NewEmail(p.fromName, p.fromEmail)
	toEmail := mail.NewEmail("", to)

	// Assuming body is plain text for simplicity, or we can treat as HTML
	message := mail.NewSingleEmail(from, subject, toEmail, body, body)

	resp, err := p.client.SendWithContext(ctx, message)
	if err != nil {
		log.Printf("SendGrid Error: %v", err)
		return &service.DeliveryResult{Success: false, ErrorMessage: err.Error()}, err
	}

	if resp.StatusCode >= 400 {
		errMsg := fmt.Sprintf("SendGrid failed with status: %d, body: %s", resp.StatusCode, resp.Body)
		log.Println(errMsg)
		return &service.DeliveryResult{Success: false, ErrorMessage: errMsg}, errors.New(errMsg)
	}

	return &service.DeliveryResult{
		Success:    true,
		ProviderID: resp.Headers["X-Message-Id"][0], // Helper access might need check
	}, nil
}
