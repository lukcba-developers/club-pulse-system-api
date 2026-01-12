package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/domain"
	"gorm.io/gorm"
)

type PostgresFacilityRepository struct {
	db *gorm.DB
}

func NewPostgresFacilityRepository(db *gorm.DB) *PostgresFacilityRepository {
	log.Println("DEBUG: Running AutoMigrate for Facilities...")
	err := db.AutoMigrate(&FacilityModel{}, &MaintenanceTaskModel{}, &EquipmentModel{})
	if err != nil {
		log.Printf("DEBUG: AutoMigrate Failed: %v", err)
	} else {
		log.Println("DEBUG: AutoMigrate Success")
	}
	return &PostgresFacilityRepository{db: db}
}

// FacilityModel mirrors domain.Facility but with GORM tags
type FacilityModel struct {
	ID             string                `gorm:"primaryKey"`
	Name           string                `gorm:"not null"`
	Type           string                `gorm:"not null"`
	Status         string                `gorm:"default:'active'"`
	Capacity       int                   `gorm:"not null"`
	HourlyRate     float64               `gorm:"not null"`
	OpeningHour    int                   `gorm:"default:8"`
	ClosingHour    int                   `gorm:"default:23"`
	GuestFee       float64               `gorm:"default:0"`
	Specifications domain.Specifications `gorm:"type:jsonb;serializer:json"` // Postgres JSONB
	Location       domain.Location       `gorm:"type:jsonb;serializer:json"`
	ClubID         string                `gorm:"index;not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (FacilityModel) TableName() string {
	return "facilities"
}

func (r *PostgresFacilityRepository) Create(ctx context.Context, facility *domain.Facility) error {
	model := FacilityModel{
		ID:             facility.ID,
		Name:           facility.Name,
		Type:           string(facility.Type),
		Status:         string(facility.Status),
		Capacity:       facility.Capacity,
		HourlyRate:     facility.HourlyRate,
		OpeningHour:    facility.OpeningHour,
		ClosingHour:    facility.ClosingHour,
		GuestFee:       facility.GuestFee,
		Specifications: facility.Specifications,
		Location:       facility.Location,
		ClubID:         facility.ClubID,
		CreatedAt:      facility.CreatedAt,
		UpdatedAt:      facility.UpdatedAt,
	}
	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *PostgresFacilityRepository) GetByID(ctx context.Context, clubID, id string) (*domain.Facility, error) {
	var model FacilityModel
	result := r.db.WithContext(ctx).Where("id = ? AND club_id = ?", id, clubID).First(&model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return r.toDomain(model), nil
}

func (r *PostgresFacilityRepository) List(ctx context.Context, clubID string, limit, offset int) ([]*domain.Facility, error) {
	var models []FacilityModel
	result := r.db.WithContext(ctx).Where("club_id = ?", clubID).Limit(limit).Offset(offset).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	facilities := make([]*domain.Facility, len(models))
	for i, m := range models {
		facilities[i] = r.toDomain(m)
	}
	return facilities, nil
}

func (r *PostgresFacilityRepository) Update(ctx context.Context, facility *domain.Facility) error {
	model := FacilityModel{
		ID:             facility.ID,
		Name:           facility.Name,
		Type:           string(facility.Type),
		Status:         string(facility.Status),
		Capacity:       facility.Capacity,
		HourlyRate:     facility.HourlyRate,
		OpeningHour:    facility.OpeningHour,
		ClosingHour:    facility.ClosingHour,
		GuestFee:       facility.GuestFee,
		Specifications: facility.Specifications,
		Location:       facility.Location,
		ClubID:         facility.ClubID,
		CreatedAt:      facility.CreatedAt,
		UpdatedAt:      time.Now(), // Update timestamp
	}
	// Save updates all fields (including zero values) which is what we want for struct replacement
	return r.db.WithContext(ctx).Save(&model).Error
}

func (r *PostgresFacilityRepository) toDomain(m FacilityModel) *domain.Facility {
	return &domain.Facility{
		ID:             m.ID,
		Name:           m.Name,
		Type:           domain.FacilityType(m.Type),
		Status:         domain.FacilityStatus(m.Status),
		Capacity:       m.Capacity,
		HourlyRate:     m.HourlyRate,
		OpeningHour:    m.OpeningHour,
		ClosingHour:    m.ClosingHour,
		GuestFee:       m.GuestFee,
		Specifications: m.Specifications,
		Location:       m.Location,
		ClubID:         m.ClubID,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

// Maintenance Logic

type MaintenanceTaskModel struct {
	ID          string    `gorm:"primaryKey"`
	FacilityID  string    `gorm:"not null;index"`
	EquipmentID *string   `gorm:"index"`
	Title       string    `gorm:"not null"`
	Description string    `gorm:"type:text"`
	Status      string    `gorm:"not null"`
	Type        string    `gorm:"not null"`
	StartTime   time.Time `gorm:"not null;index"`
	EndTime     time.Time `gorm:"not null;index"`
	CreatedBy   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (MaintenanceTaskModel) TableName() string {
	return "maintenance_tasks"
}

type EquipmentModel struct {
	ID           string `gorm:"primaryKey"`
	FacilityID   string `gorm:"not null;index"`
	Name         string `gorm:"not null"`
	Type         string `gorm:"not null"`
	Condition    string `gorm:"not null"`
	Status       string `gorm:"not null"`
	PurchaseDate *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (EquipmentModel) TableName() string {
	return "equipment"
}

// Implement MaintenanceRepository

func (r *PostgresFacilityRepository) AutoMigrateMaintenance() error {
	return r.db.AutoMigrate(&MaintenanceTaskModel{}, &EquipmentModel{})
}

func (r *PostgresFacilityRepository) CreateMaintenance(task *domain.MaintenanceTask) error {
	model := MaintenanceTaskModel{
		ID:          task.ID,
		FacilityID:  task.FacilityID,
		EquipmentID: task.EquipmentID,
		Title:       task.Title,
		Description: task.Description,
		Status:      string(task.Status),
		Type:        string(task.Type),
		StartTime:   task.StartTime,
		EndTime:     task.EndTime,
		CreatedBy:   task.CreatedBy,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
	return r.db.Create(&model).Error
}

func (r *PostgresFacilityRepository) GetMaintenanceByID(id string) (*domain.MaintenanceTask, error) {
	var model MaintenanceTaskModel
	if err := r.db.First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &domain.MaintenanceTask{
		ID:          model.ID,
		FacilityID:  model.FacilityID,
		EquipmentID: model.EquipmentID,
		Title:       model.Title,
		Description: model.Description,
		Status:      domain.MaintenanceStatus(model.Status),
		Type:        domain.MaintenanceType(model.Type),
		StartTime:   model.StartTime,
		EndTime:     model.EndTime,
		CreatedBy:   model.CreatedBy,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}, nil
}

func (r *PostgresFacilityRepository) ListMaintenanceByFacility(ctx context.Context, facilityID string) ([]*domain.MaintenanceTask, error) {
	var models []MaintenanceTaskModel
	if err := r.db.WithContext(ctx).Where("facility_id = ?", facilityID).Find(&models).Error; err != nil {
		return nil, err
	}
	tasks := make([]*domain.MaintenanceTask, len(models))
	for i, m := range models {
		tasks[i] = &domain.MaintenanceTask{
			ID:          m.ID,
			FacilityID:  m.FacilityID,
			EquipmentID: m.EquipmentID,
			Title:       m.Title,
			Description: m.Description,
			Status:      domain.MaintenanceStatus(m.Status),
			Type:        domain.MaintenanceType(m.Type),
			StartTime:   m.StartTime,
			EndTime:     m.EndTime,
			CreatedBy:   m.CreatedBy,
			CreatedAt:   m.CreatedAt,
			UpdatedAt:   m.UpdatedAt,
		}
	}
	return tasks, nil
}

func (r *PostgresFacilityRepository) HasConflict(ctx context.Context, clubID, facilityID string, startTime, endTime time.Time) (bool, error) {
	var count int64
	// Check for any maintenance task that overlaps and is active.
	// Tasks are linked to facility. We trust facilityID matches clubID via previous lookups or join if strictly necessary.
	// But maintenance tasks don't have ClubID on them explicitly, they rely on FacilityID.
	// So just filtering by FacilityID is technically enough if we trust the facilityID belongs to the club.
	// However, for strictness, we could join or check facility. But keeping it simple as Facility ownership is verified by ID.
	err := r.db.WithContext(ctx).Model(&MaintenanceTaskModel{}).
		Where("facility_id = ?", facilityID).
		Where("status IN ?", []string{string(domain.MaintenanceStatusScheduled), string(domain.MaintenanceStatusInProgress)}).
		Where("start_time < ? AND end_time > ?", endTime, startTime).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Implement EquipmentRepository

func (r *PostgresFacilityRepository) CreateEquipment(ctx context.Context, equipment *domain.Equipment) error {
	model := EquipmentModel{
		ID:           equipment.ID,
		FacilityID:   equipment.FacilityID,
		Name:         equipment.Name,
		Type:         equipment.Type,
		Condition:    string(equipment.Condition),
		Status:       equipment.Status,
		PurchaseDate: equipment.PurchaseDate,
		CreatedAt:    equipment.CreatedAt,
		UpdatedAt:    equipment.UpdatedAt,
	}
	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *PostgresFacilityRepository) GetEquipmentByID(ctx context.Context, id string) (*domain.Equipment, error) {
	var model EquipmentModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomainEquipment(model), nil
}

func (r *PostgresFacilityRepository) ListEquipmentByFacility(ctx context.Context, facilityID string) ([]*domain.Equipment, error) {
	var models []EquipmentModel
	if err := r.db.WithContext(ctx).Where("facility_id = ?", facilityID).Find(&models).Error; err != nil {
		return nil, err
	}
	equipments := make([]*domain.Equipment, len(models))
	for i, m := range models {
		equipments[i] = r.toDomainEquipment(m)
	}
	return equipments, nil
}

func (r *PostgresFacilityRepository) UpdateEquipment(ctx context.Context, equipment *domain.Equipment) error {
	model := EquipmentModel{
		ID:           equipment.ID,
		FacilityID:   equipment.FacilityID,
		Name:         equipment.Name,
		Type:         equipment.Type,
		Condition:    string(equipment.Condition),
		Status:       equipment.Status,
		PurchaseDate: equipment.PurchaseDate,
		CreatedAt:    equipment.CreatedAt,
		UpdatedAt:    time.Now(),
	}
	return r.db.WithContext(ctx).Save(&model).Error
}

func (r *PostgresFacilityRepository) LoanEquipmentAtomic(ctx context.Context, loan *domain.EquipmentLoan, equipmentID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Check and Update Equipment Status
		result := tx.Model(&EquipmentModel{}).
			Where("id = ? AND status = ?", equipmentID, "available").
			Update("status", "loaned")

		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("equipment is not available for loan")
		}

		// 2. Create Loan Record
		loanMap := map[string]interface{}{
			"id":                 loan.ID,
			"equipment_id":       loan.EquipmentID,
			"user_id":            loan.UserID,
			"loaned_at":          loan.LoanedAt,
			"expected_return_at": loan.ExpectedReturnAt,
			"status":             string(loan.Status),
			"created_at":         loan.CreatedAt,
			"updated_at":         loan.UpdatedAt,
		}

		if err := tx.Table("equipment_loans").Create(&loanMap).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *PostgresFacilityRepository) toDomainEquipment(m EquipmentModel) *domain.Equipment {
	return &domain.Equipment{
		ID:           m.ID,
		FacilityID:   m.FacilityID,
		Name:         m.Name,
		Type:         m.Type,
		Condition:    domain.EquipmentCondition(m.Condition),
		Status:       m.Status,
		PurchaseDate: m.PurchaseDate,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

// SemanticSearch performs vector similarity search using pgvector
func (r *PostgresFacilityRepository) SemanticSearch(ctx context.Context, clubID string, embedding []float32, limit int) ([]*domain.FacilityWithSimilarity, error) {
	// Convert embedding to PostgreSQL vector format
	vectorStr := float32SliceToVectorString(embedding)

	// Query using pgvector cosine distance operator <=>
	query := `
		SELECT 
			id, name, type, status, capacity, hourly_rate, 
			specifications, location, created_at, updated_at, club_id,
			1 - (embedding <=> $1::vector) as similarity
		FROM facilities 
		WHERE embedding IS NOT NULL AND status = 'active' AND club_id = $3
		ORDER BY embedding <=> $1::vector
		LIMIT $2
	`

	rows, err := r.db.WithContext(ctx).Raw(query, vectorStr, limit, clubID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*domain.FacilityWithSimilarity
	for rows.Next() {
		var m FacilityModel
		var similarity float32

		if err := rows.Scan(
			&m.ID, &m.Name, &m.Type, &m.Status, &m.Capacity, &m.HourlyRate,
			&m.Specifications, &m.Location, &m.CreatedAt, &m.UpdatedAt, &m.ClubID,
			&similarity,
		); err != nil {
			return nil, err
		}

		results = append(results, &domain.FacilityWithSimilarity{
			Facility:   r.toDomain(m),
			Similarity: similarity,
		})
	}

	return results, nil
}

// UpdateEmbedding stores the embedding vector for a facility
func (r *PostgresFacilityRepository) UpdateEmbedding(ctx context.Context, facilityID string, embedding []float32) error {
	vectorStr := float32SliceToVectorString(embedding)

	return r.db.WithContext(ctx).Exec(
		"UPDATE facilities SET embedding = $1::vector WHERE id = $2",
		vectorStr, facilityID,
	).Error
}

// float32SliceToVectorString converts a []float32 to PostgreSQL vector string format [1.2,3.4,5.6]
func float32SliceToVectorString(v []float32) string {
	if len(v) == 0 {
		return "[]"
	}

	result := "["
	for i, val := range v {
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf("%f", val)
	}
	result += "]"
	return result
}

// GetImpactedUsers returns a list of user IDs that have bookings during the maintenance window
func (r *PostgresFacilityRepository) GetImpactedUsers(facilityID string, start, end time.Time) ([]string, error) {
	var userIDs []string

	// We use standard SQL check for overlap: (StartA < EndB) and (EndA > StartB)
	// Booking Start < Task End AND Booking End > Task Start
	query := `
		SELECT DISTINCT user_id 
		FROM bookings 
		WHERE facility_id = ? 
		AND status IN ('CONFIRMED', 'PENDING')
		AND start_time < ? 
		AND end_time > ?
	`

	if err := r.db.Raw(query, facilityID, end, start).Scan(&userIDs).Error; err != nil {
		return nil, err
	}

	return userIDs, nil
}
