package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	attendanceDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/attendance/domain"
	membershipDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/domain"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
)

// PlayerStatusService maneja la lógica de negocio para obtener el estado unificado de jugadores
type PlayerStatusService struct {
	membershipRepo membershipDomain.MembershipRepository
	documentRepo   userDomain.UserDocumentRepository
	attendanceRepo attendanceDomain.AttendanceRepository
	userRepo       userDomain.UserRepository
}

// NewPlayerStatusService crea una nueva instancia del servicio
func NewPlayerStatusService(
	membershipRepo membershipDomain.MembershipRepository,
	documentRepo userDomain.UserDocumentRepository,
	attendanceRepo attendanceDomain.AttendanceRepository,
	userRepo userDomain.UserRepository,
) *PlayerStatusService {
	return &PlayerStatusService{
		membershipRepo: membershipRepo,
		documentRepo:   documentRepo,
		attendanceRepo: attendanceRepo,
		userRepo:       userRepo,
	}
}

// PlayerStatusFlags contiene los indicadores de estado de un jugador
type PlayerStatusFlags struct {
	FinancialStatus string  `json:"financial_status"` // "ACTIVE" | "DEBTOR"
	MedicalStatus   string  `json:"medical_status"`   // "VALID" | "EXPIRED" | "MISSING"
	AttendanceRate  float64 `json:"attendance_rate"`  // 0.0 - 1.0
	IsInhabilitado  bool    `json:"is_inhabilitado"`  // true si tiene deuda O apto vencido
}

// PlayerWithStatus representa un jugador con su estado unificado
type PlayerWithStatus struct {
	User        *userDomain.User  `json:"user"`
	StatusFlags PlayerStatusFlags `json:"status_flags"`
}

// GetPlayerStatus obtiene el estado unificado de un jugador
func (s *PlayerStatusService) GetPlayerStatus(ctx context.Context, clubID, userID string) (PlayerStatusFlags, error) {
	var flags PlayerStatusFlags

	// 1. Estado Financiero (desde Membership)
	userUUID, err := uuid.Parse(userID)
	if err == nil {
		memberships, err := s.membershipRepo.GetByUserID(ctx, clubID, userUUID)
		if err == nil && len(memberships) > 0 {
			membership := memberships[0]
			if membership.OutstandingBalance.IsPositive() {
				flags.FinancialStatus = "DEBTOR"
			} else {
				flags.FinancialStatus = "ACTIVE"
			}
		} else {
			flags.FinancialStatus = "NO_MEMBERSHIP"
		}
	} else {
		flags.FinancialStatus = "UNKNOWN"
	}

	// 2. Estado Médico (desde UserDocuments)
	docs, err := s.documentRepo.GetByUserID(ctx, clubID, userID)
	if err == nil {
		hasValidEMMAC := false
		for _, doc := range docs {
			if doc.Type == userDomain.DocumentTypeEMMACMedical &&
				doc.Status == userDomain.DocumentStatusValid &&
				!doc.IsExpired() {
				hasValidEMMAC = true
				break
			}
		}
		if hasValidEMMAC {
			flags.MedicalStatus = "VALID"
		} else {
			// Verificar si hay algún EMMAC (aunque esté vencido)
			hasAnyEMMAC := false
			for _, doc := range docs {
				if doc.Type == userDomain.DocumentTypeEMMACMedical {
					hasAnyEMMAC = true
					break
				}
			}
			if hasAnyEMMAC {
				flags.MedicalStatus = "EXPIRED"
			} else {
				flags.MedicalStatus = "MISSING"
			}
		}
	} else {
		flags.MedicalStatus = "MISSING"
	}

	// 3. Tasa de Asistencia (último mes)
	flags.AttendanceRate = s.calculateAttendanceRate(ctx, clubID, userID)

	// 4. Regla de Inhabilitación
	flags.IsInhabilitado = flags.FinancialStatus == "DEBTOR" ||
		flags.MedicalStatus != "VALID"

	return flags, nil
}

// GetTeamPlayersWithStatus obtiene todos los jugadores de un equipo con su estado
func (s *PlayerStatusService) GetTeamPlayersWithStatus(ctx context.Context, clubID string, userIDs []string) ([]PlayerWithStatus, error) {
	players := make([]PlayerWithStatus, 0, len(userIDs))

	for _, userID := range userIDs {
		// Obtener usuario
		user, err := s.userRepo.GetByID(ctx, clubID, userID)
		if err != nil {
			continue // Skip si no se encuentra el usuario
		}

		// Obtener estado
		status, err := s.GetPlayerStatus(ctx, clubID, userID)
		if err != nil {
			// Si hay error, usar valores por defecto
			status = PlayerStatusFlags{
				FinancialStatus: "UNKNOWN",
				MedicalStatus:   "MISSING",
				AttendanceRate:  0.0,
				IsInhabilitado:  true,
			}
		}

		players = append(players, PlayerWithStatus{
			User:        user,
			StatusFlags: status,
		})
	}

	return players, nil
}

// calculateAttendanceRate calcula el porcentaje de asistencia del último mes
func (s *PlayerStatusService) calculateAttendanceRate(ctx context.Context, clubID, userID string) float64 {
	// Calcular rango del último mes
	now := time.Now()
	oneMonthAgo := now.AddDate(0, -1, 0)

	// Consultar estadísticas reales desde el repositorio
	present, total, err := s.attendanceRepo.GetAttendanceStats(ctx, clubID, userID, oneMonthAgo, now)
	if err != nil {
		// En caso de error, retornar 0 (desconocido)
		return 0.0
	}

	if total == 0 {
		// Sin registros de asistencia, no se puede calcular
		return 0.0
	}

	return float64(present) / float64(total)
}

// GetInhabilitadoPlayers obtiene lista de jugadores inhabilitados de un equipo
func (s *PlayerStatusService) GetInhabilitadoPlayers(ctx context.Context, clubID string, userIDs []string) ([]PlayerWithStatus, error) {
	allPlayers, err := s.GetTeamPlayersWithStatus(ctx, clubID, userIDs)
	if err != nil {
		return nil, err
	}

	inhabilitados := make([]PlayerWithStatus, 0)
	for _, player := range allPlayers {
		if player.StatusFlags.IsInhabilitado {
			inhabilitados = append(inhabilitados, player)
		}
	}

	return inhabilitados, nil
}

// GetPlayerIssues retorna una lista de problemas específicos de un jugador
func (s *PlayerStatusService) GetPlayerIssues(ctx context.Context, clubID, userID string) ([]string, error) {
	status, err := s.GetPlayerStatus(ctx, clubID, userID)
	if err != nil {
		return nil, err
	}

	issues := []string{}

	if status.FinancialStatus == "DEBTOR" {
		issues = append(issues, "Tiene deuda pendiente")
	}

	switch status.MedicalStatus {
	case "EXPIRED":
		issues = append(issues, "Apto médico vencido")
	case "MISSING":
		issues = append(issues, "Apto médico faltante")
	}

	if status.AttendanceRate < 0.5 {
		issues = append(issues, "Baja asistencia (menos del 50%)")
	} else if status.AttendanceRate < 0.7 {
		issues = append(issues, "Asistencia regular (menos del 70%)")
	}

	return issues, nil
}
