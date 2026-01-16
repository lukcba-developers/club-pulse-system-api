package domain

import (
	"context"
	"time"
)

type EquipmentCondition string

const (
	EquipmentConditionExcellent EquipmentCondition = "excellent"
	EquipmentConditionGood      EquipmentCondition = "good"
	EquipmentConditionFair      EquipmentCondition = "fair"
	EquipmentConditionPoor      EquipmentCondition = "poor"
	EquipmentConditionDamaged   EquipmentCondition = "damaged"
)

type Equipment struct {
	ID          string             `json:"id"`
	FacilityID  string             `json:"facility_id"`
	Name        string             `json:"name"`
	Type        string             `json:"type"` // e.g., "Tennis Racket", "Net", "Gym Machine"
	Condition   EquipmentCondition `json:"condition"`
	Status      string             `json:"status"` // "available", "maintenance", "in_use"
	IsAvailable bool               `json:"is_available"`

	PurchaseDate *time.Time `json:"purchase_date,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EquipmentRepository interface {
	Create(ctx context.Context, clubID string, equipment *Equipment) error
	GetByID(ctx context.Context, clubID, id string) (*Equipment, error)
	ListByFacility(ctx context.Context, clubID, facilityID string) ([]*Equipment, error)
	Update(ctx context.Context, clubID string, equipment *Equipment) error
}
