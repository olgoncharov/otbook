package search

import (
	"context"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/olgoncharov/otbook/internal/repository/dto"
	"github.com/olgoncharov/otbook/internal/usecase/profile/query/search/mocks"
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
		firstNamePrefix = "Al"
		lastNamePrefix  = "Be"
	)

	var (
		ctx = context.Background()
	)

	t.Run("normal case", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		h := newTestHandler(mc)
		h.repoMock.SearchProfilesMock.
			Expect(ctx, firstNamePrefix, lastNamePrefix).
			Return([]dto.ProfileShortInfo{
				{Username: "test-user-1"},
				{Username: "test-user-2"},
				{Username: "test-user-3"},
				{Username: "test-user-4"},
				{Username: "test-user-5"},
				{Username: "test-user-6"},
			}, nil)

		actualResult, err := h.Handle(ctx, Query{
			FirstNamePrefix: firstNamePrefix,
			LastNamePrefix:  lastNamePrefix,
		})
		require.NoError(t, err)
		assert.ElementsMatch(t, []ProfileInfo{
			{Username: "test-user-1"},
			{Username: "test-user-2"},
			{Username: "test-user-3"},
			{Username: "test-user-4"},
			{Username: "test-user-5"},
			{Username: "test-user-6"},
		}, actualResult)
	})

	t.Run("repository returns error", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		h := newTestHandler(mc)
		h.repoMock.SearchProfilesMock.
			Expect(ctx, firstNamePrefix, lastNamePrefix).
			Return(nil, assert.AnError)

		actualResult, err := h.Handle(ctx, Query{
			FirstNamePrefix: firstNamePrefix,
			LastNamePrefix:  lastNamePrefix,
		})
		assert.Empty(t, actualResult)
		assert.ErrorIs(t, err, assert.AnError)
	})
}
