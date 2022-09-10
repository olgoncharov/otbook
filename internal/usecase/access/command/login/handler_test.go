package login

import (
	"context"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	repoDTO "github.com/olgoncharov/otbook/internal/repository/dto"
	repoErrors "github.com/olgoncharov/otbook/internal/repository/errors"
	"github.com/olgoncharov/otbook/internal/usecase/access/command/login/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testHandler struct {
	*Handler

	repoMock            *mocks.UserRepoMock
	passwordCheckerMock *mocks.PasswordCheckerMock
	tokenGeneratorMock  *mocks.TokenGeneratorMock
	configMock          *mocks.ConfigMock
}

func newTestHandler(mc *minimock.Controller, now time.Time) *testHandler {
	repoMock := mocks.NewUserRepoMock(mc)
	passwordCheckerMock := mocks.NewPasswordCheckerMock(mc)
	tokenGeneratorMock := mocks.NewTokenGeneratorMock(mc)
	configMock := mocks.NewConfigMock(mc)
	nowFn := func() time.Time {
		return now
	}

	return &testHandler{
		Handler: NewHandler(repoMock, passwordCheckerMock, tokenGeneratorMock, configMock, nowFn),

		repoMock:            repoMock,
		passwordCheckerMock: passwordCheckerMock,
		tokenGeneratorMock:  tokenGeneratorMock,
		configMock:          configMock,
	}
}

func TestHandler_Handle(t *testing.T) {
	const (
		username        = "test-user"
		rawPassword     = "test-password"
		hashedPassword  = "test-password-hash"
		accessToken     = "test-access-token"
		refreshToken    = "test-refresh-token"
		refreshTokenTTL = 360000
	)

	var (
		ctx = context.Background()
		now = time.Now()

		command = Command{
			Username: username,
			Password: rawPassword,
		}
	)

	t.Run("given valid username and password; expect ok", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		handler := newTestHandler(mc, now)

		handler.repoMock.
			GetPasswordHashByUsernameMock.Expect(ctx, username).Return(hashedPassword, nil).
			ReplaceRefreshTokenMock.Expect(ctx, username, repoDTO.RefreshToken{Value: refreshToken, CreatedAt: now, TTL: refreshTokenTTL}).Return(nil)

		handler.passwordCheckerMock.CheckMock.
			Expect(rawPassword, hashedPassword).
			Return(true, nil)

		handler.tokenGeneratorMock.GenerateAccessTokenMock.Expect(username).Return(accessToken, nil)
		handler.tokenGeneratorMock.GenerateRefreshTokenMock.Return(refreshToken)
		handler.configMock.JWTRefreshTokenTTLMock.Return(refreshTokenTTL)

		actualResult, err := handler.Handle(ctx, command)

		require.NoError(t, err)
		assert.Equal(t, accessToken, actualResult.AccessToken)
		assert.Equal(t, refreshToken, actualResult.RefreshToken)
	})

	t.Run("user not found by given username; expect error", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		handler := newTestHandler(mc, now)
		handler.repoMock.GetPasswordHashByUsernameMock.Expect(ctx, username).Return("", repoErrors.ErrNoRowsFound)

		actualResult, err := handler.Handle(ctx, command)

		assert.ErrorIs(t, err, ErrUserNotFound)
		assert.Nil(t, actualResult)
	})

	t.Run("user info retrieving failed; expect error", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		handler := newTestHandler(mc, now)
		handler.repoMock.GetPasswordHashByUsernameMock.Expect(ctx, username).Return("", assert.AnError)

		actualResult, err := handler.Handle(ctx, command)

		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, actualResult)
	})

	t.Run("given password is not valid; expect error", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		handler := newTestHandler(mc, now)

		handler.repoMock.GetPasswordHashByUsernameMock.Expect(ctx, username).Return(hashedPassword, nil)
		handler.passwordCheckerMock.CheckMock.
			Expect(rawPassword, hashedPassword).
			Return(false, nil)

		actualResult, err := handler.Handle(ctx, command)
		assert.ErrorIs(t, err, ErrInvalidPassword)
		assert.Nil(t, actualResult)
	})

	t.Run("password checking failed; expect error", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		handler := newTestHandler(mc, now)

		handler.repoMock.GetPasswordHashByUsernameMock.Expect(ctx, username).Return(hashedPassword, nil)
		handler.passwordCheckerMock.CheckMock.
			Expect(rawPassword, hashedPassword).
			Return(false, assert.AnError)

		actualResult, err := handler.Handle(ctx, command)
		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, actualResult)
	})

	t.Run("refresh token can't be saved in storage; expect error", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		handler := newTestHandler(mc, now)

		handler.repoMock.
			GetPasswordHashByUsernameMock.Expect(ctx, username).Return(hashedPassword, nil).
			ReplaceRefreshTokenMock.Expect(ctx, username, repoDTO.RefreshToken{Value: refreshToken, CreatedAt: now, TTL: refreshTokenTTL}).Return(assert.AnError)

		handler.passwordCheckerMock.CheckMock.
			Expect(rawPassword, hashedPassword).
			Return(true, nil)

		handler.tokenGeneratorMock.GenerateAccessTokenMock.Expect(username).Return(accessToken, nil)
		handler.tokenGeneratorMock.GenerateRefreshTokenMock.Return(refreshToken)
		handler.configMock.JWTRefreshTokenTTLMock.Return(refreshTokenTTL)

		actualResult, err := handler.Handle(ctx, command)

		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, actualResult)
	})
}
