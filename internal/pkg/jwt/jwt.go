package jwt

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	letters         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	refreshTokenLen = 32
)

type (
	TokenGenerator struct {
		cfg   generateConfig
		nowFn func() time.Time
	}

	TokenValidator struct {
		cfg validateConfig
	}

	claims struct {
		jwt.RegisteredClaims
		Username string `json:"username"`
	}

	generateConfig interface {
		JWTAccessTokenTTL() uint64
		JWTSigningKey() []byte
	}

	validateConfig interface {
		JWTSigningKey() []byte
	}
)

func NewTokenGenerator(cfg generateConfig, nowFn func() time.Time) *TokenGenerator {
	return &TokenGenerator{
		cfg:   cfg,
		nowFn: nowFn,
	}
}

func (g *TokenGenerator) GenerateAccessToken(username string) (string, error) {
	now := g.nowFn()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Second * time.Duration(g.cfg.JWTAccessTokenTTL()))),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		Username: username,
	})

	return token.SignedString(g.cfg.JWTSigningKey())
}

func (g *TokenGenerator) GenerateRefreshToken() string {
	rand.Seed(g.nowFn().UnixNano())
	token := make([]byte, refreshTokenLen)
	for i := range token {
		token[i] = letters[rand.Intn(len(letters))]
	}

	return string(token)
}

func NewTokenValidator(cfg validateConfig) *TokenValidator {
	return &TokenValidator{
		cfg: cfg,
	}
}

func (v *TokenValidator) ValidateAccessToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return v.cfg.JWTSigningKey(), nil
	})

	if err != nil {
		return "", err
	}

	if cl, ok := token.Claims.(*claims); ok {
		return cl.Username, cl.Valid()
	}

	return "", errors.New("invalid token")
}
