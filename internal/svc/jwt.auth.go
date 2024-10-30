package svc

import (
	"fmt"
	"irule-api/internal/config"
	"irule-api/internal/db/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserClaims struct {
	UserID         uuid.UUID `json:"user_id"`
	Role           string    `json:"role"`
	OrganizationId uuid.UUID `json:"org_id"`
	jwt.RegisteredClaims
}

func CreateToken(user *models.User, cfg *config.Config) (string, error) {
	claims := &UserClaims{
		UserID:         user.ID,
		Role:           user.Role,
		OrganizationId: user.OrganizationId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		claims)

	tokenString, err := token.SignedString([]byte(cfg.JWT_SECRET))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string, cfg *config.Config) (*UserClaims, error) {
	claims := &UserClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT_SECRET), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
