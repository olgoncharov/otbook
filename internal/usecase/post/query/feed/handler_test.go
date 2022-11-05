package feed

import (
	"context"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/olgoncharov/otbook/internal/repository/dto"
	repoErrors "github.com/olgoncharov/otbook/internal/repository/errors"
	"github.com/olgoncharov/otbook/internal/usecase/post/query/feed/mocks"
	"github.com/stretchr/testify/assert"
)

type testHandler struct {
	*Handler

	cacheRepoMock    *mocks.CacheRepoMock
	durableRepoMock  *mocks.DurableRepoMock
	cacheUpdaterMock *mocks.CacheUpdaterMock
	configerMock     *mocks.ConfigerMock
}

func newTestHandler(mc *minimock.Controller) *testHandler {
	cacheRepoMock := mocks.NewCacheRepoMock(mc)
	durableRepoMock := mocks.NewDurableRepoMock(mc)
	cacheUpdaterMock := mocks.NewCacheUpdaterMock(mc)
	configerMock := mocks.NewConfigerMock(mc)

	return &testHandler{
		Handler: NewHandler(cacheRepoMock, durableRepoMock, cacheUpdaterMock, configerMock),

		cacheRepoMock:    cacheRepoMock,
		durableRepoMock:  durableRepoMock,
		cacheUpdaterMock: cacheUpdaterMock,
		configerMock:     configerMock,
	}
}

func TestHandler_Handle(t *testing.T) {
	const (
		username      = "test-user"
		postFeedLimit = 5
	)

	var (
		ctx = context.Background()
		now = time.Now().Truncate(time.Second)

		postsFixture = []dto.PostShortInfo{
			{ID: 100, CreatedAt: now},
			{ID: 99, CreatedAt: now.AddDate(0, 0, -1)},
			{ID: 98, CreatedAt: now.AddDate(0, 0, -2)},
			{ID: 97, CreatedAt: now.AddDate(0, 0, -3)},
			{ID: 96, CreatedAt: now.AddDate(0, 0, -4)},
			{ID: 95, CreatedAt: now.AddDate(0, 0, -5)},
			{ID: 94, CreatedAt: now.AddDate(0, 0, -6)},
			{ID: 93, CreatedAt: now.AddDate(0, 0, -7)},
		}

		celebrities = []string{
			"celebrity-1",
			"celebrity-2",
			"celebrity-3",
		}
	)

	t.Run("cache disabled", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		h := newTestHandler(mc)
		h.configerMock.IsFeedCacheDisabledMock.Return(true)
		h.configerMock.PostFeedLimitMock.Return(postFeedLimit)
		h.durableRepoMock.GetPostFeedMock.
			Expect(ctx, username, postFeedLimit).
			Return(postsFixture[:postFeedLimit], nil)

		actualResult, err := h.Handle(ctx, Query{Username: username})
		assert.NoError(t, err)
		assert.Equal(t, postsFixture[:postFeedLimit], actualResult.Posts)
	})

	t.Run("cache miss for feed", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		h := newTestHandler(mc)
		h.configerMock.IsFeedCacheDisabledMock.Return(false)
		h.configerMock.PostFeedLimitMock.Return(postFeedLimit)

		h.cacheRepoMock.GetPostFeedMock.
			Expect(ctx, username, postFeedLimit).
			Return(nil, repoErrors.ErrCacheMiss)

		h.cacheUpdaterMock.WarmupFeedMock.Expect(ctx, username).Return()
		h.cacheUpdaterMock.WarmupCelebrityFriendsListMock.Expect(ctx, username).Return()

		h.durableRepoMock.GetPostFeedMock.
			Expect(ctx, username, postFeedLimit).
			Return(postsFixture[:postFeedLimit], nil)

		actualResult, err := h.Handle(ctx, Query{Username: username})
		assert.NoError(t, err)
		assert.Equal(t, postsFixture[:postFeedLimit], actualResult.Posts)
	})

	t.Run("cache repo returns error", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		h := newTestHandler(mc)
		h.configerMock.IsFeedCacheDisabledMock.Return(false)
		h.configerMock.PostFeedLimitMock.Return(postFeedLimit)

		h.cacheRepoMock.GetPostFeedMock.
			Expect(ctx, username, postFeedLimit).
			Return(nil, assert.AnError)
		h.durableRepoMock.GetPostFeedMock.
			Expect(ctx, username, postFeedLimit).
			Return(postsFixture[:postFeedLimit], nil)

		actualResult, err := h.Handle(ctx, Query{Username: username})
		assert.NoError(t, err)
		assert.Equal(t, postsFixture[:postFeedLimit], actualResult.Posts)
	})

	t.Run("cache miss for celebrity friends list", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		h := newTestHandler(mc)
		h.configerMock.IsFeedCacheDisabledMock.Return(false)
		h.configerMock.PostFeedLimitMock.Return(postFeedLimit)

		h.cacheRepoMock.GetPostFeedMock.
			Expect(ctx, username, postFeedLimit).
			Return(postsFixture[:postFeedLimit], nil)

		h.cacheRepoMock.GetCelebrityFriendsMock.
			Expect(ctx, username).
			Return(nil, repoErrors.ErrCacheMiss)

		h.cacheUpdaterMock.WarmupCelebrityFriendsListMock.Expect(ctx, username).Return()

		h.durableRepoMock.GetCelebrityFriendsMock.
			Expect(ctx, username).
			Return(celebrities, nil)

		h.durableRepoMock.GetPostsByFiltersMock.
			Expect(ctx, dto.PostFilters{Authors: celebrities, DateFrom: &postsFixture[postFeedLimit-1].CreatedAt}, postFeedLimit, 0).
			Return(nil, nil)

		actualResult, err := h.Handle(ctx, Query{Username: username})
		assert.NoError(t, err)
		assert.Equal(t, postsFixture[:postFeedLimit], actualResult.Posts)
	})

	t.Run("empty cached posts and not empty celebreties posts", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		h := newTestHandler(mc)
		h.configerMock.IsFeedCacheDisabledMock.Return(false)
		h.configerMock.PostFeedLimitMock.Return(postFeedLimit)

		h.cacheRepoMock.GetPostFeedMock.
			Expect(ctx, username, postFeedLimit).
			Return(nil, nil)

		h.cacheRepoMock.GetCelebrityFriendsMock.
			Expect(ctx, username).
			Return(celebrities, nil)

		h.durableRepoMock.GetPostsByFiltersMock.
			Expect(ctx, dto.PostFilters{Authors: celebrities}, postFeedLimit, 0).
			Return(postsFixture[:postFeedLimit], nil)

		actualResult, err := h.Handle(ctx, Query{Username: username})
		assert.NoError(t, err)
		assert.Equal(t, postsFixture[:postFeedLimit], actualResult.Posts)
	})

	t.Run("not empty cached posts and not empty celebreties posts", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		h := newTestHandler(mc)
		h.configerMock.IsFeedCacheDisabledMock.Return(false)
		h.configerMock.PostFeedLimitMock.Return(postFeedLimit)

		h.cacheRepoMock.GetPostFeedMock.
			Expect(ctx, username, postFeedLimit).
			Return([]dto.PostShortInfo{
				postsFixture[0],
				postsFixture[2],
				postsFixture[4],
				postsFixture[6],
			}, nil)

		h.cacheRepoMock.GetCelebrityFriendsMock.
			Expect(ctx, username).
			Return(celebrities, nil)

		h.durableRepoMock.GetPostsByFiltersMock.
			Expect(ctx, dto.PostFilters{Authors: celebrities}, postFeedLimit, 0).
			Return([]dto.PostShortInfo{
				postsFixture[1],
				postsFixture[3],
				postsFixture[5],
				postsFixture[7],
			}, nil)

		actualResult, err := h.Handle(ctx, Query{Username: username})
		assert.NoError(t, err)
		assert.Equal(t, postsFixture[:postFeedLimit], actualResult.Posts)
	})
}

func TestHandler_MergePosts(t *testing.T) {
	var (
		now = time.Now().Truncate(time.Second)

		mc = minimock.NewController(t)
		h  = newTestHandler(mc)
	)

	h.configerMock.PostFeedLimitMock.Return(5)

	testCases := []struct {
		name            string
		postsFeed       []dto.PostShortInfo
		celebrityPosts  []dto.PostShortInfo
		wantedResultIDs []uint64
	}{
		{
			name: "no celebrity posts",
			postsFeed: []dto.PostShortInfo{
				{ID: 3, CreatedAt: now},
				{ID: 2, CreatedAt: now.AddDate(0, 0, -1)},
				{ID: 1, CreatedAt: now.AddDate(0, 0, -2)},
			},
			wantedResultIDs: []uint64{3, 2, 1},
		},
		{
			name: "empty input feed, not empty celebrity posts",
			celebrityPosts: []dto.PostShortInfo{
				{ID: 6, CreatedAt: now},
				{ID: 5, CreatedAt: now.AddDate(0, 0, -1)},
				{ID: 4, CreatedAt: now.AddDate(0, 0, -2)},
				{ID: 3, CreatedAt: now.AddDate(0, 0, -3)},
				{ID: 2, CreatedAt: now.AddDate(0, 0, -4)},
				{ID: 1, CreatedAt: now.AddDate(0, 0, -5)},
			},
			wantedResultIDs: []uint64{6, 5, 4, 3, 2},
		},
		{
			name: "input feed is full, there are no celebrity posts for inserting",
			postsFeed: []dto.PostShortInfo{
				{ID: 9, CreatedAt: now},
				{ID: 8, CreatedAt: now.AddDate(0, 0, -1)},
				{ID: 7, CreatedAt: now.AddDate(0, 0, -2)},
				{ID: 6, CreatedAt: now.AddDate(0, 0, -3)},
				{ID: 5, CreatedAt: now.AddDate(0, 0, -4)},
			},
			celebrityPosts: []dto.PostShortInfo{
				{ID: 4, CreatedAt: now.AddDate(0, 0, -5)},
				{ID: 3, CreatedAt: now.AddDate(0, 0, -6)},
				{ID: 2, CreatedAt: now.AddDate(0, 0, -7)},
				{ID: 1, CreatedAt: now.AddDate(0, 0, -8)},
			},
			wantedResultIDs: []uint64{9, 8, 7, 6, 5},
		},
		{
			name: "input feed is full, there are celebrity posts for inserting",
			postsFeed: []dto.PostShortInfo{
				{ID: 9, CreatedAt: now},
				{ID: 8, CreatedAt: now.AddDate(0, 0, -1)},
				{ID: 7, CreatedAt: now.AddDate(0, 0, -2)},
				{ID: 6, CreatedAt: now.AddDate(0, 0, -3)},
				{ID: 5, CreatedAt: now.AddDate(0, 0, -4)},
			},
			celebrityPosts: []dto.PostShortInfo{
				{ID: 4, CreatedAt: now.AddDate(0, 0, -1)},
				{ID: 3, CreatedAt: now.AddDate(0, 0, -2)},
				{ID: 2, CreatedAt: now.AddDate(0, 0, -3)},
				{ID: 1, CreatedAt: now.AddDate(0, 0, -4)},
			},
			wantedResultIDs: []uint64{9, 8, 4, 7, 3},
		},
		{
			name: "input feed is not full, there are celebrity posts for inserting",
			postsFeed: []dto.PostShortInfo{
				{ID: 7, CreatedAt: now},
				{ID: 6, CreatedAt: now.AddDate(0, 0, -1)},
				{ID: 5, CreatedAt: now.AddDate(0, 0, -2)},
			},
			celebrityPosts: []dto.PostShortInfo{
				{ID: 4, CreatedAt: now.AddDate(0, 0, -3)},
				{ID: 3, CreatedAt: now.AddDate(0, 0, -4)},
				{ID: 2, CreatedAt: now.AddDate(0, 0, -5)},
				{ID: 1, CreatedAt: now.AddDate(0, 0, -6)},
			},
			wantedResultIDs: []uint64{7, 6, 5, 4, 3},
		},
		{
			name: "input feed is totally replaced by celebrity posts",
			postsFeed: []dto.PostShortInfo{
				{ID: 5, CreatedAt: now.AddDate(0, 0, -6)},
				{ID: 4, CreatedAt: now.AddDate(0, 0, -7)},
				{ID: 3, CreatedAt: now.AddDate(0, 0, -8)},
				{ID: 2, CreatedAt: now.AddDate(0, 0, -9)},
				{ID: 1, CreatedAt: now.AddDate(0, 0, -10)},
			},
			celebrityPosts: []dto.PostShortInfo{
				{ID: 11, CreatedAt: now},
				{ID: 10, CreatedAt: now.AddDate(0, 0, -1)},
				{ID: 9, CreatedAt: now.AddDate(0, 0, -2)},
				{ID: 8, CreatedAt: now.AddDate(0, 0, -3)},
				{ID: 7, CreatedAt: now.AddDate(0, 0, -4)},
				{ID: 6, CreatedAt: now.AddDate(0, 0, -5)},
			},
			wantedResultIDs: []uint64{11, 10, 9, 8, 7},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mergedPosts := h.mergePosts(tc.postsFeed, tc.celebrityPosts)
			postIDs := make([]uint64, len(mergedPosts))
			for i, p := range mergedPosts {
				postIDs[i] = p.ID
			}

			assert.Equal(t, tc.wantedResultIDs, postIDs)
		})
	}
}
