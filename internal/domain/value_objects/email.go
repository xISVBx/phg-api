package valueobjects

import (
	"errors"
	"net/mail"
	"strings"
)

type Email struct {
	value string
}

func NewEmail(raw string) (Email, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return Email{}, errors.New("el email es requerido")
	}

	addr, err := mail.ParseAddress(raw)
	if err != nil {
		return Email{}, errors.New("el email no es válido")
	}

	if addr.Address != raw {
		return Email{}, errors.New("el email no es válido")
	}

	return Email{value: raw}, nil
}

func (e Email) String() string {
	return e.value
}
