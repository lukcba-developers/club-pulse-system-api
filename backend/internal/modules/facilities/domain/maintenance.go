package domain

import (
	"time"
)

type MaintenanceStatus string

const (
	MaintenanceStatusScheduled  MaintenanceStatus = "scheduled"
	MaintenanceStatusInProgress MaintenanceStatus = "in_progress"
	MaintenanceStatusCompleted  MaintenanceStatus = "completed"
	MaintenanceStatusCancelled  MaintenanceStatus = "cancelled"
)

type MaintenanceType string

const (
	MaintenanceTypePreventive MaintenanceType = "preventive"
	MaintenanceTypeCorrective MaintenanceType = "corrective"
	MaintenanceTypeCleaning   MaintenanceType = "cleaning"
)

type MaintenanceTask struct {
	ID          string            `json:"id"`
	FacilityID  string            `json:"facility_id"`
	EquipmentID *string           `json:"equipment_id,omitempty"` // Nullable for facility-only maintenance
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      MaintenanceStatus `json:"status"`
	Type        MaintenanceType   `json:"type"`

	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`

	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Repository Interface extension for Maintenance
type MaintenanceRepository interface {
	Create(task *MaintenanceTask) error
	GetByID(id string) (*MaintenanceTask, error)
	ListByFacility(facilityID string) ([]*MaintenanceTask, error)
	// HasConflict checks if there is any active maintenance overlapping with the given time range
	HasConflict(facilityID string, startTime, endTime time.Time) (bool, error)
	// GetImpactedUsers returns a list of user IDs that have bookings during the maintenance window
	GetImpactedUsers(facilityID string, start, end time.Time) ([]string, error)
}
