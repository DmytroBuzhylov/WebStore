package service

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"time"
)

type Claims struct {
	UserID string `json:"sub"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type TokenGenerator interface {
	Generate(id, role, secret string) (string, error)
}

type AccessTokenGenerator struct{}

func (g *AccessTokenGenerator) Generate(id, role, secret string) (string, error) {
	claims := Claims{
		UserID: id,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "shop-api",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

type RefreshTokenGenerator struct{}

func (g *RefreshTokenGenerator) Generate(id, role, secret string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   id,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "shop-api",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func generateToken(generator TokenGenerator, id, role, secret string) (string, error) {
	return generator.Generate(id, role, secret)
}

type TokenChecker interface {
	Validate(token, secret string) (*Claims, error)
}

type AccessTokenChecker struct{}

func (c *AccessTokenChecker) Validate(token, secret string) (*Claims, error) {
	claims := &Claims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неверный алгоритм подписи")
		}
		return []byte(secret), nil
	})

	if err != nil || !parsedToken.Valid {
		return nil, errors.New("токен недействителен")
	}
	return claims, nil
}

type RefreshTokenChecker struct{}

func (c *RefreshTokenChecker) Validate(token, secret string) (*Claims, error) {
	return &Claims{}, nil
}

type StoreToken interface {
	StoreToken(token string, r *redis.Client) error
	IsValid(token string, r *redis.Client) bool
}

type StoreRefreshToken struct{}

func (t *StoreRefreshToken) StoreToken(token string, r *redis.Client) error {
	return r.Set(context.Background(), token, "active", 7*24*time.Hour).Err()
}

func (t *StoreRefreshToken) IsValid(token string, r *redis.Client) bool {
	val, err := r.Get(context.Background(), token).Result()
	return err == nil && val == "active"
}

func RefreshAccessToken(refreshToken string) {}
