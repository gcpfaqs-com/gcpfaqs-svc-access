package types

import "errors"

var ErrInvalidUsername = errors.New("invalid username must be a non-empty string between 5 and 75 characters")

type Username struct {
	value string
}

func NewUsername(s string) (*Username, error) {
	if len(s) < 5 || len(s) > 75 {
		return nil, ErrInvalidUsername
	}

	return &Username{
		value: s,
	}, nil
}

func (t *Username) Value() string {
	return t.value
}
