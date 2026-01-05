package domain

import (
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
	HourlyRate     float64        `json:"hourly_rate"`
	OpeningHour    int            `json:"opening_hour"`   // 0-23
	ClosingHour    int            `json:"closing_hour"`   // 0-23
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
	Create(facility *Facility) error
	GetByID(clubID, id string) (*Facility, error)
	List(clubID string, limit, offset int) ([]*Facility, error)
	Update(facility *Facility) error

	// Maintenance Extensions
	HasConflict(clubID, facilityID string, startTime, endTime time.Time) (bool, error)

	// Semantic Search Extensions
	SemanticSearch(clubID string, embedding []float32, limit int) ([]*FacilityWithSimilarity, error)
	UpdateEmbedding(facilityID string, embedding []float32) error

	// Equipment Management
	CreateEquipment(equipment *Equipment) error
	GetEquipmentByID(id string) (*Equipment, error)
	ListEquipmentByFacility(facilityID string) ([]*Equipment, error)
	UpdateEquipment(equipment *Equipment) error
}
