package create

import (
	"context"
	"time"

	"github.com/olgoncharov/otbook/internal/entity"
)

type (
	postRepo interface {
		CreatePost(ctx context.Context, post entity.Post) (uint64, error)
	}

	Handler struct {
		repo  postRepo
		nowFn func() time.Time
	}

	Command struct {
		Author string
		Title  string
		Text   string
	}
)

func NewHandler(repo postRepo, nowFn func() time.Time) *Handler {
	return &Handler{
		repo:  repo,
		nowFn: nowFn,
	}
}

func (h *Handler) Handle(ctx context.Context, command Command) (uint64, error) {
	return h.repo.CreatePost(ctx, entity.Post{
		Author:    command.Author,
		Title:     command.Title,
		Text:      command.Text,
		CreatedAt: h.nowFn(),
	})
}
