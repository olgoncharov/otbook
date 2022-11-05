package cacheupdater

import (
	"context"
	"os"
	"sync"

	"github.com/olgoncharov/otbook/internal/entity"
	"github.com/olgoncharov/otbook/internal/repository/dto"
	"github.com/rs/zerolog"
)

type (
	durableRepo interface {
		GetPostFeedWithoutCelebrities(ctx context.Context, username string, limit uint) ([]dto.PostShortInfo, error)
		GetCelebrityFriends(ctx context.Context, username string) ([]string, error)
		GetFollowersOfUser(ctx context.Context, username string) ([]string, error)
		GetProfileByUsername(ctx context.Context, username string) (*entity.Profile, error)
	}

	cacheRepo interface {
		GetPostFeed(ctx context.Context, username string, limit uint) ([]dto.PostShortInfo, error)
		SetPostFeed(ctx context.Context, username string, posts []dto.PostShortInfo) error
		SetCelebrityFriends(ctx context.Context, usercname string, friends []string) error
		PushPostToFeed(ctx context.Context, post entity.Post, usernames []string, feedLimit uint) error
		DeletePostsFromFeed(ctx context.Context, username string, posts []dto.PostShortInfo) error
		AddCelebrityFriend(ctx context.Context, user string, friend string) error
		DeleteCelebrityFriend(ctx context.Context, user string, friend string) error
	}

	configer interface {
		PostFeedLimit() int
	}

	updateFriendsMsg struct {
		user          string
		addedFriend   string
		deletedFriend string
	}

	CacheUpdater struct {
		durableRepo durableRepo
		cacheRepo   cacheRepo
		cfg         configer

		warmupFeedQueue    chan string
		warmupFriendsQueue chan string
		pushPostQueue      chan entity.Post
		updateFriendsQueue chan updateFriendsMsg
		stopCh             chan struct{}

		wg sync.WaitGroup

		logger zerolog.Logger
	}
)

func NewCacheUpdater(durRepo durableRepo, cRepo cacheRepo, cfg configer) *CacheUpdater {
	return &CacheUpdater{
		durableRepo: durRepo,
		cacheRepo:   cRepo,
		cfg:         cfg,

		warmupFeedQueue:    make(chan string, 10_000),
		warmupFriendsQueue: make(chan string, 10_000),
		pushPostQueue:      make(chan entity.Post, 10_000),
		updateFriendsQueue: make(chan updateFriendsMsg, 10_000),
		stopCh:             make(chan struct{}),

		logger: zerolog.New(os.Stdout).With().Timestamp().Str("component", "cacheUpdater").Logger(),
	}
}

func (c *CacheUpdater) Run(ctx context.Context) {
	c.runFeedWarmupWorker(ctx)
	c.runCelebrityFriendsWarmupWorker(ctx)
	c.runPostPushWorker(ctx)
	c.runFriendsUpdateWorker(ctx)
}

func (c *CacheUpdater) Stop() error {
	close(c.stopCh)

	c.wg.Wait()

	close(c.warmupFeedQueue)
	close(c.warmupFriendsQueue)
	close(c.pushPostQueue)
	close(c.updateFriendsQueue)

	return nil
}
