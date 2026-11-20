package types

import "errors"

var ErrInvalidName = errors.New("invalid Name must be a non-empty string between 5 and 75 characters")

type Name struct {
	value string
}

func NewName(s string) (*Name, error) {
	if len(s) < 5 || len(s) > 75 {
		return nil, ErrInvalidName
	}

	return &Name{
		value: s,
	}, nil
}

func (c *Name) Value() string {
	return c.value
}
