package add

import (
	"context"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/olgoncharov/otbook/internal/usecase/friends/command/add/mocks"
	"github.com/stretchr/testify/assert"
)

type testHandler struct {
	*Handler

	usersRepoMock   *mocks.UsersRepoMock
	friendsRepoMock *mocks.FriendsRepoMock
}

func newTestHandler(mc *minimock.Controller) *testHandler {
	usersRepoMock := mocks.NewUsersRepoMock(mc)
	friendsRepoMock := mocks.NewFriendsRepoMock(mc)

	return &testHandler{
		Handler: NewHandler(usersRepoMock, friendsRepoMock),

		usersRepoMock:   usersRepoMock,
		friendsRepoMock: friendsRepoMock,
	}
}

func TestHandler_Handle(t *testing.T) {
	const (
		username1 = "test-user-1"
		username2 = "test-user-2"
	)

	var (
		ctx = context.Background()
	)

	t.Run("normal case", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		h := newTestHandler(mc)
		h.usersRepoMock.CheckUsersExistenceMock.Expect(ctx, username1, username2).Return(
			map[string]bool{
				username1: true,
				username2: true,
			}, nil)
		h.friendsRepoMock.AddFriendMock.Expect(ctx, username1, username2).Return(nil)

		err := h.Handle(ctx, Command{User: username1, NewFriend: username2})
		assert.NoError(t, err)
	})

	t.Run("first user not found", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		h := newTestHandler(mc)
		h.usersRepoMock.CheckUsersExistenceMock.Expect(ctx, username1, username2).Return(
			map[string]bool{
				username1: false,
				username2: true,
			}, nil)

		err := h.Handle(ctx, Command{User: username1, NewFriend: username2})
		assert.ErrorIs(t, err, ErrUserNotFound)
	})

	t.Run("second user not found", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		h := newTestHandler(mc)
		h.usersRepoMock.CheckUsersExistenceMock.Expect(ctx, username1, username2).Return(
			map[string]bool{
				username1: true,
				username2: false,
			}, nil)

		err := h.Handle(ctx, Command{User: username1, NewFriend: username2})
		assert.ErrorIs(t, err, ErrUserNotFound)
	})

	t.Run("user repo return error", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		h := newTestHandler(mc)
		h.usersRepoMock.CheckUsersExistenceMock.Expect(ctx, username1, username2).Return(nil, assert.AnError)

		err := h.Handle(ctx, Command{User: username1, NewFriend: username2})
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("friends repo returns error", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		h := newTestHandler(mc)
		h.usersRepoMock.CheckUsersExistenceMock.Expect(ctx, username1, username2).Return(
			map[string]bool{
				username1: true,
				username2: true,
			}, nil)
		h.friendsRepoMock.AddFriendMock.Expect(ctx, username1, username2).Return(assert.AnError)

		err := h.Handle(ctx, Command{User: username1, NewFriend: username2})
		assert.ErrorIs(t, err, assert.AnError)
	})
}
