package feed

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/olgoncharov/otbook/internal/repository/dto"
	repoErrors "github.com/olgoncharov/otbook/internal/repository/errors"
)

type (
	cacheRepo interface {
		GetPostFeed(ctx context.Context, username string, limit uint) ([]dto.PostShortInfo, error)
		GetCelebrityFriends(ctx context.Context, username string) ([]string, error)
	}

	durableRepo interface {
		GetPostFeed(ctx context.Context, username string, limit uint) ([]dto.PostShortInfo, error)
		GetPostsByFilters(ctx context.Context, filters dto.PostFilters, limit uint, offset uint) ([]dto.PostShortInfo, error)
		GetCelebrityFriends(ctx context.Context, username string) ([]string, error)
	}

	cacheUpdater interface {
		WarmupFeed(ctx context.Context, username string)
		WarmupCelebrityFriendsList(ctx context.Context, username string)
	}

	configer interface {
		IsFeedCacheDisabled() bool
		PostFeedLimit() int
	}

	Handler struct {
		cacheRepo    cacheRepo
		durableRepo  durableRepo
		cacheUpdater cacheUpdater
		cfg          configer
	}

	Query struct {
		Username string
	}

	Result struct {
		Posts []dto.PostShortInfo
	}
)

func NewHandler(cRepo cacheRepo, dRepo durableRepo, cUpdater cacheUpdater, cfg configer) *Handler {
	return &Handler{
		cacheRepo:    cRepo,
		durableRepo:  dRepo,
		cacheUpdater: cUpdater,
		cfg:          cfg,
	}
}

func (h *Handler) Handle(ctx context.Context, query Query) (*Result, error) {
	if h.cfg.IsFeedCacheDisabled() {
		return h.getFeedFromDurableStorage(ctx, query.Username)
	}

	cachedPosts, err := h.cacheRepo.GetPostFeed(ctx, query.Username, uint(h.cfg.PostFeedLimit()))
	if errors.Is(err, repoErrors.ErrCacheMiss) {
		h.cacheUpdater.WarmupFeed(ctx, query.Username)
		h.cacheUpdater.WarmupCelebrityFriendsList(ctx, query.Username)
	}
	if err != nil {
		return h.getFeedFromDurableStorage(ctx, query.Username)
	}

	celebrityFriends, err := h.getCelebrityFriends(ctx, query.Username)
	if err != nil {
		return nil, fmt.Errorf("can't get celebrity friends of user: %w", err)
	}

	var celebrityPosts []dto.PostShortInfo
	if len(celebrityFriends) > 0 {
		filters := dto.PostFilters{
			Authors: celebrityFriends,
		}
		if len(cachedPosts) == int(h.cfg.PostFeedLimit()) {
			filters.DateFrom = &cachedPosts[len(cachedPosts)-1].CreatedAt
		}

		celebrityPosts, err = h.durableRepo.GetPostsByFilters(ctx, filters, uint(h.cfg.PostFeedLimit()), 0)
		if err != nil {
			return nil, fmt.Errorf("can't retrieve celebrity posts from storage: %w", err)
		}
	}

	return &Result{
		Posts: h.mergePosts(cachedPosts, celebrityPosts),
	}, nil
}

func (h *Handler) getFeedFromDurableStorage(ctx context.Context, username string) (*Result, error) {
	posts, err := h.durableRepo.GetPostFeed(ctx, username, uint(h.cfg.PostFeedLimit()))
	if err != nil {
		return nil, fmt.Errorf("can't retrieve posts from storage: %w", err)
	}

	return &Result{
		Posts: posts,
	}, nil
}

func (h *Handler) getCelebrityFriends(ctx context.Context, username string) ([]string, error) {
	friends, err := h.cacheRepo.GetCelebrityFriends(ctx, username)
	if errors.Is(err, repoErrors.ErrCacheMiss) {
		h.cacheUpdater.WarmupCelebrityFriendsList(ctx, username)
	}

	if err == nil {
		return friends, nil
	}

	friends, err = h.durableRepo.GetCelebrityFriends(ctx, username)
	if err != nil {
		return nil, err
	}

	return friends, nil
}

func (h *Handler) mergePosts(postFeed, celebrityPosts []dto.PostShortInfo) []dto.PostShortInfo {
	if len(celebrityPosts) == 0 {
		return postFeed
	}

	if len(postFeed) == 0 {
		if len(celebrityPosts) > h.cfg.PostFeedLimit() {
			return celebrityPosts[:h.cfg.PostFeedLimit()]
		}
		return celebrityPosts
	}

	minPostDate := postFeed[len(postFeed)-1].CreatedAt
	for _, celebrityPost := range celebrityPosts {
		if celebrityPost.CreatedAt.After(minPostDate) ||
			celebrityPost.CreatedAt.Equal(minPostDate) ||
			len(postFeed) < h.cfg.PostFeedLimit() {

			postFeed = append(postFeed, celebrityPost)
			continue
		}

		break
	}

	sort.Slice(postFeed, func(i, j int) bool {
		if postFeed[i].CreatedAt.Equal(postFeed[j].CreatedAt) {
			return postFeed[i].ID > postFeed[j].ID
		}
		return postFeed[i].CreatedAt.After(postFeed[j].CreatedAt)
	})

	if len(postFeed) > h.cfg.PostFeedLimit() {
		return postFeed[:h.cfg.PostFeedLimit()]
	}

	return postFeed
}
