package singlepost

import (
	"context"
	"errors"
	"fmt"

	"github.com/olgoncharov/otbook/internal/entity"
	repoErrors "github.com/olgoncharov/otbook/internal/repository/errors"
)

type (
	postRepo interface {
		GetPostByID(ctx context.Context, id uint64) (*entity.Post, error)
	}

	Handler struct {
		repo postRepo
	}

	Query struct {
		ID uint64
	}
)

func NewHandler(repo postRepo) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) Handle(ctx context.Context, query Query) (*entity.Post, error) {
	post, err := h.repo.GetPostByID(ctx, query.ID)
	if errors.Is(err, repoErrors.ErrNoRowsFound) {
		return nil, ErrPostNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("can't retrieve post from storage: %w", err)
	}

	return post, nil
}
