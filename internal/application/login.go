package application

import (
	"errors"
	"gcpaccesss-svc-accesss/internal/domain/persistence"
	"gcpaccesss-svc-accesss/internal/domain/secure"
)

type LoginAccessUseCase struct {
	repository persistence.AccessRepository
	hasher     secure.Hasher
	tokener    secure.Tokener
}

func NewLoginAccessUseCase(repo persistence.AccessRepository, hasher secure.Hasher, tokener secure.Tokener) *LoginAccessUseCase {
	return &LoginAccessUseCase{repository: repo, hasher: hasher, tokener: tokener}
}

func (uc *LoginAccessUseCase) Login(username, password string) (*string, error) {
	access, err := uc.repository.GetFromUsername(username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := uc.hasher.Verify(access.Password.Value(), password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := uc.tokener.Generate(access.ID.Value())
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &token, nil
}
