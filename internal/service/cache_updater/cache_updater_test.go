package cacheupdater

import (
	"context"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/olgoncharov/otbook/internal/entity"
	"github.com/olgoncharov/otbook/internal/repository/dto"
	"github.com/olgoncharov/otbook/internal/service/cache_updater/mocks"
)

type testCacheUpdater struct {
	*CacheUpdater

	durableRepoMock *mocks.DurableRepoMock
	cacheRepoMock   *mocks.CacheRepoMock
	configerMock    *mocks.ConfigerMock
}

func newTestCacheUpdater(mc *minimock.Controller) *testCacheUpdater {
	durableRepoMock := mocks.NewDurableRepoMock(mc)
	cacheRepoMock := mocks.NewCacheRepoMock(mc)
	configerMock := mocks.NewConfigerMock(mc)

	return &testCacheUpdater{
		CacheUpdater: NewCacheUpdater(durableRepoMock, cacheRepoMock, configerMock),

		durableRepoMock: durableRepoMock,
		cacheRepoMock:   cacheRepoMock,
		configerMock:    configerMock,
	}
}

func TestCacheUpdater(t *testing.T) {
	t.Parallel()

	const (
		username      = "test-user"
		postFeedLimit = 1000
	)

	var (
		ctx = context.Background()
	)

	t.Run("warmup post feed", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		updater := newTestCacheUpdater(mc)

		posts := []dto.PostShortInfo{
			{ID: 1},
			{ID: 2},
			{ID: 3},
		}

		updater.configerMock.PostFeedLimitMock.Return(postFeedLimit)
		updater.durableRepoMock.GetPostFeedWithoutCelebritiesMock.
			Expect(ctx, username, postFeedLimit).
			Return(posts, nil)

		updater.cacheRepoMock.SetPostFeedMock.Expect(ctx, username, posts).Return(nil)

		updater.Run(ctx)
		updater.WarmupFeed(ctx, username)
		time.Sleep(time.Second)
		updater.Stop()
	})

	t.Run("warmup friends", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		updater := newTestCacheUpdater(mc)

		celebrityFriends := []string{"friend-1", "friend-2", "friend-3"}
		updater.durableRepoMock.GetCelebrityFriendsMock.
			Expect(ctx, username).
			Return(celebrityFriends, nil)

		updater.cacheRepoMock.SetCelebrityFriendsMock.
			Expect(ctx, username, celebrityFriends).
			Return(nil)

		updater.Run(ctx)
		updater.WarmupCelebrityFriendsList(ctx, username)
		time.Sleep(time.Second)
		updater.Stop()
	})

	t.Run("push post from not celebrity", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		updater := newTestCacheUpdater(mc)

		post := entity.Post{
			ID:     1,
			Author: "test-author",
		}
		followers := []string{"follower-1", "follower-2", "follower-3"}

		updater.configerMock.PostFeedLimitMock.Return(postFeedLimit)
		updater.durableRepoMock.GetFollowersOfUserMock.
			Expect(ctx, post.Author).
			Return(followers, nil)

		updater.durableRepoMock.GetProfileByUsernameMock.
			Expect(ctx, post.Author).
			Return(&entity.Profile{Username: post.Author, IsCelebrity: false}, nil)

		updater.cacheRepoMock.PushPostToFeedMock.
			Expect(ctx, post, followers, postFeedLimit).
			Return(nil)

		updater.Run(ctx)
		updater.PushPostToFeed(ctx, post)
		time.Sleep(time.Second)
		updater.Stop()
	})

	t.Run("push post from celebrity", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		updater := newTestCacheUpdater(mc)

		post := entity.Post{
			ID:     1,
			Author: "test-author",
		}

		updater.durableRepoMock.GetProfileByUsernameMock.
			Expect(ctx, post.Author).
			Return(&entity.Profile{Username: post.Author, IsCelebrity: true}, nil)

		updater.Run(ctx)
		updater.PushPostToFeed(ctx, post)
		time.Sleep(time.Second)
		updater.Stop()
	})

	t.Run("add friend", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		updater := newTestCacheUpdater(mc)

		friendProfile := entity.Profile{
			Username:    "friend",
			IsCelebrity: true,
		}
		posts := []dto.PostShortInfo{
			{ID: 1},
			{ID: 2},
			{ID: 3},
		}

		updater.configerMock.PostFeedLimitMock.Return(postFeedLimit)
		updater.durableRepoMock.GetProfileByUsernameMock.
			Expect(ctx, friendProfile.Username).
			Return(&friendProfile, nil)

		updater.durableRepoMock.GetPostFeedWithoutCelebritiesMock.
			Expect(ctx, username, postFeedLimit).
			Return(posts, nil)

		updater.cacheRepoMock.AddCelebrityFriendMock.
			Expect(ctx, username, friendProfile.Username).
			Return(nil)

		updater.cacheRepoMock.SetPostFeedMock.Expect(ctx, username, posts).Return(nil)

		updater.Run(ctx)
		updater.AddFriend(ctx, username, friendProfile.Username)
		time.Sleep(time.Second)
		updater.Stop()
	})

	t.Run("delete friend", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		updater := newTestCacheUpdater(mc)

		friendProfile := entity.Profile{
			Username:    "friend",
			IsCelebrity: true,
		}

		feed := []dto.PostShortInfo{
			{ID: 1, Author: friendProfile.Username},
			{ID: 2, Author: "somebody"},
		}

		updater.configerMock.PostFeedLimitMock.Return(postFeedLimit)
		updater.durableRepoMock.GetProfileByUsernameMock.
			Expect(ctx, friendProfile.Username).
			Return(&friendProfile, nil)

		updater.cacheRepoMock.DeleteCelebrityFriendMock.
			Expect(ctx, username, friendProfile.Username).
			Return(nil)

		updater.cacheRepoMock.GetPostFeedMock.
			Expect(ctx, username, postFeedLimit).
			Return(feed, nil)

		updater.cacheRepoMock.DeletePostsFromFeedMock.
			Expect(ctx, username, feed[:1]).
			Return(nil)

		updater.Run(ctx)
		updater.DeleteFriend(ctx, username, friendProfile.Username)
		time.Sleep(time.Second)
		updater.Stop()
	})
}
