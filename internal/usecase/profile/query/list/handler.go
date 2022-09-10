package list

import (
	"context"
	"fmt"

	"github.com/olgoncharov/otbook/internal/entity"
)

type (
	profileRepo interface {
		GetAllProfiles(ctx context.Context, limit, offset uint) ([]entity.Profile, error)
		GetProfilesCount(ctx context.Context) (uint, error)
	}

	Handler struct {
		repo profileRepo
	}

	Query struct {
		Limit  uint
		Offset uint
	}

	Result struct {
		Profiles   []entity.Profile
		TotalCount uint
	}
)

func NewHandler(repo profileRepo) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) Handle(ctx context.Context, query Query) (*Result, error) {
	profiles, err := h.repo.GetAllProfiles(ctx, query.Limit, query.Offset)

	if err != nil {
		return nil, fmt.Errorf("can't retrieve objects from storage: %w", err)
	}

	totalCount, err := h.repo.GetProfilesCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't calculate profiles count: %w", err)
	}

	return &Result{
		Profiles:   profiles,
		TotalCount: totalCount,
	}, nil
}
