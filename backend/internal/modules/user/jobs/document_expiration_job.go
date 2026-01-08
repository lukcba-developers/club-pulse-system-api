package jobs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/notification/service"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
)

// DocumentExpirationJob maneja la verificación y notificación de documentos vencidos
type DocumentExpirationJob struct {
	docRepo      domain.UserDocumentRepository
	notifService *service.NotificationService
}

// NewDocumentExpirationJob crea una nueva instancia del job
func NewDocumentExpirationJob(docRepo domain.UserDocumentRepository, notifService *service.NotificationService) *DocumentExpirationJob {
	return &DocumentExpirationJob{
		docRepo:      docRepo,
		notifService: notifService,
	}
}

// Run ejecuta el job de verificación de vencimientos
func (j *DocumentExpirationJob) Run() {
	log.Println("🔍 Ejecutando job de vencimiento de documentos...")

	// 1. Notificar documentos que vencen en 30 días
	j.notifyExpiringDocuments(30, "⚠️ Tu documento vence en 30 días")

	// 2. Notificar documentos que vencen en 7 días
	j.notifyExpiringDocuments(7, "🚨 Tu documento vence en 7 días")

	// 3. Marcar documentos vencidos
	j.markExpiredDocuments()

	log.Println("✅ Job de vencimiento de documentos completado")
}

// notifyExpiringDocuments notifica a usuarios sobre documentos próximos a vencer
func (j *DocumentExpirationJob) notifyExpiringDocuments(days int, title string) {
	docs, err := j.docRepo.GetExpiringDocuments("", days)
	if err != nil {
		log.Printf("❌ Error obteniendo documentos que vencen en %d días: %v\n", days, err)
		return
	}

	log.Printf("📋 Encontrados %d documentos que vencen en %d días\n", len(docs), days)

	for _, doc := range docs {
		message := fmt.Sprintf("Tu %s vence el %s. Por favor, renuévalo pronto.",
			getDocumentTypeName(doc.Type),
			doc.ExpirationDate.Format("02/01/2006"))

		notification := service.Notification{
			RecipientID: doc.UserID,
			Type:        service.NotificationTypeEmail,
			Subject:     title,
			Message:     message,
		}

		if err := j.notifService.Send(context.Background(), notification); err != nil {
			log.Printf("❌ Error enviando notificación a usuario %s: %v\n", doc.UserID, err)
		} else {
			log.Printf("✉️ Notificación enviada a usuario %s sobre %s\n", doc.UserID, doc.Type)
		}
	}
}

// markExpiredDocuments marca documentos vencidos como EXPIRED
func (j *DocumentExpirationJob) markExpiredDocuments() {
	docs, err := j.docRepo.GetExpiredDocuments("")
	if err != nil {
		log.Printf("❌ Error obteniendo documentos vencidos: %v\n", err)
		return
	}

	log.Printf("📋 Encontrados %d documentos vencidos\n", len(docs))

	for _, doc := range docs {
		// Actualizar estado a EXPIRED
		doc.Status = domain.DocumentStatusExpired
		if err := j.docRepo.Update(&doc); err != nil {
			log.Printf("❌ Error actualizando documento %s: %v\n", doc.ID, err)
			continue
		}

		// Notificar al usuario
		message := fmt.Sprintf("Tu %s ha vencido. Por favor, sube un nuevo documento.",
			getDocumentTypeName(doc.Type))

		notification := service.Notification{
			RecipientID: doc.UserID,
			Type:        service.NotificationTypeEmail,
			Subject:     "❌ Documento vencido",
			Message:     message,
		}

		if err := j.notifService.Send(context.Background(), notification); err != nil {
			log.Printf("❌ Error enviando notificación a usuario %s: %v\n", doc.UserID, err)
		} else {
			log.Printf("✉️ Notificación de vencimiento enviada a usuario %s sobre %s\n", doc.UserID, doc.Type)
		}
	}
}

// getDocumentTypeName retorna el nombre legible del tipo de documento
func getDocumentTypeName(docType domain.DocumentType) string {
	switch docType {
	case domain.DocumentTypeDNIFront:
		return "DNI (Frente)"
	case domain.DocumentTypeDNIBack:
		return "DNI (Dorso)"
	case domain.DocumentTypeEMMACMedical:
		return "Apto Médico (EMMAC)"
	case domain.DocumentTypeLeagueForm:
		return "Formulario de Liga"
	case domain.DocumentTypeInsurance:
		return "Seguro"
	default:
		return string(docType)
	}
}

// RunPeriodically ejecuta el job periódicamente
func (j *DocumentExpirationJob) RunPeriodically(interval time.Duration, stop <-chan bool) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Printf("⏰ Job de vencimiento de documentos iniciado (intervalo: %v)\n", interval)

	// Ejecutar inmediatamente al inicio
	j.Run()

	for {
		select {
		case <-ticker.C:
			j.Run()
		case <-stop:
			log.Println("🛑 Job de vencimiento de documentos detenido")
			return
		}
	}
}
