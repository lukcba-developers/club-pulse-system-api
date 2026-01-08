package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/team/domain"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestTravelEvent_CalculateCostPerPerson(t *testing.T) {
	tests := []struct {
		name           string
		actualCost     decimal.Decimal
		estimatedCost  decimal.Decimal
		confirmedCount int
		expectedResult decimal.Decimal
	}{
		{
			name:           "Con costo real y 4 confirmados",
			actualCost:     decimal.NewFromInt(10000),
			estimatedCost:  decimal.NewFromInt(8000),
			confirmedCount: 4,
			expectedResult: decimal.NewFromInt(2500), // 10000 / 4
		},
		{
			name:           "Sin costo real, usa estimado",
			actualCost:     decimal.Zero,
			estimatedCost:  decimal.NewFromInt(8000),
			confirmedCount: 4,
			expectedResult: decimal.NewFromInt(2000), // 8000 / 4
		},
		{
			name:           "Sin confirmados, retorna cero",
			actualCost:     decimal.NewFromInt(10000),
			estimatedCost:  decimal.NewFromInt(8000),
			confirmedCount: 0,
			expectedResult: decimal.Zero,
		},
		{
			name:           "Un solo confirmado",
			actualCost:     decimal.NewFromInt(5000),
			estimatedCost:  decimal.Zero,
			confirmedCount: 1,
			expectedResult: decimal.NewFromInt(5000),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &domain.TravelEvent{
				ID:            uuid.New(),
				ActualCost:    tt.actualCost,
				EstimatedCost: tt.estimatedCost,
			}

			result := event.CalculateCostPerPerson(tt.confirmedCount)
			assert.True(t, tt.expectedResult.Equal(result), "Expected %s, got %s", tt.expectedResult, result)
		})
	}
}

func TestTravelEvent_IsOpen(t *testing.T) {
	tests := []struct {
		name          string
		departureDate time.Time
		expected      bool
	}{
		{
			name:          "Fecha futura - Abierto",
			departureDate: time.Now().AddDate(0, 0, 7),
			expected:      true,
		},
		{
			name:          "Fecha pasada - Cerrado",
			departureDate: time.Now().AddDate(0, 0, -1),
			expected:      false,
		},
		{
			name:          "Mismo momento - Cerrado",
			departureDate: time.Now().Add(-1 * time.Minute),
			expected:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &domain.TravelEvent{
				ID:            uuid.New(),
				DepartureDate: tt.departureDate,
			}

			result := event.IsOpen()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTravelEvent_IsFull(t *testing.T) {
	maxParticipants := 15

	tests := []struct {
		name            string
		maxParticipants *int
		currentCount    int
		expected        bool
	}{
		{
			name:            "Sin límite de participantes",
			maxParticipants: nil,
			currentCount:    100,
			expected:        false,
		},
		{
			name:            "Por debajo del límite",
			maxParticipants: &maxParticipants,
			currentCount:    10,
			expected:        false,
		},
		{
			name:            "En el límite exacto",
			maxParticipants: &maxParticipants,
			currentCount:    15,
			expected:        true,
		},
		{
			name:            "Sobre el límite",
			maxParticipants: &maxParticipants,
			currentCount:    20,
			expected:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &domain.TravelEvent{
				ID:              uuid.New(),
				MaxParticipants: tt.maxParticipants,
			}

			result := event.IsFull(tt.currentCount)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEventTypes(t *testing.T) {
	assert.Equal(t, domain.EventType("TRAVEL"), domain.EventTypeTravel)
	assert.Equal(t, domain.EventType("MATCH"), domain.EventTypeMatch)
	assert.Equal(t, domain.EventType("TOURNAMENT"), domain.EventTypeTournament)
	assert.Equal(t, domain.EventType("TRAINING"), domain.EventTypeTraining)
}

func TestRSVPStatuses(t *testing.T) {
	assert.Equal(t, domain.RSVPStatus("PENDING"), domain.RSVPStatusPending)
	assert.Equal(t, domain.RSVPStatus("CONFIRMED"), domain.RSVPStatusConfirmed)
	assert.Equal(t, domain.RSVPStatus("DECLINED"), domain.RSVPStatusDeclined)
}
