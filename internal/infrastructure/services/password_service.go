package services

import (
	"errors"
	appsvc "photogallery/api_go/internal/application/interfaces/services"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

type PasswordService struct{}

var _ appsvc.IPasswordService = (*PasswordService)(nil)

func NewPasswordService() *PasswordService { return &PasswordService{} }

func (p *PasswordService) Hash(raw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (p *PasswordService) Compare(hash, raw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(raw)) == nil
}

func (p *PasswordService) Validate(raw string) error {
	if len(raw) < 8 {
		return errors.New("la contraseña debe tener al menos 8 caracteres")
	}

	if len(raw) > 72 {
		return errors.New("la contraseña es demasiado larga")
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool

	for _, r := range raw {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case strings.ContainsRune(`!@#$%^&*()-_=+[]{}|;:,.<>?/\'"~`, r):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New("la contraseña debe contener al menos una mayúscula")
	}
	if !hasLower {
		return errors.New("la contraseña debe contener al menos una minúscula")
	}
	if !hasDigit {
		return errors.New("la contraseña debe contener al menos un número")
	}
	if !hasSpecial {
		return errors.New("la contraseña debe contener al menos un carácter especial")
	}

	return nil
}