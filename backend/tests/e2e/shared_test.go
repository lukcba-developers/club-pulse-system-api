package e2e_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/notification/service"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/domain"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/database"
	"gorm.io/gorm"
)

// SharedMockNotifier is a generic mock for the notification service and refund service.
type SharedMockNotifier struct{}

func (m *SharedMockNotifier) Send(ctx context.Context, n service.Notification) error {
	return nil
}

func (m *SharedMockNotifier) Refund(ctx context.Context, clubID string, referenceID uuid.UUID, referenceType string) error {
	return nil
}

// ensure we satisfy the interfaces
var _ service.NotificationSender = (*SharedMockNotifier)(nil)

// RecordingMockPaymentGateway records calls for verification
type RecordingMockPaymentGateway struct {
	RefundCalledWith []string
}

func (m *RecordingMockPaymentGateway) CreatePreference(ctx context.Context, p *domain.Payment, e, d string) (string, error) {
	return "http://mock.mp/pref", nil
}

func (m *RecordingMockPaymentGateway) ProcessWebhook(ctx context.Context, pl interface{}) (*domain.Payment, error) {
	return &domain.Payment{Status: domain.PaymentStatusCompleted, ExternalID: "ext-123"}, nil
}

func (m *RecordingMockPaymentGateway) ValidateWebhook(req *http.Request) error {
	return nil
}

func (m *RecordingMockPaymentGateway) Refund(ctx context.Context, externalID string) error {
	m.RefundCalledWith = append(m.RefundCalledWith, externalID)
	return nil
}

// SetupTestDB returns a transactional DB and a cleanup function.
// This ensures each test runs in isolation and rolls back changes.
func SetupTestDB(t *testing.T) *gorm.DB {
	database.InitDB()
	db := database.GetDB()

	tx := db.Begin()
	t.Cleanup(func() {
		tx.Rollback()
	})

	return tx
}
