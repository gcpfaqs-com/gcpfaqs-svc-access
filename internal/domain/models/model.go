package models

import "gcpaccesss-svc-accesss/internal/domain/types"

type Access struct {
	ID       *types.ID       `json:"id"`
	Username *types.Username `json:"username"`
	Password *types.Password `json:"password"`
	Name     *types.Name     `json:"name"`
}

func NewAccess(username string, password string, name string) (*Access, error) {
	accessUsername, err := types.NewUsername(username)
	if err != nil {
		return nil, err
	}

	accessPassword, err := types.NewPassword(password)
	if err != nil {
		return nil, err
	}

	accessName, err := types.NewName(name)
	if err != nil {
		return nil, err
	}

	accessID, err := types.NewID()
	if err != nil {
		return nil, err
	}

	return &Access{
		ID:       accessID,
		Username: accessUsername,
		Password: accessPassword,
		Name:     accessName,
	}, nil
}
