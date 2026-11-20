package persistence

import (
	"gcpaccesss-svc-accesss/internal/domain/models"
	"gcpaccesss-svc-accesss/internal/domain/types"
)

type AccessRepository interface {
	Save(access *models.Access) error
	Get(id *types.ID) (*models.Access, error)
	Delete(id *types.ID) error
	Update(access *models.Access) error
	List() ([]*models.Access, error)
	GetFromUsername(username string) (*models.Access, error)
}
