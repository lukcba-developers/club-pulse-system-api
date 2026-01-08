package repository

import (
	"errors"

	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/store/domain"
	"gorm.io/gorm"
)

type PostgresStoreRepository struct {
	db *gorm.DB
}

func NewPostgresStoreRepository(db *gorm.DB) *PostgresStoreRepository {
	return &PostgresStoreRepository{db: db}
}

func (r *PostgresStoreRepository) CreateProduct(product *domain.Product) error {
	return r.db.Create(product).Error
}

func (r *PostgresStoreRepository) GetProduct(id string) (*domain.Product, error) {
	var product domain.Product
	if err := r.db.Where("id = ?", id).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *PostgresStoreRepository) UpdateProduct(product *domain.Product) error {
	return r.db.Save(product).Error
}

func (r *PostgresStoreRepository) ListProducts(clubID string, category string) ([]domain.Product, error) {
	var products []domain.Product
	query := r.db.Where("club_id = ? AND is_active = ?", clubID, true)
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *PostgresStoreRepository) CreateOrder(order *domain.Order) error {
	return r.db.Create(order).Error
}

func (r *PostgresStoreRepository) CreateOrderWithStockUpdate(order *domain.Order, items []domain.OrderItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Decrement Stock
		for _, item := range items {
			result := tx.Model(&domain.Product{}).
				Where("id = ? AND stock_quantity >= ?", item.ProductID, item.Quantity).
				UpdateColumn("stock_quantity", gorm.Expr("stock_quantity - ?", item.Quantity))

			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return errors.New("insufficient stock for product: " + item.ProductID.String())
			}
		}

		// 2. Create Order
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		return nil
	})
}

// DecreaseStock decrements the stock of a product by quantity.
// It ensures stock doesn't go below zero (optional business rule).
func (r *PostgresStoreRepository) DecreaseStock(productID string, quantity int) error {
	result := r.db.Model(&domain.Product{}).
		Where("id = ? AND stock_quantity >= ?", productID, quantity).
		UpdateColumn("stock_quantity", gorm.Expr("stock_quantity - ?", quantity))

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("insufficient stock or product not found")
	}
	return nil
}
