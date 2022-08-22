package list

import (
	"context"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/olgoncharov/otbook/internal/entity"
	"github.com/olgoncharov/otbook/internal/usecase/friends/query/list/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testHandler struct {
	*Handler

	repoMock *mocks.FriendsRepoMock
}

func newTestHandler(mc *minimock.Controller) *testHandler {
	friendsRepoMock := mocks.NewFriendsRepoMock(mc)

	return &testHandler{
		Handler: NewHandler(friendsRepoMock),

		repoMock: friendsRepoMock,
	}
}

func TestHandler_Handle(t *testing.T) {
	const (
		username   = "test-user"
		limit      = 50
		offset     = 500
		totalCount = 1000
	)

	var (
		ctx = context.Background()

		friends = []entity.Profile{
			{Username: "friend-1"},
			{Username: "friend-2"},
			{Username: "friend-3"},
			{Username: "friend-4"},
			{Username: "friend-5"},
			{Username: "friend-6"},
		}
	)

	t.Run("normal case", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		h := newTestHandler(mc)
		h.repoMock.
			GetFriendsOfUserMock.Expect(ctx, username, limit, offset).Return(friends, nil).
			GetCountOfFriendsMock.Expect(ctx, username).Return(totalCount, nil)

		actualResult, err := h.Handle(ctx, Query{
			Username: username,
			Limit:    limit,
			Offset:   offset,
		})

		require.NoError(t, err)
		assert.ElementsMatch(t, friends, actualResult.Friends)
		assert.EqualValues(t, totalCount, actualResult.TotalCount)
	})

	t.Run("can't retrieve friends from storage", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		h := newTestHandler(mc)
		h.repoMock.GetFriendsOfUserMock.Expect(ctx, username, limit, offset).Return(nil, assert.AnError)

		actualResult, err := h.Handle(ctx, Query{
			Username: username,
			Limit:    limit,
			Offset:   offset,
		})

		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, actualResult)
	})

	t.Run("can't calculate total count of friends", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		h := newTestHandler(mc)
		h.repoMock.
			GetFriendsOfUserMock.Expect(ctx, username, limit, offset).Return(friends, nil).
			GetCountOfFriendsMock.Expect(ctx, username).Return(0, assert.AnError)

		actualResult, err := h.Handle(ctx, Query{
			Username: username,
			Limit:    limit,
			Offset:   offset,
		})

		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, actualResult)
	})
}
