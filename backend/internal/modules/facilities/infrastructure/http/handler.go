package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/application"
)

type FacilityHandler struct {
	useCases *application.FacilityUseCases
}

func NewFacilityHandler(useCases *application.FacilityUseCases) *FacilityHandler {
	return &FacilityHandler{
		useCases: useCases,
	}
}

func RegisterRoutes(r *gin.RouterGroup, handler *FacilityHandler, authMiddleware gin.HandlerFunc) {
	facilities := r.Group("/facilities")
	// Public routes? Or protected? Let's protect create, allow public read potentially.
	// For MVP, protect all write ops.

	facilities.GET("", handler.List)
	facilities.GET("/:id", handler.Get)

	protected := facilities.Group("")
	protected.Use(authMiddleware)
	{
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

	facility, err := h.useCases.CreateFacility(dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, facility)
}

func (h *FacilityHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	facilities, err := h.useCases.ListFacilities(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, facilities)
}

func (h *FacilityHandler) Get(c *gin.Context) {
	id := c.Param("id")
	facility, err := h.useCases.GetFacility(id)
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

	facility, err := h.useCases.UpdateFacility(id, dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, facility)
}
