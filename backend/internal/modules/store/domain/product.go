package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	// Status Constants
	ProductStatusActive   = "ACTIVE"
	ProductStatusInactive = "INACTIVE"
)

type Product struct {
	ID            uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ClubID        string          `json:"club_id" gorm:"not null;index"`
	Name          string          `json:"name" gorm:"not null"`
	Description   string          `json:"description,omitempty"`
	Price         decimal.Decimal `json:"price" gorm:"type:decimal(10,2);not null"`
	StockQuantity int             `json:"stock_quantity" gorm:"default:0"`
	SKU           string          `json:"sku,omitempty"`
	Category      string          `json:"category,omitempty"` // Merch, Buffet, Equipment
	IsActive      bool            `json:"is_active" gorm:"default:true"`
	Status        string          `json:"status" gorm:"-"` // Mapped from IsActive
	ImageURL      string          `json:"image_url,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	DeletedAt     *time.Time      `json:"deleted_at,omitempty" gorm:"index"`
}

// AfterFind is a GORM hook that populates the computed Status field
func (p *Product) AfterFind(tx *gorm.DB) (err error) {
	p.Status = ProductStatusInactive
	if p.IsActive {
		p.Status = ProductStatusActive
	}
	return
}

type Order struct {
	ID          uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ClubID      string          `json:"club_id" gorm:"not null;index"`
	UserID      *string         `json:"user_id" gorm:"index"`
	GuestName   string          `json:"guest_name,omitempty"`
	GuestEmail  string          `json:"guest_email,omitempty"`
	TotalAmount decimal.Decimal `json:"total_amount" gorm:"type:decimal(10,2);not null"`
	Status      string          `json:"status" gorm:"default:'PAID'"`     // PAID, PENDING, CANCELLED
	Items       datatypes.JSON  `json:"items" gorm:"type:jsonb;not null"` // [{product_id, qty, unit_price}]
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   *time.Time      `json:"deleted_at,omitempty"`
}

type OrderItem struct {
	ProductID uuid.UUID       `json:"product_id"`
	Quantity  int             `json:"quantity"`
	UnitPrice decimal.Decimal `json:"unit_price"`
}

type StoreRepository interface {
	CreateProduct(ctx context.Context, product *Product) error
	GetProduct(ctx context.Context, clubID, id string) (*Product, error)
	UpdateProduct(ctx context.Context, product *Product) error
	ListProducts(ctx context.Context, clubID string, category string) ([]Product, error)
	CreateOrder(ctx context.Context, order *Order) error
	CreateOrderWithStockUpdate(ctx context.Context, order *Order, items []OrderItem) error
	DecreaseStock(ctx context.Context, clubID, productID string, quantity int) error
}
