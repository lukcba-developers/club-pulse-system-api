package e2e_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	storeApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/store/application"
	storeDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/store/domain"
	storeHttp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/store/infrastructure/http"
	storeRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/store/infrastructure/repository"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreFlow(t *testing.T) {
	// 1. Setup
	gin.SetMode(gin.TestMode)
	database.InitDB()
	db := database.GetDB()

	// Clean state
	_ = db.Migrator().DropTable(&storeDomain.Product{}, &storeDomain.Order{}, &storeDomain.OrderItem{})
	_ = db.AutoMigrate(&storeDomain.Product{}, &storeDomain.Order{}, &storeDomain.OrderItem{})

	repo := storeRepo.NewPostgresStoreRepository(db)
	uc := storeApp.NewStoreUseCases(repo)
	h := storeHttp.NewStoreHandler(uc)

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("userID", "store-user-id")
		c.Set("clubID", "test-club-store")
		c.Next()
	})

	storeHttp.RegisterRoutes(r.Group("/api/v1"), h, func(c *gin.Context) { c.Next() }, func(c *gin.Context) { c.Next() })

	// 2. Data Setup
	prodID := uuid.New()
	db.Create(&storeDomain.Product{
		ID:          prodID,
		ClubID:      "test-club-store",
		Name:        "Protein Bar",
		Description: "Chocolate",
		Price:       250.0,
		Stock:       100,
		Category:    "Food",
	})

	// 3. Test: List Products
	t.Run("List Products", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/store/products", nil)
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Protein Bar")
	})

	// 4. Test: Purchase
	t.Run("Purchase Product", func(t *testing.T) {
		// Payload structure depends on DTO. Assuming standard keys.
		body := `{
			"items": [
				{
					"product_id": "` + prodID.String() + `",
					"quantity": 2
				}
			]
		}`
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/store/purchase", strings.NewReader(body))
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)

		// Verify Stock Deduction
		var p storeDomain.Product
		db.First(&p, "id = ?", prodID)
		assert.Equal(t, 98, p.Stock)
	})
}
