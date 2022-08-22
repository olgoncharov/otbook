package refreshtoken

import (
	"context"
	"errors"
	"fmt"
	"time"

	repoDTO "github.com/olgoncharov/otbook/internal/repository/dto"
	repoErrors "github.com/olgoncharov/otbook/internal/repository/errors"
)

type (
	Handler struct {
		repo       tokenRepo
		tGenerator tokenGenerator
		cfg        config
		nowFn      func() time.Time
	}

	Command struct {
		Username     string
		RefreshToken string
	}

	Result struct {
		AccessToken  string
		RefreshToken string
	}
)

func NewHandler(repo tokenRepo, tGenerator tokenGenerator, cfg config, nowFn func() time.Time) *Handler {
	return &Handler{
		repo:       repo,
		tGenerator: tGenerator,
		cfg:        cfg,
		nowFn:      nowFn,
	}
}

func (h *Handler) Handle(ctx context.Context, command Command) (*Result, error) {
	tokenInfo, err := h.repo.GetRefreshTokenForUser(ctx, command.Username)
	if errors.Is(err, repoErrors.ErrNoRowsFound) {
		return nil, ErrInvalidTokenGiven
	}

	if err != nil {
		return nil, fmt.Errorf("can't retrieve refresh token from storage: %w", err)
	}

	if tokenInfo.Value != command.RefreshToken {
		// Token mismatching mean it might have been stolen. Remove all user tokens by security reasons
		_ = h.repo.DeleteRefreshTokenForUser(ctx, command.Username)
		return nil, ErrInvalidTokenGiven
	}

	expirationDate := tokenInfo.CreatedAt.Add(time.Duration(tokenInfo.TTL) * time.Second)
	if expirationDate.Before(h.nowFn()) {
		_ = h.repo.DeleteRefreshTokenForUser(ctx, command.Username)
		return nil, ErrExpiredTokenGiven
	}

	accessToken, err := h.tGenerator.GenerateAccessToken(command.Username)
	if err != nil {
		return nil, fmt.Errorf("can't generate access token: %w", err)
	}

	refreshToken := h.tGenerator.GenerateRefreshToken()

	err = h.repo.ReplaceRefreshToken(ctx, command.Username, repoDTO.RefreshToken{
		Value:     refreshToken,
		CreatedAt: h.nowFn(),
		TTL:       h.cfg.JWTRefreshTokenTTL(),
	})
	if err != nil {
		return nil, fmt.Errorf("can't save user session in storage: %w", err)
	}

	return &Result{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
