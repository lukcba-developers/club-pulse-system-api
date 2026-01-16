package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/notification/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockEmailProvider struct {
	mock.Mock
}

func (m *MockEmailProvider) SendEmail(ctx context.Context, to, subject, body string) (*service.DeliveryResult, error) {
	args := m.Called(ctx, to, subject, body)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.DeliveryResult), args.Error(1)
}

type MockSMSProvider struct {
	mock.Mock
}

func (m *MockSMSProvider) SendSMS(ctx context.Context, to, body string) (*service.DeliveryResult, error) {
	args := m.Called(ctx, to, body)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.DeliveryResult), args.Error(1)
}

// --- Tests ---

func TestNotificationService_Send(t *testing.T) {
	t.Run("Send Email Success", func(t *testing.T) {
		emailMock := new(MockEmailProvider)
		smsMock := new(MockSMSProvider)
		svc := service.NewNotificationService(emailMock, smsMock)

		emailMock.On("SendEmail", mock.Anything, "user@example.com", "Subject", "Body").
			Return(&service.DeliveryResult{Success: true}, nil).Once()

		err := svc.Send(context.Background(), service.Notification{
			RecipientID: "user@example.com",
			Type:        service.NotificationTypeEmail,
			Title:       "Subject",
			Body:        "Body",
		})
		assert.NoError(t, err)
		emailMock.AssertExpectations(t)
	})

	t.Run("Send SMS Success", func(t *testing.T) {
		emailMock := new(MockEmailProvider)
		smsMock := new(MockSMSProvider)
		svc := service.NewNotificationService(emailMock, smsMock)

		smsMock.On("SendSMS", mock.Anything, "+123456789", "SMS Body").
			Return(&service.DeliveryResult{Success: true}, nil).Once()

		err := svc.Send(context.Background(), service.Notification{
			RecipientID: "+123456789",
			Type:        service.NotificationTypeSMS,
			Body:        "SMS Body",
		})
		assert.NoError(t, err)
		smsMock.AssertExpectations(t)
	})

	t.Run("Send Email Error", func(t *testing.T) {
		emailMock := new(MockEmailProvider)
		svc := service.NewNotificationService(emailMock, nil)

		emailMock.On("SendEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("smtp error")).Once()

		err := svc.Send(context.Background(), service.Notification{
			RecipientID: "fail@example.com",
			Type:        service.NotificationTypeEmail,
			Title:       "Sub",
			Body:        "Msg",
		})
		assert.Error(t, err)
		assert.Equal(t, "smtp error", err.Error())
	})

	t.Run("Fallback to Mock if Provider Nil", func(t *testing.T) {
		svc := service.NewNotificationService(nil, nil)

		// Should just log to stdout and return nil error
		err := svc.Send(context.Background(), service.Notification{
			RecipientID: "mock@example.com",
			Type:        service.NotificationTypeEmail,
			Title:       "Mock Subject",
			Body:        "Mock Body",
		})
		assert.NoError(t, err)
	})

	t.Run("Unsupported Type", func(t *testing.T) {
		svc := service.NewNotificationService(nil, nil)
		err := svc.Send(context.Background(), service.Notification{
			Type: "UNKNOWN",
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported notification type")
	})

	t.Run("Backward Compatibility Mapping", func(t *testing.T) {
		emailMock := new(MockEmailProvider)
		svc := service.NewNotificationService(emailMock, nil)

		emailMock.On("SendEmail", mock.Anything, "old@example.com", "Old Subject", "Old Message").
			Return(&service.DeliveryResult{Success: true}, nil).Once()

		err := svc.Send(context.Background(), service.Notification{
			RecipientID: "old@example.com",
			Type:        service.NotificationTypeEmail,
			Subject:     "Old Subject", // Deprecated field
			Message:     "Old Message", // Deprecated field
		})
		assert.NoError(t, err)
		emailMock.AssertExpectations(t)
	})
}
