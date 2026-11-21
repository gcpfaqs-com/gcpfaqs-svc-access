package postgres

import (
	"errors"
	"gcpaccesss-svc-accesss/internal/domain/models"
	"gcpaccesss-svc-accesss/internal/domain/persistence"
	"gcpaccesss-svc-accesss/internal/domain/types"
	"time"

	"gorm.io/gorm"
)

type PostgresAccessRepository struct {
	db *gorm.DB
}

type accessRecord struct {
	ID        string    `gorm:"primaryKey;type:uuid;not null"`
	Username  string    `gorm:"not null"`
	Password  string    `gorm:"not null"`
	Name      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (accessRecord) TableName() string {
	return "accesss"
}

func NewPostgresAccessRepository(db *gorm.DB) persistence.AccessRepository {
	if db == nil {
		panic("postgres repository requires a valid gorm db instance")
	}

	if err := db.AutoMigrate(&accessRecord{}); err != nil {
		panic(err)
	}

	return &PostgresAccessRepository{db: db}
}

func (r *PostgresAccessRepository) Save(access *models.Access) error {
	if access == nil {
		return errors.New("invalid access: access is nil")
	}
	if r.db == nil {
		return errors.New("database connection is nil")
	}
	if access.ID == nil ||
		access.Username == nil ||
		access.Password == nil ||
		access.Name == nil {
		return errors.New("invalid access: access is incomplete")
	}

	payload := accessRecord{
		ID:        access.ID.Value(),
		Username:  access.Username.Value(),
		Password:  access.Password.Value(),
		Name:      access.Name.Value(),
		CreatedAt: access.CreatedAt,
		UpdatedAt: access.UpdatedAt,
	}

	return r.db.Save(&payload).Error
}

func (r *PostgresAccessRepository) Get(id *types.ID) (*models.Access, error) {
	if id == nil {
		return nil, errors.New("invalid access: access id is nil")
	}
	if r.db == nil {
		return nil, errors.New("database connection is nil")
	}

	var payload accessRecord
	result := r.db.Where("id = ?", id.Value()).First(&payload)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}

	accessID, err := types.NewIDFromString(payload.ID)
	if err != nil {
		return nil, err
	}

	accessUsername, err := types.NewUsername(payload.Username)
	if err != nil {
		return nil, err
	}

	accessPassword, err := types.NewPassword(payload.Password)
	if err != nil {
		return nil, err
	}

	accessName, err := types.NewName(payload.Name)
	if err != nil {
		return nil, err
	}

	return &models.Access{
		ID:        accessID,
		Username:  accessUsername,
		Password:  accessPassword,
		Name:      accessName,
		CreatedAt: payload.CreatedAt,
		UpdatedAt: payload.UpdatedAt,
	}, nil
}

func (r *PostgresAccessRepository) GetFromUsername(username string) (*models.Access, error) {
	if username == "" {
		return nil, errors.New("username is empty")
	}
	if r.db == nil {
		return nil, errors.New("database connection is nil")
	}

	var payload accessRecord
	result := r.db.Where("username = ?", username).First(&payload)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}

	accessID, err := types.NewIDFromString(payload.ID)
	if err != nil {
		return nil, err
	}

	accessUsername, err := types.NewUsername(payload.Username)
	if err != nil {
		return nil, err
	}

	accessPassword, err := types.NewPassword(payload.Password)
	if err != nil {
		return nil, err
	}

	accessName, err := types.NewName(payload.Name)
	if err != nil {
		return nil, err
	}

	return &models.Access{
		ID:        accessID,
		Username:  accessUsername,
		Password:  accessPassword,
		Name:      accessName,
		CreatedAt: payload.CreatedAt,
		UpdatedAt: payload.UpdatedAt,
	}, nil
}

func (r *PostgresAccessRepository) Update(access *models.Access) error {
	if access == nil {
		return errors.New("invalid access: access is nil")
	}
	if r.db == nil {
		return errors.New("database connection is nil")
	}
	if access.ID == nil ||
		access.Username == nil ||
		access.Password == nil ||
		access.Name == nil {
		return errors.New("invalid access: access is incomplete")
	}

	payload := accessRecord{
		ID:        access.ID.Value(),
		Username:  access.Username.Value(),
		Password:  access.Password.Value(),
		Name:      access.Name.Value(),
		CreatedAt: access.CreatedAt,
		UpdatedAt: access.UpdatedAt,
	}

	result := r.db.Model(&accessRecord{}).Where("id = ?", payload.ID).Updates(payload)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
