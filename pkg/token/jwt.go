package token

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/revandpratama/lognest/config"
)

type CustomClaims struct {
	UserID       string `json:"user_id"`
	Email        string `json:"email"`
	RoleID       uint   `json:"role_id"`
	Provider     string `json:"provider,omitempty"` // Optional if OAuth
	SessionID    string `json:"sid,omitempty"`      // Optional, for token tracking
	MFACompleted bool   `json:"mfa,omitempty"`
	jwt.RegisteredClaims
}

func ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(config.ENV.JWT_SECRET), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}

func ParseExpiredToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(config.ENV.JWT_SECRET), nil
	}, jwt.WithoutClaimsValidation())

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("unexpected signing method")
	}

	// Manually validate expiration if needed
	// For refresh flow, we expect it to be expired
	if claims.ExpiresAt == nil {
		return nil, errors.New("missing exp in token")
	}

	return claims, nil
}
