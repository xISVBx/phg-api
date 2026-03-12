package services

import "github.com/google/uuid"

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

type Claims struct {
	UserID   uuid.UUID
	Username string
}

type IJWTService interface {
	Generate(userID uuid.UUID, username string) (*TokenPair, error)
	Parse(token string) (*Claims, error)
}

type IPasswordService interface {
	Hash(raw string) (string, error)
	Compare(hash, raw string) bool
	Validate(raw string) error
}
