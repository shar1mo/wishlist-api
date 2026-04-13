package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secretKey []byte
	ttl       time.Duration
}

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func NewJWTManager(secret string, ttl time.Duration) *JWTManager {
	return &JWTManager{
		secretKey: []byte(secret),
		ttl:       ttl,
	}
}

func (m *JWTManager) Generate(userID int64) (string, error) {
	now := time.Now()

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", userID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(m.secretKey)
}

func (m *JWTManager) Parse(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return m.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}