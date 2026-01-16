package domain

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Enums

type FacilityType string

const (
	FacilityTypeCourt FacilityType = "court"
	FacilityTypePool  FacilityType = "pool"
	FacilityTypeGym   FacilityType = "gym"
	FacilityTypeField FacilityType = "field"
)

type FacilityStatus string

const (
	FacilityStatusActive      FacilityStatus = "active"
	FacilityStatusMaintenance FacilityStatus = "maintenance"
	FacilityStatusClosed      FacilityStatus = "closed"
)

// JSONB Structures

type Specifications struct {
	SurfaceType *string  `json:"surface_type,omitempty"`
	Lighting    bool     `json:"lighting"`
	Covered     bool     `json:"covered"`
	Equipment   []string `json:"equipment,omitempty"` // Basic inventory list
}

// Value method for GORM storage
func (s Specifications) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan method for GORM storage
func (s *Specifications) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, s)
}

type Location struct {
	Name        string `json:"name"` // e.g. "Main Building"
	Description string `json:"description,omitempty"`
}

func (l Location) Value() (driver.Value, error) {
	return json.Marshal(l)
}

func (l *Location) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, l)
}

// Main Entity

type Facility struct {
	ID             string         `json:"id"`
	ClubID         string         `json:"club_id" gorm:"index;not null"`
	Name           string         `json:"name"`
	Type           FacilityType   `json:"type"`
	Status         FacilityStatus `json:"status"`
	Capacity       int            `json:"capacity"`
	Description    string         `json:"description"`
	HourlyRate     float64        `json:"hourly_rate"`
	OpeningTime    string         `json:"opening_time"` // HH:MM
	ClosingTime    string         `json:"closing_time"` // HH:MM
	GuestFee       float64        `json:"guest_fee"`
	Specifications Specifications `json:"specifications"` // Stored as JSONB
	Location       Location       `json:"location"`       // Stored as JSONB

	// Semantic Search (pgvector)
	Embedding []float32 `json:"-" gorm:"-"` // Managed by 002_pgvector_indexes.sql

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FacilityWithSimilarity represents a facility with its search similarity score
type FacilityWithSimilarity struct {
	Facility   *Facility
	Similarity float32
}

// Repository Interface

type FacilityRepository interface {
	Create(ctx context.Context, facility *Facility) error
	GetByID(ctx context.Context, clubID, id string) (*Facility, error)
	GetByIDForUpdate(ctx context.Context, clubID, id string) (*Facility, error)
	List(ctx context.Context, clubID string, limit, offset int) ([]*Facility, error)
	Update(ctx context.Context, facility *Facility) error

	// Maintenance Extensions
	HasConflict(ctx context.Context, clubID, facilityID string, startTime, endTime time.Time) (bool, error)
	ListMaintenanceByFacility(ctx context.Context, clubID, facilityID string) ([]*MaintenanceTask, error)

	// Semantic Search Extensions
	SemanticSearch(ctx context.Context, clubID string, embedding []float32, limit int) ([]*FacilityWithSimilarity, error)
	UpdateEmbedding(ctx context.Context, facilityID string, embedding []float32) error

	// Equipment Management
	CreateEquipment(ctx context.Context, clubID string, equipment *Equipment) error
	GetEquipmentByID(ctx context.Context, clubID, id string) (*Equipment, error)
	ListEquipmentByFacility(ctx context.Context, clubID, facilityID string) ([]*Equipment, error)
	UpdateEquipment(ctx context.Context, clubID string, equipment *Equipment) error
	LoanEquipmentAtomic(ctx context.Context, loan *EquipmentLoan, equipmentID string) error
}
