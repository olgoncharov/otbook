package create

import (
	"context"
	"fmt"
	"time"

	"github.com/olgoncharov/otbook/internal/entity"
)

type (
	postRepo interface {
		CreatePost(ctx context.Context, post entity.Post) (uint64, error)
	}

	cacheUpdater interface {
		PushPostToFeed(ctx context.Context, post entity.Post)
	}

	Handler struct {
		repo         postRepo
		cacheUpdater cacheUpdater
		nowFn        func() time.Time
	}

	Command struct {
		Author string
		Title  string
		Text   string
	}
)

func NewHandler(repo postRepo, cUpdater cacheUpdater, nowFn func() time.Time) *Handler {
	return &Handler{
		repo:         repo,
		cacheUpdater: cUpdater,
		nowFn:        nowFn,
	}
}

func (h *Handler) Handle(ctx context.Context, command Command) (uint64, error) {
	newPost := entity.Post{
		Author:    command.Author,
		Title:     command.Title,
		Text:      command.Text,
		CreatedAt: h.nowFn(),
	}

	postID, err := h.repo.CreatePost(ctx, newPost)
	if err != nil {
		return 0, fmt.Errorf("can't create post: %w", err)
	}

	newPost.ID = postID
	h.cacheUpdater.PushPostToFeed(ctx, newPost)

	return postID, nil
}
