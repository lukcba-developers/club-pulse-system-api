package service

import "context"

type NotificationType string

const (
	NotificationTypeEmail NotificationType = "EMAIL"
	NotificationTypePush  NotificationType = "PUSH"
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

// ConsoleNotificationSender mocks sending notifications by logging to console
type ConsoleNotificationSender struct{}

func NewConsoleNotificationSender() *ConsoleNotificationSender {
	return &ConsoleNotificationSender{}
}

func (s *ConsoleNotificationSender) Send(ctx context.Context, n Notification) error {
	// In a real implementation this would call SendGrid/Firebase
	// fmt.Printf("[MOCK NOTIFICATION] To: %s | Via: %s | Subject: %s | Body: %s\n", n.RecipientID, n.Type, n.Subject, n.Message)
	return nil
}
