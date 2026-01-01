package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/attendance/domain"
	"gorm.io/gorm"
)

type PostgresAttendanceRepository struct {
	db *gorm.DB
}

func NewPostgresAttendanceRepository(db *gorm.DB) *PostgresAttendanceRepository {
	return &PostgresAttendanceRepository{db: db}
}

// Table Models
type AttendanceListModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Date      time.Time
	Group     string `gorm:"column:group_name"`
	CoachID   string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	// Relations
	Records []AttendanceRecordModel `gorm:"foreignKey:AttendanceListID"`
}

func (AttendanceListModel) TableName() string {
	return "attendance_lists"
}

type AttendanceRecordModel struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	AttendanceListID uuid.UUID
	UserID           string
	Status           string // PRESENT, ABSENT, LATE
	Notes            string
}

func (AttendanceRecordModel) TableName() string {
	return "attendance_records"
}

func (r *PostgresAttendanceRepository) CreateList(list *domain.AttendanceList) error {
	model := AttendanceListModel{
		ID:        list.ID,
		Date:      list.Date,
		Group:     list.Group,
		CoachID:   list.CoachID,
		CreatedAt: list.CreatedAt,
		UpdatedAt: list.UpdatedAt,
	}
	return r.db.Create(&model).Error
}

func (r *PostgresAttendanceRepository) GetListByID(id uuid.UUID) (*domain.AttendanceList, error) {
	var model AttendanceListModel
	if err := r.db.Preload("Records").First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.mapToDomain(&model), nil
}

func (r *PostgresAttendanceRepository) GetListByGroupAndDate(group string, date time.Time) (*domain.AttendanceList, error) {
	var model AttendanceListModel
	// Assuming date match matches the day.
	// We might need strict equality or range if timestamp includes time.
	// For simplicity, let's assume we store truncated dates or query range.
	// Here I'll verify exact match assuming logic truncates it.
	if err := r.db.Preload("Records").Where("group_name = ? AND date = ?", group, date).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		// Try column 'group' if 'group_name' fails. SQL maps 'Group' struct field to 'group' column usually but 'Group' is keyword.
		// Let's rely on explicit column naming in Schema or GORM default.
		// GORM `Group` -> `group`. `group` is reserved in SQL. Best to map to `group_name` in DB.
		return nil, err
	}
	return r.mapToDomain(&model), nil
}

func (r *PostgresAttendanceRepository) UpsertRecord(record *domain.AttendanceRecord) error {
	model := AttendanceRecordModel{
		ID:               record.ID,
		AttendanceListID: record.AttendanceListID,
		UserID:           record.UserID,
		Status:           string(record.Status),
		Notes:            record.Notes,
	}
	// On Conflict Update
	// Postgres: ON CONFLICT (id) DO UPDATE
	// But we might want unique (list_id, user_id).
	// Let's assume ID is provided or we query first.
	// Ideally Logic handles ID generation.
	return r.db.Save(&model).Error
}

func (r *PostgresAttendanceRepository) UpdateRecord(record *domain.AttendanceRecord) error {
	return r.UpsertRecord(record)
}

func (r *PostgresAttendanceRepository) mapToDomain(model *AttendanceListModel) *domain.AttendanceList {
	records := make([]domain.AttendanceRecord, len(model.Records))
	for i, rec := range model.Records {
		records[i] = domain.AttendanceRecord{
			ID:               rec.ID,
			AttendanceListID: rec.AttendanceListID,
			UserID:           rec.UserID,
			Status:           domain.AttendanceStatus(rec.Status),
			Notes:            rec.Notes,
		}
	}
	return &domain.AttendanceList{
		ID:        model.ID,
		Date:      model.Date,
		Group:     model.Group,
		CoachID:   model.CoachID,
		Records:   records,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
