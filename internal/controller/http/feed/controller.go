package feed

import (
	"context"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/olgoncharov/otbook/internal/controller/http/utils"
	"github.com/olgoncharov/otbook/internal/usecase/post/query/feed"
	"github.com/rs/zerolog"
)

type (
	useCase interface {
		Handle(ctx context.Context, query feed.Query) (*feed.Result, error)
	}
	linkBuilder interface {
		BuildProfileLink(username string) string
		BuildPostLink(postID uint64) string
	}

	Controller struct {
		useCase     useCase
		linkBuilder linkBuilder
		logger      zerolog.Logger
	}
)

func NewController(uCase useCase, lBuilder linkBuilder, logger zerolog.Logger) *Controller {
	return &Controller{
		useCase:     uCase,
		linkBuilder: lBuilder,
		logger:      logger,
	}
}

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	result, err := c.useCase.Handle(r.Context(), feed.Query{
		Username: utils.GetUsernameFromContext(r.Context()),
	})

	if err != nil {
		c.logger.Error().Err(err).Msg("")
		utils.WriteJSONError(w, utils.ErrInternal.Error(), http.StatusInternalServerError)

		return
	}

	resp := make([]postInfo, 0, len(result.Posts))
	for _, post := range result.Posts {
		resp = append(resp, postInfo{
			ID:        post.ID,
			Author:    post.Author,
			Title:     post.Title,
			CreatedAt: post.CreatedAt,
			Links: map[string]string{
				"self":   c.linkBuilder.BuildPostLink(post.ID),
				"author": c.linkBuilder.BuildProfileLink(post.Author),
			},
		})
	}

	responseEncoder := jsoniter.NewEncoder(w)
	responseEncoder.Encode(resp)
}
