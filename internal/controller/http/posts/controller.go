package posts

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
	"github.com/olgoncharov/otbook/internal/controller/http/utils"
	"github.com/olgoncharov/otbook/internal/entity"
	"github.com/olgoncharov/otbook/internal/usecase/post/command/create"
	"github.com/olgoncharov/otbook/internal/usecase/post/query/full_list"
	singlepost "github.com/olgoncharov/otbook/internal/usecase/post/query/single_post"
	"github.com/rs/zerolog"
)

type (
	createPostUseCase interface {
		Handle(ctx context.Context, command create.Command) (uint64, error)
	}
	getPostsListUseCase interface {
		Handle(ctx context.Context, query full_list.Query) (*full_list.Result, error)
	}
	getPostUseCase interface {
		Handle(ctx context.Context, query singlepost.Query) (*entity.Post, error)
	}

	linkBuilder interface {
		BuildProfileLink(username string) string
		BuildPostLink(postID uint64) string
	}

	ListController struct {
		getUseCase    getPostsListUseCase
		createUseCase createPostUseCase
		linkBuilder   linkBuilder
		logger        zerolog.Logger
	}

	ObjectController struct {
		getUseCase  getPostUseCase
		linkBuilder linkBuilder
		logger      zerolog.Logger
	}
)

func NewListController(
	getUseCase getPostsListUseCase,
	createUseCase createPostUseCase,
	lBuilder linkBuilder,
	logger zerolog.Logger,
) *ListController {
	return &ListController{
		getUseCase:    getUseCase,
		createUseCase: createUseCase,
		linkBuilder:   lBuilder,
		logger:        logger,
	}
}

func NewObjectController(getUseCase getPostUseCase, lBuilder linkBuilder, logger zerolog.Logger) *ObjectController {
	return &ObjectController{
		getUseCase:  getUseCase,
		linkBuilder: lBuilder,
		logger:      logger,
	}
}

func (c *ListController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c.serveGet(w, r)
	case http.MethodPost:
		c.servePost(w, r)
	}
}

func (c *ListController) serveGet(w http.ResponseWriter, r *http.Request) {
	limit, offset := utils.GetLimitOffsetFromURL(r.URL)

	result, err := c.getUseCase.Handle(r.Context(), full_list.Query{
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		c.logger.Error().Err(err).Msg("")
		utils.WriteJSONError(w, utils.ErrInternal.Error(), http.StatusInternalServerError)

		return
	}

	resp := listResponse{
		List:       make([]postInfo, len(result.Posts)),
		TotalCount: result.TotalCount,
	}

	for i, p := range result.Posts {
		resp.List[i] = postInfo{
			ID:        p.ID,
			Author:    p.Author,
			Title:     p.Title,
			CreatedAt: p.CreatedAt,
			Links: map[string]string{
				"self":   c.linkBuilder.BuildPostLink(p.ID),
				"author": c.linkBuilder.BuildProfileLink(p.Author),
			},
		}
	}

	responseEncoder := jsoniter.NewEncoder(w)
	responseEncoder.Encode(resp)
}

func (c *ListController) servePost(w http.ResponseWriter, r *http.Request) {
	bodyDecoder := jsoniter.NewDecoder(r.Body)
	body := &createRequest{}
	err := bodyDecoder.Decode(body)
	if err != nil {
		c.logger.Error().Err(err).Msg("")
		utils.WriteJSONError(w, "invalid request body", http.StatusBadRequest)

		return
	}

	err = body.validate()
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	postID, err := c.createUseCase.Handle(r.Context(), create.Command{
		Author: utils.GetUsernameFromContext(r.Context()),
		Title:  body.Title,
		Text:   body.Text,
	})

	if err != nil {
		c.logger.Error().Err(err).Msg("")
		utils.WriteJSONError(w, utils.ErrInternal.Error(), http.StatusInternalServerError)

		return
	}

	responseEncoder := jsoniter.NewEncoder(w)
	responseEncoder.Encode(createResponse{ID: postID})
}

func (c *ObjectController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqVars := mux.Vars(r)

	postID, err := strconv.ParseUint(reqVars["id"], 10, 64)
	if err != nil {
		c.logger.Error().Err(err).Msg("")
		utils.WriteJSONError(w, "invalid post id", http.StatusBadRequest)

		return
	}

	post, err := c.getUseCase.Handle(r.Context(), singlepost.Query{ID: postID})
	if errors.Is(err, singlepost.ErrPostNotFound) {
		utils.WriteJSONError(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		c.logger.Error().Err(err).Msg("")
		utils.WriteJSONError(w, utils.ErrInternal.Error(), http.StatusInternalServerError)

		return
	}

	responseEncoder := jsoniter.NewEncoder(w)
	responseEncoder.Encode(singlePostResponse{
		ID:        post.ID,
		Author:    post.Author,
		Title:     post.Title,
		Text:      post.Text,
		CreatedAt: post.CreatedAt,
		Links: map[string]string{
			"author": c.linkBuilder.BuildProfileLink(post.Author),
		},
	})
}
