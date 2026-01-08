package application

import (
	"fmt"
	"io"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
)

// TeamMember representa un miembro del equipo para la exportación
type TeamMember struct {
	ID        string
	Name      string
	DNI       string
	BirthDate *time.Time
	Documents []domain.UserDocument
}

// LeagueExportService maneja la generación de PDFs para la Liga
type LeagueExportService struct {
	docRepo domain.UserDocumentRepository
}

// NewLeagueExportService crea una nueva instancia del servicio
func NewLeagueExportService(docRepo domain.UserDocumentRepository) *LeagueExportService {
	return &LeagueExportService{
		docRepo: docRepo,
	}
}

// GenerateLeagueFolder genera un PDF con la "Carpeta de Liga"
// Incluye:
// - Página 1: Lista de Buena Fe (tabla resumen)
// - Páginas siguientes: DNI + Apto Físico de cada jugador
func (s *LeagueExportService) GenerateLeagueFolder(teamName string, members []TeamMember, writer io.Writer) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)

	// Página 1: Lista de Buena Fe
	s.addSummaryPage(pdf, teamName, members)

	// Páginas siguientes: Documentos de cada jugador
	for _, member := range members {
		s.addPlayerDocumentsPage(pdf, member)
	}

	// Escribir PDF al writer
	return pdf.Output(writer)
}

// addSummaryPage agrega la página de resumen (Lista de Buena Fe)
func (s *LeagueExportService) addSummaryPage(pdf *gofpdf.Fpdf, teamName string, members []TeamMember) {
	pdf.AddPage()

	// Título
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "LISTA DE BUENA FE", "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Nombre del equipo
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 8, teamName, "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Fecha de generación
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 6, fmt.Sprintf("Fecha: %s", time.Now().Format("02/01/2006")), "", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Tabla de jugadores
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(200, 200, 200)

	// Headers
	pdf.CellFormat(10, 8, "N°", "1", 0, "C", true, 0, "")
	pdf.CellFormat(60, 8, "Apellido y Nombre", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "DNI", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "Fecha Nac.", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 8, "Estado Doc.", "1", 1, "C", true, 0, "")

	// Filas
	pdf.SetFont("Arial", "", 9)
	for i, member := range members {
		// Determinar estado documental
		status := s.getDocumentStatus(member)
		statusColor := s.getStatusColor(status)

		pdf.CellFormat(10, 7, fmt.Sprintf("%d", i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(60, 7, member.Name, "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 7, member.DNI, "1", 0, "C", false, 0, "")

		birthDate := "-"
		if member.BirthDate != nil {
			birthDate = member.BirthDate.Format("02/01/2006")
		}
		pdf.CellFormat(30, 7, birthDate, "1", 0, "C", false, 0, "")

		// Estado con color
		pdf.SetFillColor(statusColor[0], statusColor[1], statusColor[2])
		pdf.CellFormat(40, 7, status, "1", 1, "C", true, 0, "")
	}

	// Pie de página
	pdf.Ln(10)
	pdf.SetFont("Arial", "I", 8)
	pdf.CellFormat(0, 5, fmt.Sprintf("Total de jugadores: %d", len(members)), "", 1, "L", false, 0, "")
}

// addPlayerDocumentsPage agrega una página con los documentos de un jugador
func (s *LeagueExportService) addPlayerDocumentsPage(pdf *gofpdf.Fpdf, member TeamMember) {
	pdf.AddPage()

	// Título con nombre del jugador
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, member.Name, "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// DNI
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 6, fmt.Sprintf("DNI: %s", member.DNI), "", 1, "L", false, 0, "")
	pdf.Ln(5)

	// Buscar documentos
	hasDNI := false
	hasEMMAC := false

	for _, doc := range member.Documents {
		if doc.Status != domain.DocumentStatusValid {
			continue
		}

		switch doc.Type {
		case domain.DocumentTypeDNIFront, domain.DocumentTypeDNIBack:
			if !hasDNI {
				pdf.SetFont("Arial", "B", 11)
				pdf.CellFormat(0, 7, "DNI:", "", 1, "L", false, 0, "")
				pdf.SetFont("Arial", "", 9)
				pdf.CellFormat(0, 6, fmt.Sprintf("Archivo: %s", doc.FileURL), "", 1, "L", false, 0, "")
				pdf.Ln(3)
				hasDNI = true

				// TODO: Insertar imagen del DNI
				// pdf.Image(doc.FileURL, 15, pdf.GetY(), 90, 0, false, "", 0, "")
			}

		case domain.DocumentTypeEMMACMedical:
			if !hasEMMAC {
				pdf.SetFont("Arial", "B", 11)
				pdf.CellFormat(0, 7, "Apto Médico (EMMAC):", "", 1, "L", false, 0, "")
				pdf.SetFont("Arial", "", 9)
				pdf.CellFormat(0, 6, fmt.Sprintf("Archivo: %s", doc.FileURL), "", 1, "L", false, 0, "")

				if doc.ExpirationDate != nil {
					pdf.CellFormat(0, 6, fmt.Sprintf("Vencimiento: %s", doc.ExpirationDate.Format("02/01/2006")), "", 1, "L", false, 0, "")
				}
				pdf.Ln(3)
				hasEMMAC = true

				// TODO: Insertar imagen del apto médico
				// pdf.Image(doc.FileURL, 15, pdf.GetY(), 90, 0, false, "", 0, "")
			}
		}
	}

	// Advertencias si faltan documentos
	if !hasDNI || !hasEMMAC {
		pdf.Ln(5)
		pdf.SetFont("Arial", "B", 10)
		pdf.SetTextColor(255, 0, 0)
		pdf.CellFormat(0, 7, "⚠ DOCUMENTACIÓN INCOMPLETA", "", 1, "L", false, 0, "")
		pdf.SetTextColor(0, 0, 0)
		pdf.SetFont("Arial", "", 9)

		if !hasDNI {
			pdf.CellFormat(0, 6, "- Falta DNI", "", 1, "L", false, 0, "")
		}
		if !hasEMMAC {
			pdf.CellFormat(0, 6, "- Falta Apto Médico", "", 1, "L", false, 0, "")
		}
	}
}

// getDocumentStatus determina el estado documental de un jugador
func (s *LeagueExportService) getDocumentStatus(member TeamMember) string {
	hasDNI := false
	hasEMMAC := false

	for _, doc := range member.Documents {
		if doc.Status != domain.DocumentStatusValid {
			continue
		}

		if doc.Type == domain.DocumentTypeDNIFront || doc.Type == domain.DocumentTypeDNIBack {
			hasDNI = true
		}

		if doc.Type == domain.DocumentTypeEMMACMedical && !doc.IsExpired() {
			hasEMMAC = true
		}
	}

	if hasDNI && hasEMMAC {
		return "✓ Completo"
	} else if hasDNI || hasEMMAC {
		return "⚠ Incompleto"
	}
	return "✗ Sin Docs"
}

// getStatusColor retorna el color RGB según el estado
func (s *LeagueExportService) getStatusColor(status string) [3]int {
	switch status {
	case "✓ Completo":
		return [3]int{144, 238, 144} // Verde claro
	case "⚠ Incompleto":
		return [3]int{255, 255, 153} // Amarillo claro
	default:
		return [3]int{255, 182, 193} // Rojo claro
	}
}

// FilterEligibleMembers filtra solo los miembros con documentación válida
func (s *LeagueExportService) FilterEligibleMembers(members []TeamMember) []TeamMember {
	eligible := []TeamMember{}

	for _, member := range members {
		status := s.getDocumentStatus(member)
		if status == "✓ Completo" {
			eligible = append(eligible, member)
		}
	}

	return eligible
}

// GetMemberDocuments obtiene los documentos de un miembro
func (s *LeagueExportService) GetMemberDocuments(clubID, userID string) ([]domain.UserDocument, error) {
	return s.docRepo.GetByUserID(clubID, userID)
}
