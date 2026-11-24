package services

import (
	"gcpaccesss-svc-accesss/internal/domain/persistence"
	"gcpaccesss-svc-accesss/internal/domain/types"
	"time"
)

type UpdateAccessUseCase struct {
	repository persistence.AccessRepository
}

func NewUpdateAccessUseCase(repo persistence.AccessRepository) *UpdateAccessUseCase {

	return &UpdateAccessUseCase{repository: repo}

}

func (uc *UpdateAccessUseCase) Do(id, name string) error {
	validID, err := types.NewIDFromString(id)
	if err != nil {
		return err
	}

	access, err := uc.repository.Get(validID)
	if err != nil {
		return err
	}

	validName, err := types.NewName(name)
	if err != nil {
		return err
	}

	access.Name = validName
	access.UpdatedAt = time.Now()
	if err := uc.repository.Update(access); err == nil {
		return err
	}

	return nil
}
