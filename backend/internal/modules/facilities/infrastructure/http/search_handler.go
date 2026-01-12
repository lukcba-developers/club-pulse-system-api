package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/application"
)

// SearchHandler handles semantic search for facilities
type SearchHandler struct {
	searchUseCase *application.SemanticSearchUseCase
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(searchUseCase *application.SemanticSearchUseCase) *SearchHandler {
	return &SearchHandler{
		searchUseCase: searchUseCase,
	}
}

// RegisterSearchRoutes registers the semantic search routes
func RegisterSearchRoutes(r *gin.RouterGroup, handler *SearchHandler) {
	facilities := r.Group("/facilities")
	{
		// GET /api/v1/facilities/search?q=canchas+techadas&limit=10
		facilities.GET("/search", handler.Search)

		// POST /api/v1/facilities/embeddings/generate - Admin only, generates embeddings for all facilities
		facilities.POST("/embeddings/generate", handler.GenerateEmbeddings)
	}
}

// Search performs semantic search on facilities
// @Summary Search facilities using natural language
// @Description Search for facilities using semantic similarity (e.g., "canchas para lluvia", "piscina nocturna")
// @Tags Facilities
// @Produce json
// @Param q query string true "Search query in natural language"
// @Param limit query int false "Maximum results (default 10)"
// @Success 200 {array} application.SemanticSearchResult
// @Router /facilities/search [get]
func (h *SearchHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing search query",
			"message": "Please provide a search query using the 'q' parameter",
			"example": "/facilities/search?q=canchas+techadas+para+lluvia",
		})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	clubID := c.GetString("clubID")
	results, err := h.searchUseCase.Search(c.Request.Context(), clubID, query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Search failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"query":   query,
		"count":   len(results),
		"results": results,
	})
}

// GenerateEmbeddings generates embeddings for all facilities (admin operation)
// @Summary Generate embeddings for all facilities
// @Description Batch operation to generate and store embeddings for all facilities
// @Tags Facilities
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /facilities/embeddings/generate [post]
func (h *SearchHandler) GenerateEmbeddings(c *gin.Context) {
	clubID := c.GetString("clubID")
	count, err := h.searchUseCase.GenerateAllEmbeddings(c.Request.Context(), clubID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate embeddings",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Embeddings generated successfully",
		"processed": count,
	})
}
