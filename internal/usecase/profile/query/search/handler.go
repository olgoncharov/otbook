package search

import (
	"context"
	"fmt"

	"github.com/olgoncharov/otbook/internal/repository/dto"
)

type (
	profileRepo interface {
		SearchProfiles(ctx context.Context, firstNamePrefix, lastNamePrefix string) ([]dto.ProfileShortInfo, error)
	}

	Handler struct {
		repo profileRepo
	}

	Query struct {
		FirstNamePrefix string
		LastNamePrefix  string
	}

	ProfileInfo struct {
		Username  string
		FirstName string
		LastName  string
	}
)

func NewHandler(repo profileRepo) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) Handle(ctx context.Context, query Query) ([]ProfileInfo, error) {
	profiles, err := h.repo.SearchProfiles(ctx, query.FirstNamePrefix, query.LastNamePrefix)

	if err != nil {
		return nil, fmt.Errorf("can't retrieve objects from storage: %w", err)
	}

	result := make([]ProfileInfo, len(profiles))
	for i, p := range profiles {
		result[i] = ProfileInfo{
			Username:  p.Username,
			FirstName: p.FirstName,
			LastName:  p.LastName,
		}
	}

	return result, nil
}
