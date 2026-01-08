package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/application"
)

// LeagueExportHandler maneja las peticiones de exportación para la Liga
type LeagueExportHandler struct {
	exportService *application.LeagueExportService
	userRepo      interface {
		GetByID(clubID, id string) (interface{}, error)
		List(clubID string, limit, offset int, filters map[string]interface{}) ([]interface{}, error)
	}
}

// NewLeagueExportHandler crea una nueva instancia del handler
func NewLeagueExportHandler(exportService *application.LeagueExportService) *LeagueExportHandler {
	return &LeagueExportHandler{
		exportService: exportService,
	}
}

// ExportLeagueFolder genera y descarga el PDF de la Carpeta de Liga
// GET /teams/:teamId/league-export
func (h *LeagueExportHandler) ExportLeagueFolder(c *gin.Context) {
	clubID := c.GetString("club_id")
	teamID := c.Param("teamId")

	// Verificar permisos (solo admins y coaches)
	role := c.GetString("role")
	if role != "ADMIN" && role != "SUPER_ADMIN" && role != "COACH" {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permiso para exportar carpetas de liga"})
		return
	}

	// TODO: Obtener miembros del equipo desde el repositorio de Team
	// Por ahora, usamos datos de ejemplo
	members := []application.TeamMember{
		// Esto debería venir de la base de datos
	}

	// Obtener documentos de cada miembro
	for i := range members {
		docs, err := h.exportService.GetMemberDocuments(clubID, members[i].ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo documentos"})
			return
		}
		members[i].Documents = docs
	}

	// Filtrar solo miembros elegibles (opcional, basado en query param)
	onlyEligible := c.Query("only_eligible") == "true"
	if onlyEligible {
		members = h.exportService.FilterEligibleMembers(members)
	}

	// Generar PDF
	teamName := fmt.Sprintf("Equipo %s", teamID) // TODO: Obtener nombre real del equipo

	// Configurar headers para descarga
	filename := fmt.Sprintf("carpeta_liga_%s_%s.pdf", teamID, uuid.New().String()[:8])
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	// Generar y escribir PDF directamente a la respuesta
	if err := h.exportService.GenerateLeagueFolder(teamName, members, c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generando PDF"})
		return
	}
}

// RegisterLeagueExportRoutes registra las rutas de exportación
func RegisterLeagueExportRoutes(router *gin.RouterGroup, handler *LeagueExportHandler, authMiddleware gin.HandlerFunc) {
	teams := router.Group("/teams")
	teams.Use(authMiddleware)
	{
		teams.GET("/:teamId/league-export", handler.ExportLeagueFolder)
	}
}
