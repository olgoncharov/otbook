package list

import (
	"context"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/olgoncharov/otbook/internal/entity"
	"github.com/olgoncharov/otbook/internal/usecase/profile/query/list/mocks"
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
		limit      = 50
		offset     = 100
		totalCount = 1000
	)

	var (
		ctx = context.Background()

		profiles = []entity.Profile{
			{Username: "test-user-1"},
			{Username: "test-user-2"},
			{Username: "test-user-3"},
			{Username: "test-user-4"},
			{Username: "test-user-5"},
			{Username: "test-user-6"},
		}
	)

	t.Run("normal case", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		h := newTestHandler(mc)
		h.repoMock.
			GetAllProfilesMock.Expect(ctx, limit, offset).Return(profiles, nil).
			GetProfilesCountMock.Expect(ctx).Return(totalCount, nil)

		actualResult, err := h.Handle(ctx, Query{Limit: limit, Offset: offset})
		require.NoError(t, err)
		assert.ElementsMatch(t, profiles, actualResult.Profiles)
	})

	t.Run("GetAllProfiles returns error", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		h := newTestHandler(mc)
		h.repoMock.GetAllProfilesMock.Expect(ctx, limit, offset).Return(nil, assert.AnError)

		actualResult, err := h.Handle(ctx, Query{Limit: limit, Offset: offset})
		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, actualResult)
	})

	t.Run("GetProfilesCount returns error", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		h := newTestHandler(mc)
		h.repoMock.
			GetAllProfilesMock.Expect(ctx, limit, offset).Return(profiles, nil).
			GetProfilesCountMock.Expect(ctx).Return(0, assert.AnError)

		actualResult, err := h.Handle(ctx, Query{Limit: limit, Offset: offset})
		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, actualResult)
	})
}
