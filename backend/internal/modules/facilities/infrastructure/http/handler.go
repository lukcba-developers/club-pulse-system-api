package http

import (
	"net/http"
	"strconv"

	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/application"
	platformRedis "github.com/lukcba/club-pulse-system-api/backend/internal/platform/redis"
)

type FacilityHandler struct {
	useCases *application.FacilityUseCases
}

func NewFacilityHandler(useCases *application.FacilityUseCases) *FacilityHandler {
	return &FacilityHandler{
		useCases: useCases,
	}
}

func RegisterRoutes(r *gin.RouterGroup, handler *FacilityHandler, authMiddleware, tenantMiddleware gin.HandlerFunc) {
	facilities := r.Group("/facilities")
	// Public routes? Or protected? Let's protect create, allow public read potentially.
	// For MVP, protect all write ops.

	protected := facilities.Group("")
	protected.Use(authMiddleware, tenantMiddleware)
	{
		protected.GET("", handler.List)
		protected.GET("/:id", handler.Get)
		protected.POST("", handler.Create)
		protected.PUT("/:id", handler.Update)
	}
}

func (h *FacilityHandler) Create(c *gin.Context) {
	var dto application.CreateFacilityDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clubID := c.GetString("clubID")
	facility, err := h.useCases.CreateFacility(clubID, dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.invalidateCache(c, clubID)
	c.JSON(http.StatusCreated, facility)
}

// List godoc
// @Summary      List facilities
// @Description  Get a list of facilities for the authenticated club
// @Tags         facilities
// @Accept       json
// @Produce      json
// @Param        limit   query      int  false  "Limit"
// @Param        offset  query      int  false  "Offset"
// @Success      200     {array}    domain.Facility
// @Failure      500     {object}   map[string]string
// @Router       /facilities [get]
func (h *FacilityHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	clubID := c.GetString("clubID")

	// Cache Key
	cacheKey := fmt.Sprintf("facilities:list:%s:%d:%d", clubID, limit, offset)
	ctx := c.Request.Context()

	// 1. Try Cache
	cached, err := platformRedis.GetClient().Get(ctx, cacheKey)
	if err == nil && cached != "" {
		c.Header("X-Cache", "HIT")
		c.Header("Content-Type", "application/json")
		c.String(http.StatusOK, cached)
		return
	}

	// 2. DB Lookups
	facilities, err := h.useCases.ListFacilities(clubID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. Set Cache
	data, _ := json.Marshal(facilities)
	_ = platformRedis.GetClient().Set(ctx, cacheKey, string(data), 5*time.Minute)

	c.Header("X-Cache", "MISS")
	c.JSON(http.StatusOK, facilities)
}

func (h *FacilityHandler) Get(c *gin.Context) {
	id := c.Param("id")
	clubID := c.GetString("clubID")
	facility, err := h.useCases.GetFacility(clubID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if facility == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Facility not found"})
		return
	}
	c.JSON(http.StatusOK, facility)
}

func (h *FacilityHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var dto application.UpdateFacilityDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clubID := c.GetString("clubID")
	facility, err := h.useCases.UpdateFacility(clubID, id, dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.invalidateCache(c, clubID)
	c.JSON(http.StatusOK, facility)
}

func (h *FacilityHandler) invalidateCache(c *gin.Context, clubID string) {
	ctx := c.Request.Context()
	pattern := fmt.Sprintf("facilities:list:%s:*", clubID)
	keys, _ := platformRedis.GetClient().Scan(ctx, pattern)
	if len(keys) > 0 {
		_ = platformRedis.GetClient().Del(ctx, keys...)
	}
}
