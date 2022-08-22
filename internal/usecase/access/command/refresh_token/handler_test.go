package refreshtoken

import (
	"context"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	repoDTO "github.com/olgoncharov/otbook/internal/repository/dto"
	"github.com/olgoncharov/otbook/internal/usecase/access/command/refresh_token/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testHandler struct {
	*Handler

	repoMock           *mocks.TokenRepoMock
	tokenGeneratorMock *mocks.TokenGeneratorMock
	configMock         *mocks.ConfigMock
}

func newTestHandler(mc *minimock.Controller, now time.Time) *testHandler {
	repoMock := mocks.NewTokenRepoMock(mc)
	tokenGeneratorMock := mocks.NewTokenGeneratorMock(mc)
	configMock := mocks.NewConfigMock(mc)
	nowFn := func() time.Time {
		return now
	}

	return &testHandler{
		Handler: NewHandler(repoMock, tokenGeneratorMock, configMock, nowFn),

		repoMock:           repoMock,
		tokenGeneratorMock: tokenGeneratorMock,
		configMock:         configMock,
	}
}

func TestHandler_Handle(t *testing.T) {
	const (
		username        = "test-user"
		oldRefreshToken = "test-old-refresh-token"
		newRefreshToken = "test-new-refresh-token"
		newAccessToken  = "test-new-access-token"
		refreshTokenTTL = 360000
	)

	var (
		ctx                 = context.Background()
		now                 = time.Now()
		notExpiredTokenTime = now.Add(time.Duration(-refreshTokenTTL/2) * time.Second)
		expiredTokenTime    = now.Add(time.Duration(-refreshTokenTTL*2) * time.Second)

		command = Command{
			Username:     username,
			RefreshToken: oldRefreshToken,
		}
	)

	t.Run("given valid refresh token; expect ok", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		handler := newTestHandler(mc, now)

		handler.repoMock.GetRefreshTokenForUserMock.Expect(ctx, username).Return(
			&repoDTO.RefreshToken{
				Value:     oldRefreshToken,
				CreatedAt: notExpiredTokenTime,
				TTL:       refreshTokenTTL,
			}, nil,
		)
		handler.repoMock.ReplaceRefreshTokenMock.Expect(ctx, username, repoDTO.RefreshToken{Value: newRefreshToken, CreatedAt: now, TTL: refreshTokenTTL}).Return(nil)

		handler.tokenGeneratorMock.GenerateAccessTokenMock.Expect(username).Return(newAccessToken, nil)
		handler.tokenGeneratorMock.GenerateRefreshTokenMock.Return(newRefreshToken)
		handler.configMock.JWTRefreshTokenTTLMock.Return(refreshTokenTTL)

		actualResult, err := handler.Handle(ctx, command)

		require.NoError(t, err)
		assert.Equal(t, newAccessToken, actualResult.AccessToken)
		assert.Equal(t, newRefreshToken, actualResult.RefreshToken)
	})

	t.Run("given token doesn't match with stored; expect error", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		handler := newTestHandler(mc, now)

		handler.repoMock.GetRefreshTokenForUserMock.Expect(ctx, username).Return(
			&repoDTO.RefreshToken{
				Value:     "some-different-token",
				CreatedAt: notExpiredTokenTime,
				TTL:       refreshTokenTTL,
			}, nil,
		)
		handler.repoMock.DeleteRefreshTokenForUserMock.Expect(ctx, username).Return(nil)

		actualResult, err := handler.Handle(ctx, command)

		assert.ErrorIs(t, err, ErrInvalidTokenGiven)
		assert.Nil(t, actualResult)
	})

	t.Run("given token is expired; expect error", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		handler := newTestHandler(mc, now)

		handler.repoMock.GetRefreshTokenForUserMock.Expect(ctx, username).Return(
			&repoDTO.RefreshToken{
				Value:     oldRefreshToken,
				CreatedAt: expiredTokenTime,
				TTL:       refreshTokenTTL,
			}, nil,
		)
		handler.repoMock.DeleteRefreshTokenForUserMock.Expect(ctx, username).Return(nil)

		actualResult, err := handler.Handle(ctx, command)

		assert.ErrorIs(t, err, ErrExpiredTokenGiven)
		assert.Nil(t, actualResult)
	})

	t.Run("can't retrieve refresh token from storage; expect error", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		handler := newTestHandler(mc, now)
		handler.repoMock.GetRefreshTokenForUserMock.Expect(ctx, username).Return(nil, assert.AnError)

		actualResult, err := handler.Handle(ctx, command)

		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, actualResult)
	})

	t.Run("can't save user session in storage; expect error", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		handler := newTestHandler(mc, now)

		handler.repoMock.GetRefreshTokenForUserMock.Expect(ctx, username).Return(
			&repoDTO.RefreshToken{
				Value:     oldRefreshToken,
				CreatedAt: notExpiredTokenTime,
				TTL:       refreshTokenTTL,
			}, nil,
		)
		handler.repoMock.ReplaceRefreshTokenMock.Expect(ctx, username, repoDTO.RefreshToken{Value: newRefreshToken, CreatedAt: now, TTL: refreshTokenTTL}).Return(assert.AnError)

		handler.tokenGeneratorMock.GenerateAccessTokenMock.Expect(username).Return(newAccessToken, nil)
		handler.tokenGeneratorMock.GenerateRefreshTokenMock.Return(newRefreshToken)
		handler.configMock.JWTRefreshTokenTTLMock.Return(refreshTokenTTL)

		actualResult, err := handler.Handle(ctx, command)

		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, actualResult)
	})
}
