package application

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/store/domain"
	"gorm.io/datatypes"
)

type StoreUseCases struct {
	repo domain.StoreRepository
}

func NewStoreUseCases(repo domain.StoreRepository) *StoreUseCases {
	return &StoreUseCases{repo: repo}
}

func (uc *StoreUseCases) GetCatalog(ctx context.Context, clubID string, category string) ([]domain.Product, error) {
	return uc.repo.ListProducts(ctx, clubID, category)
}

type PurchaseRequest struct {
	ClubID     string             `json:"club_id"`
	UserID     *string            `json:"user_id"`
	GuestName  string             `json:"guest_name"`
	GuestEmail string             `json:"guest_email"`
	Items      []domain.OrderItem `json:"items"`
}

func (uc *StoreUseCases) PurchaseItems(ctx context.Context, req PurchaseRequest) (*domain.Order, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("cannot purchase empty cart")
	}

	// Validation: Guest or User
	if req.UserID == nil && req.GuestEmail == "" {
		return nil, errors.New("must provide UserID or GuestEmail")
	}

	var totalAmount float64
	var orderItems []domain.OrderItem

	// Validate stock and calculate total (Optimistic check, real deduction happens in repo transaction)
	// In a high concurrence scenario, repo should handle "UPDATE ... WHERE stock > qty"
	for _, item := range req.Items {
		product, err := uc.repo.GetProduct(ctx, item.ProductID.String())
		if err != nil {
			return nil, errors.New("product not found: " + item.ProductID.String())
		}
		if product.StockQuantity < item.Quantity {
			return nil, errors.New("insufficient stock for product: " + product.Name)
		}
		if !product.IsActive {
			return nil, errors.New("product is not active: " + product.Name)
		}

		item.UnitPrice = product.Price
		totalAmount += item.UnitPrice * float64(item.Quantity)
		orderItems = append(orderItems, item)
	}

	itemsJSON, err := json.Marshal(orderItems)
	if err != nil {
		return nil, err
	}

	order := &domain.Order{
		ID:          uuid.New(),
		ClubID:      req.ClubID,
		UserID:      req.UserID,
		GuestName:   req.GuestName,
		GuestEmail:  req.GuestEmail,
		TotalAmount: totalAmount,
		Status:      "PAID", // Assuming instant payment for this phase
		Items:       datatypes.JSON(itemsJSON),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Atomic Purchase
	if err := uc.repo.CreateOrderWithStockUpdate(ctx, order, orderItems); err != nil {
		return nil, err
	}

	return order, nil
}
