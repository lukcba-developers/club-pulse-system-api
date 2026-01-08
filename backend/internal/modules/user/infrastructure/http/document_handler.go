package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
)

// DocumentHandler maneja las peticiones HTTP relacionadas con documentos de usuario
type DocumentHandler struct {
	docRepo            domain.UserDocumentRepository
	eligibilityService *application.EligibilityService
}

// NewDocumentHandler crea una nueva instancia del handler
func NewDocumentHandler(docRepo domain.UserDocumentRepository, eligibilityService *application.EligibilityService) *DocumentHandler {
	return &DocumentHandler{
		docRepo:            docRepo,
		eligibilityService: eligibilityService,
	}
}

// UploadDocumentRequest representa la solicitud de upload de documento
type UploadDocumentRequest struct {
	Type           string  `form:"type" binding:"required"`
	ExpirationDate *string `form:"expiration_date"` // Formato: YYYY-MM-DD
}

// UploadDocument maneja el upload de un documento
// POST /users/:userId/documents
func (h *DocumentHandler) UploadDocument(c *gin.Context) {
	clubID := c.GetString("club_id")
	userID := c.Param("userId")
	currentUserID := c.GetString("user_id")

	// Verificar que el usuario solo puede subir sus propios documentos (a menos que sea admin)
	role := c.GetString("role")
	if userID != currentUserID && role != "ADMIN" && role != "SUPER_ADMIN" {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permiso para subir documentos de otro usuario"})
		return
	}

	var req UploadDocumentRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener el archivo
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Archivo requerido"})
		return
	}

	// TODO: Subir archivo a MinIO/S3
	// Por ahora, guardamos la ruta local como placeholder
	fileURL := "/uploads/documents/" + uuid.New().String() + "_" + file.Filename

	// Parsear fecha de vencimiento si existe
	var expirationDate *time.Time
	if req.ExpirationDate != nil && *req.ExpirationDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.ExpirationDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido. Use YYYY-MM-DD"})
			return
		}
		expirationDate = &parsed
	}

	// Crear documento
	doc := &domain.UserDocument{
		ClubID:         clubID,
		UserID:         userID,
		Type:           domain.DocumentType(req.Type),
		FileURL:        fileURL,
		Status:         domain.DocumentStatusPending,
		ExpirationDate: expirationDate,
	}

	if err := h.docRepo.Create(doc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar documento"})
		return
	}

	c.JSON(http.StatusCreated, doc)
}

// ListDocuments lista todos los documentos de un usuario
// GET /users/:userId/documents
func (h *DocumentHandler) ListDocuments(c *gin.Context) {
	clubID := c.GetString("club_id")
	userID := c.Param("userId")

	docs, err := h.docRepo.GetByUserID(clubID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener documentos"})
		return
	}

	c.JSON(http.StatusOK, docs)
}

// GetDocument obtiene un documento específico
// GET /users/:userId/documents/:docId
func (h *DocumentHandler) GetDocument(c *gin.Context) {
	clubID := c.GetString("club_id")
	docID := c.Param("docId")

	docUUID, err := uuid.Parse(docID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de documento inválido"})
		return
	}

	doc, err := h.docRepo.GetByID(clubID, docUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Documento no encontrado"})
		return
	}

	c.JSON(http.StatusOK, doc)
}

// ValidateDocumentRequest representa la solicitud de validación
type ValidateDocumentRequest struct {
	Approve bool   `json:"approve" binding:"required"`
	Notes   string `json:"notes"`
}

// ValidateDocument valida o rechaza un documento
// PUT /users/:userId/documents/:docId/validate
func (h *DocumentHandler) ValidateDocument(c *gin.Context) {
	clubID := c.GetString("club_id")
	docID := c.Param("docId")
	validatorID := c.GetString("user_id")

	// Solo admins pueden validar documentos
	role := c.GetString("role")
	if role != "ADMIN" && role != "SUPER_ADMIN" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Solo administradores pueden validar documentos"})
		return
	}

	var req ValidateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.eligibilityService.ValidateDocument(clubID, docID, validatorID, req.Approve, req.Notes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener el documento actualizado
	docUUID, _ := uuid.Parse(docID)
	doc, _ := h.docRepo.GetByID(clubID, docUUID)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Documento validado exitosamente",
		"document": doc,
	})
}

// DeleteDocument elimina un documento
// DELETE /users/:userId/documents/:docId
func (h *DocumentHandler) DeleteDocument(c *gin.Context) {
	clubID := c.GetString("club_id")
	userID := c.Param("userId")
	docID := c.Param("docId")
	currentUserID := c.GetString("user_id")

	// Verificar permisos
	role := c.GetString("role")
	if userID != currentUserID && role != "ADMIN" && role != "SUPER_ADMIN" {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permiso para eliminar documentos de otro usuario"})
		return
	}

	docUUID, err := uuid.Parse(docID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de documento inválido"})
		return
	}

	if err := h.docRepo.Delete(clubID, docUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar documento"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Documento eliminado exitosamente"})
}

// CheckEligibility verifica la elegibilidad de un usuario
// GET /users/:userId/eligibility
func (h *DocumentHandler) CheckEligibility(c *gin.Context) {
	clubID := c.GetString("club_id")
	userID := c.Param("userId")

	result, err := h.eligibilityService.CheckEligibility(clubID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar elegibilidad"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetDocumentSummary obtiene un resumen del estado de los documentos
// GET /users/:userId/documents/summary
func (h *DocumentHandler) GetDocumentSummary(c *gin.Context) {
	clubID := c.GetString("club_id")
	userID := c.Param("userId")

	summary, err := h.eligibilityService.GetDocumentSummary(clubID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener resumen"})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// RegisterDocumentRoutes registra las rutas de documentos
func RegisterDocumentRoutes(router *gin.RouterGroup, handler *DocumentHandler, authMiddleware gin.HandlerFunc) {
	users := router.Group("/users")
	users.Use(authMiddleware)
	{
		// Rutas de documentos
		users.POST("/:userId/documents", handler.UploadDocument)
		users.GET("/:userId/documents", handler.ListDocuments)
		users.GET("/:userId/documents/summary", handler.GetDocumentSummary)
		users.GET("/:userId/documents/:docId", handler.GetDocument)
		users.PUT("/:userId/documents/:docId/validate", handler.ValidateDocument)
		users.DELETE("/:userId/documents/:docId", handler.DeleteDocument)

		// Elegibilidad
		users.GET("/:userId/eligibility", handler.CheckEligibility)
	}
}
