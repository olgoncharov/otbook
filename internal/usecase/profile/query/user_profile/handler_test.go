package user_profile

import (
	"context"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/olgoncharov/otbook/internal/entity"
	repoErrors "github.com/olgoncharov/otbook/internal/repository/errors"
	"github.com/olgoncharov/otbook/internal/usecase/profile/query/user_profile/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	const (
		username = "test-user"
	)

	var (
		ctx = context.Background()

		profile = &entity.Profile{
			Username:  username,
			FirstName: "Ivan",
			LastName:  "Ivanov",
			City:      "Moscow",
			Sex:       "Male",
			Hobby:     "Guitar",
		}
	)

	t.Run("normal case", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		h := newTestHandler(mc)
		h.repoMock.GetProfileByUsernameMock.Expect(ctx, username).Return(profile, nil)

		actualResult, err := h.Handle(ctx, Query{Username: username})

		require.NoError(t, err)
		assert.Equal(t, profile, actualResult)
	})

	t.Run("user with given username not found", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		h := newTestHandler(mc)
		h.repoMock.GetProfileByUsernameMock.Expect(ctx, username).Return(nil, repoErrors.ErrNoRowsFound)

		actualResult, err := h.Handle(ctx, Query{Username: username})
		assert.ErrorIs(t, err, ErrUserNotFound)
		assert.Nil(t, actualResult)
	})

	t.Run("repo returns error", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		h := newTestHandler(mc)
		h.repoMock.GetProfileByUsernameMock.Expect(ctx, username).Return(nil, assert.AnError)

		actualResult, err := h.Handle(ctx, Query{Username: username})
		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, actualResult)
	})
}
