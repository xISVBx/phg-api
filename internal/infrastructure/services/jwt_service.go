package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	appsvc "photogallery/api_go/internal/application/interfaces/services"
)

type JWTService struct {
	secret []byte
	exp    int64
}

type customClaims struct {
	UserID   string `json:"uid"`
	Username string `json:"uname"`
	jwt.RegisteredClaims
}

var _ appsvc.IJWTService = (*JWTService)(nil)

func NewJWTService(secret string, expSeconds int64) *JWTService {
	return &JWTService{secret: []byte(secret), exp: expSeconds}
}

func (s *JWTService) Generate(userID uuid.UUID, username string) (*appsvc.TokenPair, error) {
	now := time.Now().UTC()
	accessExp := now.Add(time.Duration(s.exp) * time.Second)
	refreshExp := now.Add(time.Duration(s.exp) * 24 * time.Second)

	accessClaims := customClaims{
		UserID:   userID.String(),
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExp),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   userID.String(),
			ID:        uuid.NewString(),
		},
	}
	accessJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	access, err := accessJWT.SignedString(s.secret)
	if err != nil {
		return nil, err
	}

	refreshClaims := customClaims{
		UserID:   userID.String(),
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExp),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   userID.String(),
			ID:        uuid.NewString(),
		},
	}
	refreshJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refresh, err := refreshJWT.SignedString(s.secret)
	if err != nil {
		return nil, err
	}

	return &appsvc.TokenPair{AccessToken: access, RefreshToken: refresh, ExpiresIn: s.exp}, nil
}

func (s *JWTService) Parse(token string) (*appsvc.Claims, error) {
	parsed, err := jwt.ParseWithClaims(token, &customClaims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("invalid signing method")
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := parsed.Claims.(*customClaims)
	if !ok || !parsed.Valid {
		return nil, errors.New("invalid token")
	}
	uid, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, err
	}
	return &appsvc.Claims{UserID: uid, Username: claims.Username}, nil
}
