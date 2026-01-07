package domain

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"` // "credit", "debit", "manual_debt", "cantina_charge"
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

type TransactionHistory []Transaction

type Wallet struct {
	ID           uuid.UUID          `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID       string             `json:"user_id" gorm:"type:varchar(100);not null;uniqueIndex"`
	Balance      float64            `json:"balance" gorm:"default:0.0"`
	Points       int                `json:"points" gorm:"default:0"`
	Transactions TransactionHistory `json:"transactions" gorm:"type:jsonb;serializer:json"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}
