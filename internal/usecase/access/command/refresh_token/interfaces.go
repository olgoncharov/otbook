package refreshtoken

import (
	"context"

	repoDTO "github.com/olgoncharov/otbook/internal/repository/dto"
)

type (
	tokenRepo interface {
		GetRefreshTokenForUser(ctx context.Context, username string) (*repoDTO.RefreshToken, error)
		DeleteRefreshTokenForUser(ctx context.Context, username string) error
		ReplaceRefreshToken(ctx context.Context, username string, newToken repoDTO.RefreshToken) error
	}

	tokenGenerator interface {
		GenerateAccessToken(username string) (string, error)
		GenerateRefreshToken() string
	}

	config interface {
		JWTRefreshTokenTTL() uint64
	}
)
