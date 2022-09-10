package update

import (
	"context"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/olgoncharov/otbook/internal/entity"
	"github.com/olgoncharov/otbook/internal/usecase/profile/command/update/mocks"
	"github.com/stretchr/testify/assert"
)

type testHandler struct {
	*Handler

	repoMock *mocks.ProfileRepoMock
}

func newTestHandler(mc *minimock.Controller) *testHandler {
	repoMock := mocks.NewProfileRepoMock(mc)

	return &testHandler{
		Handler: NewHandler(repoMock),

		repoMock: repoMock,
	}
}

func TestHandler_Handle(t *testing.T) {
	var (
		ctx = context.Background()

		command = Command{
			Username:  "test_user",
			FirstName: "Ivan",
			LastName:  "Ivanov",
			Birthdate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			City:      "Moscow",
			Sex:       "Male",
			Hobby:     "Guitar",
		}
	)

	t.Run("normal case", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		h := newTestHandler(mc)
		h.repoMock.CheckUsersExistenceMock.Return(map[string]bool{command.Username: true}, nil)
		h.repoMock.UpdateProfileMock.Inspect(func(_ context.Context, profile entity.Profile) {
			assert.Equal(t, command.Username, profile.Username)
			assert.Equal(t, command.FirstName, profile.FirstName)
			assert.Equal(t, command.LastName, profile.LastName)
			assert.Equal(t, command.Birthdate, profile.Birthdate)
			assert.Equal(t, command.City, profile.City)
			assert.Equal(t, command.Sex, profile.Sex)
			assert.Equal(t, command.Hobby, profile.Hobby)
		}).Return(nil)
		err := h.Handle(ctx, command)
		assert.NoError(t, err)
	})

	t.Run("user not found", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		h := newTestHandler(mc)
		h.repoMock.CheckUsersExistenceMock.Return(map[string]bool{command.Username: false}, nil)

		err := h.Handle(ctx, command)
		assert.ErrorIs(t, err, ErrUserNotFound)
	})

	t.Run("repo returns error", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		h := newTestHandler(mc)
		h.repoMock.CheckUsersExistenceMock.Return(map[string]bool{command.Username: true}, nil)
		h.repoMock.UpdateProfileMock.Inspect(func(_ context.Context, profile entity.Profile) {
			assert.Equal(t, command.Username, profile.Username)
			assert.Equal(t, command.FirstName, profile.FirstName)
			assert.Equal(t, command.LastName, profile.LastName)
			assert.Equal(t, command.Birthdate, profile.Birthdate)
			assert.Equal(t, command.City, profile.City)
			assert.Equal(t, command.Sex, profile.Sex)
			assert.Equal(t, command.Hobby, profile.Hobby)
		}).Return(assert.AnError)

		err := h.Handle(ctx, command)
		assert.ErrorIs(t, err, assert.AnError)
	})
}
