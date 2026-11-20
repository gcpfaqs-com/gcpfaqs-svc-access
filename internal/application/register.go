package application

import (
	"gcpaccesss-svc-accesss/internal/domain/models"
	"gcpaccesss-svc-accesss/internal/domain/persistence"
	"gcpaccesss-svc-accesss/internal/domain/secure"
	"gcpaccesss-svc-accesss/internal/domain/types"
)

type RegisterAccessUseCase struct {
	repository persistence.AccessRepository
	hasher     secure.Hasher
	tokener    secure.Tokener
}

func NewRegisterAccessUseCase(
	repo persistence.AccessRepository,
	hasher secure.Hasher,
	tokener secure.Tokener) *RegisterAccessUseCase {

	return &RegisterAccessUseCase{
		repository: repo,
		hasher:     hasher,
		tokener:    tokener,
	}

}

func (uc *RegisterAccessUseCase) Register(username, password, name string) (string, error) {
	access, err := models.NewAccess(username, password, name)
	if err != nil {
		return "", err
	}

	hashedPassword, err := uc.hasher.Hash(password)
	if err != nil {
		return "", err
	}

	access.Password, err = types.NewPassword(hashedPassword)
	if err != nil {
		return "", err
	}

	err = uc.repository.Save(access)
	if err != nil {
		return "", err
	}

	return access.ID.Value(), nil
}
