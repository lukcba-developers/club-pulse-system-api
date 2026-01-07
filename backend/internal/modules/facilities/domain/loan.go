package domain

import (
	"time"
)

type LoanStatus string

const (
	LoanStatusActive   LoanStatus = "ACTIVE"
	LoanStatusReturned LoanStatus = "RETURNED"
	LoanStatusOverdue  LoanStatus = "OVERDUE"
	LoanStatusLost     LoanStatus = "LOST"
)

type EquipmentLoan struct {
	ID                string     `json:"id"`
	EquipmentID       string     `json:"equipment_id"`
	UserID            string     `json:"user_id"`
	LoanedAt          time.Time  `json:"loaned_at"`
	ExpectedReturnAt  time.Time  `json:"expected_return_at"`
	ReturnedAt        *time.Time `json:"returned_at,omitempty"`
	Status            LoanStatus `json:"status"`
	ConditionOnReturn string     `json:"condition_on_return,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// LoanDisplay is a DTO for displaying loan details with related entity names
type LoanDisplay struct {
	ID               string     `json:"id"`
	EquipmentName    string     `json:"equipment_name"`
	UserName         string     `json:"user_name"`
	LoanedAt         time.Time  `json:"loaned_at"`
	ExpectedReturnAt time.Time  `json:"expected_return_at"`
	Status           LoanStatus `json:"status"`
}

type LoanRepository interface {
	Create(loan *EquipmentLoan) error
	GetByID(id string) (*EquipmentLoan, error)
	ListByUser(userID string) ([]*EquipmentLoan, error)
	ListByStatus(status LoanStatus) ([]*EquipmentLoan, error)
	Update(loan *EquipmentLoan) error
}
