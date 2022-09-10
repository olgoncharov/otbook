package user_profile

import (
	"context"
	"errors"
	"fmt"

	"github.com/olgoncharov/otbook/internal/entity"
	repoErrors "github.com/olgoncharov/otbook/internal/repository/errors"
)

type (
	profileRepo interface {
		GetProfileByUsername(ctx context.Context, username string) (*entity.Profile, error)
	}

	Handler struct {
		repo profileRepo
	}

	Query struct {
		Username string
	}
)

func NewHandler(repo profileRepo) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) Handle(ctx context.Context, query Query) (*entity.Profile, error) {
	profile, err := h.repo.GetProfileByUsername(ctx, query.Username)

	if errors.Is(err, repoErrors.ErrNoRowsFound) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("can't retrieve profile from storage: %w", err)
	}

	return profile, nil
}
