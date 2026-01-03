package service

import (
	"context"
	"fmt"
)

type NotificationType string

const (
	NotificationTypeEmail NotificationType = "EMAIL"
	NotificationTypePush  NotificationType = "PUSH"
	NotificationTypeSMS   NotificationType = "SMS"
)

type Notification struct {
	RecipientID string
	Type        NotificationType
	Subject     string
	Message     string
}

type NotificationSender interface {
	Send(ctx context.Context, notification Notification) error
}

type NotificationService struct {
	emailProvider EmailProvider
	smsProvider   SMSProvider
}

func NewNotificationService(email EmailProvider, sms SMSProvider) *NotificationService {
	return &NotificationService{
		emailProvider: email,
		smsProvider:   sms,
	}
}

func (s *NotificationService) Send(ctx context.Context, n Notification) error {
	var err error

	switch n.Type {
	case NotificationTypeEmail:
		if s.emailProvider != nil {
			_, err = s.emailProvider.SendEmail(ctx, n.RecipientID, n.Subject, n.Message)
		} else {
			// Fallback logging if provider not configured
			fmt.Printf("[MOCK EMAIL] To: %s | Subject: %s\n", n.RecipientID, n.Subject)
		}
	case NotificationTypePush:
		// Push not yet implemented, specific provider needed
		fmt.Printf("[MOCK PUSH] To: %s | Subject: %s\n", n.RecipientID, n.Subject)
	case NotificationTypeSMS:
		// Assuming we add SMS type constant
		if s.smsProvider != nil {
			_, err = s.smsProvider.SendSMS(ctx, n.RecipientID, n.Message)
		} else {
			fmt.Printf("[MOCK SMS] To: %s | Body: %s\n", n.RecipientID, n.Message)
		}
	default:
		return fmt.Errorf("unsupported notification type: %s", n.Type)
	}

	return err
}

// Deprecated: ConsoleMock kept for backward compat if needed temporarily
type ConsoleNotificationSender struct{}

func NewConsoleNotificationSender() *ConsoleNotificationSender { return &ConsoleNotificationSender{} }
func (s *ConsoleNotificationSender) Send(ctx context.Context, n Notification) error {
	fmt.Printf("[CONSOLE FALLBACK] %s\n", n.Message)
	return nil
}
