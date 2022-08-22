package login

import (
	"context"

	repoDTO "github.com/olgoncharov/otbook/internal/repository/dto"
)

type (
	userRepo interface {
		GetPasswordHashByUsername(ctx context.Context, username string) (string, error)
		ReplaceRefreshToken(ctx context.Context, username string, newToken repoDTO.RefreshToken) error
	}

	passwordChecker interface {
		Check(input, hash string) (bool, error)
	}

	tokenGenerator interface {
		GenerateAccessToken(username string) (string, error)
		GenerateRefreshToken() string
	}

	config interface {
		JWTRefreshTokenTTL() uint64
	}
)
