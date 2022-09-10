package login

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
		repo       userRepo
		checker    passwordChecker
		tGenerator tokenGenerator
		cfg        config
		nowFn      func() time.Time
	}

	Command struct {
		Username string
		Password string
	}

	Result struct {
		AccessToken  string
		RefreshToken string
	}
)

func NewHandler(repo userRepo, checker passwordChecker, tGenerator tokenGenerator, cfg config, nowFn func() time.Time) *Handler {
	return &Handler{
		repo:       repo,
		checker:    checker,
		tGenerator: tGenerator,
		cfg:        cfg,
		nowFn:      nowFn,
	}
}

func (h *Handler) Handle(ctx context.Context, command Command) (*Result, error) {
	hash, err := h.repo.GetPasswordHashByUsername(ctx, command.Username)

	if errors.Is(err, repoErrors.ErrNoRowsFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("can't retrieve user info from storage: %w", err)
	}

	isValid, err := h.checker.Check(command.Password, hash)
	if err != nil {
		return nil, fmt.Errorf("can't check password: %w", err)
	}

	if !isValid {
		return nil, ErrInvalidPassword
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
