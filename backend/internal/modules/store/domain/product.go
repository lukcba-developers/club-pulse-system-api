package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Product struct {
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ClubID        string     `json:"club_id" gorm:"not null;index"`
	Name          string     `json:"name" gorm:"not null"`
	Description   string     `json:"description,omitempty"`
	Price         float64    `json:"price" gorm:"type:decimal(10,2);not null"`
	StockQuantity int        `json:"stock_quantity" gorm:"default:0"`
	SKU           string     `json:"sku,omitempty"`
	Category      string     `json:"category,omitempty"` // Merch, Buffet, Equipment
	IsActive      bool       `json:"is_active" gorm:"default:true"`
	ImageURL      string     `json:"image_url,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

type Order struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ClubID      string         `json:"club_id" gorm:"not null;index"`
	UserID      *string        `json:"user_id" gorm:"index"`
	GuestName   string         `json:"guest_name,omitempty"`
	GuestEmail  string         `json:"guest_email,omitempty"`
	TotalAmount float64        `json:"total_amount" gorm:"type:decimal(10,2);not null"`
	Status      string         `json:"status" gorm:"default:'PAID'"`     // PAID, PENDING, CANCELLED
	Items       datatypes.JSON `json:"items" gorm:"type:jsonb;not null"` // [{product_id, qty, unit_price}]
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   *time.Time     `json:"deleted_at,omitempty"`
}

type OrderItem struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
	UnitPrice float64   `json:"unit_price"`
}

type StoreRepository interface {
	CreateProduct(product *Product) error
	GetProduct(id string) (*Product, error)
	UpdateProduct(product *Product) error
	ListProducts(clubID string, category string) ([]Product, error)
	CreateOrder(order *Order) error
	DecreaseStock(productID string, quantity int) error
}
