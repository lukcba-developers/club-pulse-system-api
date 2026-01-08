package application

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
)

// EligibilityService maneja la lógica de negocio para verificar la elegibilidad de un usuario
type EligibilityService struct {
	docRepo domain.UserDocumentRepository
}

// NewEligibilityService crea una nueva instancia del servicio
func NewEligibilityService(docRepo domain.UserDocumentRepository) *EligibilityService {
	return &EligibilityService{
		docRepo: docRepo,
	}
}

// EligibilityResult contiene el resultado de la verificación de elegibilidad
type EligibilityResult struct {
	IsEligible bool     `json:"is_eligible"`
	Issues     []string `json:"issues,omitempty"`
	HasDNI     bool     `json:"has_dni"`
	HasEMMAC   bool     `json:"has_emmac"`
}

// CheckEligibility verifica si un usuario es elegible para participar
// Un usuario es elegible si tiene:
// 1. DNI (Frente y/o Dorso) válido
// 2. Apto Médico (EMMAC) válido y no vencido
func (s *EligibilityService) CheckEligibility(clubID, userID string) (*EligibilityResult, error) {
	docs, err := s.docRepo.GetByUserID(clubID, userID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo documentos: %w", err)
	}

	result := &EligibilityResult{
		IsEligible: true,
		Issues:     []string{},
		HasDNI:     false,
		HasEMMAC:   false,
	}

	// Verificar DNI
	for _, doc := range docs {
		if (doc.Type == domain.DocumentTypeDNIFront || doc.Type == domain.DocumentTypeDNIBack) &&
			doc.Status == domain.DocumentStatusValid {
			result.HasDNI = true
			break
		}
	}

	// Verificar EMMAC (Apto Médico)
	for _, doc := range docs {
		if doc.Type == domain.DocumentTypeEMMACMedical {
			if doc.Status == domain.DocumentStatusExpired || doc.IsExpired() {
				result.Issues = append(result.Issues, "Apto físico vencido")
			} else if doc.Status == domain.DocumentStatusValid {
				result.HasEMMAC = true
			} else if doc.Status == domain.DocumentStatusPending {
				result.Issues = append(result.Issues, "Apto físico pendiente de validación")
			} else if doc.Status == domain.DocumentStatusRejected {
				result.Issues = append(result.Issues, "Apto físico rechazado")
			}
			break
		}
	}

	// Agregar issues si faltan documentos
	if !result.HasDNI {
		result.Issues = append(result.Issues, "DNI faltante o no validado")
	}
	if !result.HasEMMAC {
		// Solo agregar si no hay otro issue relacionado con EMMAC
		hasEMMACIssue := false
		for _, issue := range result.Issues {
			if issue == "Apto físico vencido" || issue == "Apto físico pendiente de validación" || issue == "Apto físico rechazado" {
				hasEMMACIssue = true
				break
			}
		}
		if !hasEMMACIssue {
			result.Issues = append(result.Issues, "Apto físico faltante")
		}
	}

	// El usuario es elegible solo si tiene ambos documentos válidos
	result.IsEligible = result.HasDNI && result.HasEMMAC && len(result.Issues) == 0

	return result, nil
}

// GetDocumentSummary obtiene un resumen del estado de los documentos de un usuario
func (s *EligibilityService) GetDocumentSummary(clubID, userID string) (map[domain.DocumentType]domain.DocumentStatus, error) {
	docs, err := s.docRepo.GetByUserID(clubID, userID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo documentos: %w", err)
	}

	summary := make(map[domain.DocumentType]domain.DocumentStatus)

	// Obtener el documento más reciente de cada tipo
	for _, doc := range docs {
		// Si el documento está vencido, actualizar su estado
		if doc.IsExpired() && doc.Status == domain.DocumentStatusValid {
			summary[doc.Type] = domain.DocumentStatusExpired
		} else {
			// Solo agregar si no existe o si es más reciente
			if _, exists := summary[doc.Type]; !exists {
				summary[doc.Type] = doc.Status
			}
		}
	}

	return summary, nil
}

// ValidateDocument valida o rechaza un documento
func (s *EligibilityService) ValidateDocument(clubID string, docID string, validatorID string, approve bool, notes string) error {
	docUUID, err := parseUUID(docID)
	if err != nil {
		return fmt.Errorf("ID de documento inválido: %w", err)
	}

	doc, err := s.docRepo.GetByID(clubID, docUUID)
	if err != nil {
		return fmt.Errorf("documento no encontrado: %w", err)
	}

	// Verificar que el documento puede ser validado
	if !doc.CanBeValidated() {
		return fmt.Errorf("el documento no puede ser validado (estado: %s, vencido: %v)", doc.Status, doc.IsExpired())
	}

	// Actualizar el documento
	now := time.Now()
	doc.ValidatedAt = &now
	doc.ValidatedBy = &validatorID

	if approve {
		doc.Status = domain.DocumentStatusValid
		doc.RejectionNotes = ""
	} else {
		doc.Status = domain.DocumentStatusRejected
		doc.RejectionNotes = notes
	}

	return s.docRepo.Update(doc)
}

// Helper function para parsear UUID
func parseUUID(id string) (uuid.UUID, error) {
	return uuid.Parse(id)
}
