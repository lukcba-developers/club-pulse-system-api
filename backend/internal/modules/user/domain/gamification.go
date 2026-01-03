package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"` // "credit", "debit"
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

type TransactionHistory []Transaction

func (t TransactionHistory) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *TransactionHistory) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, t)
}

type Wallet struct {
	ID           uuid.UUID          `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID       string             `json:"user_id" gorm:"type:varchar(100);not null;uniqueIndex"`
	Balance      float64            `json:"balance" gorm:"default:0.0"`
	Points       int                `json:"points" gorm:"default:0"`
	Transactions TransactionHistory `json:"transactions" gorm:"type:jsonb"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}
