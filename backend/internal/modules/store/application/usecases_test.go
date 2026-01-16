package application_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/store/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/store/domain"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockStoreRepo struct {
	mock.Mock
}

func (m *MockStoreRepo) CreateProduct(ctx context.Context, product *domain.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}
func (m *MockStoreRepo) GetProduct(ctx context.Context, clubID, id string) (*domain.Product, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}
func (m *MockStoreRepo) UpdateProduct(ctx context.Context, product *domain.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}
func (m *MockStoreRepo) ListProducts(ctx context.Context, clubID string, category string) ([]domain.Product, error) {
	args := m.Called(ctx, clubID, category)
	return args.Get(0).([]domain.Product), args.Error(1)
}
func (m *MockStoreRepo) CreateOrder(ctx context.Context, order *domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}
func (m *MockStoreRepo) CreateOrderWithStockUpdate(ctx context.Context, order *domain.Order, items []domain.OrderItem) error {
	args := m.Called(ctx, order, items)
	return args.Error(0)
}
func (m *MockStoreRepo) DecreaseStock(ctx context.Context, clubID, productID string, quantity int) error {
	args := m.Called(ctx, clubID, productID, quantity)
	return args.Error(0)
}

// --- Tests ---

func TestStoreUseCases_GetCatalog(t *testing.T) {
	repo := new(MockStoreRepo)
	uc := application.NewStoreUseCases(repo)
	clubID := "c1"

	t.Run("Success", func(t *testing.T) {
		products := []domain.Product{{Name: "P1", Price: decimal.NewFromInt(10)}}
		repo.On("ListProducts", mock.Anything, clubID, "MERCH").Return(products, nil).Once()

		res, err := uc.GetCatalog(context.TODO(), clubID, "MERCH")
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		repo.AssertExpectations(t)
	})
}

func TestStoreUseCases_PurchaseItems(t *testing.T) {
	repo := new(MockStoreRepo)
	uc := application.NewStoreUseCases(repo)
	clubID := "c1"
	pID := uuid.New()
	userID := "u1"

	t.Run("Success: Purchase", func(t *testing.T) {
		req := application.PurchaseRequest{
			ClubID: clubID,
			UserID: &userID,
			Items:  []domain.OrderItem{{ProductID: pID, Quantity: 2}},
		}

		repo.On("GetProduct", mock.Anything, clubID, pID.String()).Return(&domain.Product{
			ID: pID, Name: "Soda", Price: decimal.NewFromInt(5), StockQuantity: 10, IsActive: true,
		}, nil).Once()

		repo.On("CreateOrderWithStockUpdate", mock.Anything, mock.MatchedBy(func(o *domain.Order) bool {
			// 2 * 5 = 10
			return o.TotalAmount.Equal(decimal.NewFromInt(10)) && o.Status == "PAID"
		}), mock.Anything).Return(nil).Once()

		order, err := uc.PurchaseItems(context.TODO(), req)
		assert.NoError(t, err)
		assert.NotNil(t, order)
		assert.True(t, order.TotalAmount.Equal(decimal.NewFromInt(10)))
		repo.AssertExpectations(t)
	})

	t.Run("Fail: Out of Stock", func(t *testing.T) {
		req := application.PurchaseRequest{
			ClubID: clubID,
			UserID: &userID,
			Items:  []domain.OrderItem{{ProductID: pID, Quantity: 20}},
		}

		repo.On("GetProduct", mock.Anything, clubID, pID.String()).Return(&domain.Product{
			ID: pID, Name: "Soda", Price: decimal.NewFromInt(5), StockQuantity: 10, IsActive: true,
		}, nil).Once()

		_, err := uc.PurchaseItems(context.TODO(), req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insufficient stock")
	})

	t.Run("Fail: Product Inactive", func(t *testing.T) {
		req := application.PurchaseRequest{
			ClubID: clubID,
			UserID: &userID,
			Items:  []domain.OrderItem{{ProductID: pID, Quantity: 1}},
		}

		repo.On("GetProduct", mock.Anything, clubID, pID.String()).Return(&domain.Product{
			ID: pID, Name: "Soda", Price: decimal.NewFromInt(5), StockQuantity: 10, IsActive: false,
		}, nil).Once()

		_, err := uc.PurchaseItems(context.TODO(), req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not active")
	})
}
