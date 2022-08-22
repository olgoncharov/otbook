package list

import (
	"context"
	"fmt"

	"github.com/olgoncharov/otbook/internal/entity"
)

type (
	friendsRepo interface {
		GetFriendsOfUser(ctx context.Context, username string, limit, offset uint) ([]entity.Profile, error)
		GetCountOfFriends(ctx context.Context, username string) (uint, error)
	}

	Handler struct {
		repo friendsRepo
	}

	Query struct {
		Username string
		Limit    uint
		Offset   uint
	}

	Result struct {
		Friends    []entity.Profile
		TotalCount uint
	}
)

func NewHandler(repo friendsRepo) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) Handle(ctx context.Context, command Query) (*Result, error) {
	friends, err := h.repo.GetFriendsOfUser(ctx, command.Username, command.Limit, command.Offset)
	if err != nil {
		return nil, fmt.Errorf("can't retrieve friends from storage: %w", err)
	}

	totalCount, err := h.repo.GetCountOfFriends(ctx, command.Username)
	if err != nil {
		return nil, fmt.Errorf("can't calculate total count of friends: %w", err)
	}

	return &Result{
		Friends:    friends,
		TotalCount: totalCount,
	}, nil
}
