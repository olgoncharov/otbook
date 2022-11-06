package full_list

import (
	"context"
	"fmt"

	"github.com/olgoncharov/otbook/internal/repository/dto"
)

type (
	postRepo interface {
		GetAllPosts(ctx context.Context, limit, offset uint) ([]dto.PostShortInfo, error)
		GetPostsCount(ctx context.Context) (uint, error)
	}

	Handler struct {
		repo postRepo
	}

	Query struct {
		Limit  uint
		Offset uint
	}

	Result struct {
		Posts      []dto.PostShortInfo
		TotalCount uint
	}
)

func NewHandler(repo postRepo) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) Handle(ctx context.Context, query Query) (*Result, error) {
	posts, err := h.repo.GetAllPosts(ctx, query.Limit, query.Offset)
	if err != nil {
		return nil, fmt.Errorf("can't retrieve posts from storage: %w", err)
	}

	totalCount, err := h.repo.GetPostsCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get posts count: %w", err)
	}

	return &Result{
		Posts:      posts,
		TotalCount: totalCount,
	}, nil
}
