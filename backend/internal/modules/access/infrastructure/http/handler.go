package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/access/application"
)

type AccessHandler struct {
	useCases *application.AccessUseCases
}

func NewAccessHandler(useCases *application.AccessUseCases) *AccessHandler {
	return &AccessHandler{useCases: useCases}
}

func (h *AccessHandler) ValidateEntry(c *gin.Context) {
	var req application.EntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log, err := h.useCases.RequestEntry(c.Request.Context(), req)
	if err != nil {
		// Even if error (e.g. denied), we might return 200 with Denied status,
		// but RequestEntry returns error only on system failure or logging failure usually.
		// If it returns a log with DENIED status, it's a success in execution.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Explicitly check status to return 403 if Denied?
	// Usually access control devices expect 200 OK + payload saying "Granted/Denied"
	// or 403 Forbidden. Let's send 200 with data.

	statusCode := http.StatusOK
	if log.Status == "DENIED" {
		statusCode = http.StatusForbidden
	}

	c.JSON(statusCode, gin.H{"data": log})
}

func RegisterRoutes(r *gin.RouterGroup, handler *AccessHandler, authMiddleware gin.HandlerFunc) {
	access := r.Group("/access")
	access.Use(authMiddleware)
	{
		access.POST("/entry", handler.ValidateEntry)
	}
}
