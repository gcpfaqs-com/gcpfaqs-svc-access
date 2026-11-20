package types

import "errors"

var ErrInvalidPassword = errors.New("invalid password must be a non-empty string between 5 and 60 characters")

type Password struct {
	value string
}

func NewPassword(s string) (*Password, error) {
	if len(s) < 5 || len(s) > 60 {
		return nil, ErrInvalidPassword
	}

	return &Password{
		value: s,
	}, nil
}

func (c *Password) Value() string {
	return c.value
}
