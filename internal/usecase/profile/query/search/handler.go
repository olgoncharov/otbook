package search

import (
	"context"
	"fmt"

	"github.com/olgoncharov/otbook/internal/entity"
)

type (
	profileRepo interface {
		SearchProfiles(ctx context.Context, firstNamePrefix, lastNamePrefix string) ([]entity.Profile, error)
	}

	Handler struct {
		repo profileRepo
	}

	Query struct {
		FirstNamePrefix string
		LastNamePrefix  string
	}
)

func NewHandler(repo profileRepo) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) Handle(ctx context.Context, query Query) ([]entity.Profile, error) {
	profiles, err := h.repo.SearchProfiles(ctx, query.FirstNamePrefix, query.LastNamePrefix)

	if err != nil {
		return nil, fmt.Errorf("can't retrieve objects from storage: %w", err)
	}

	return profiles, nil
}
