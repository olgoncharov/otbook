package create

import (
	"context"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	repoDTO "github.com/olgoncharov/otbook/internal/repository/dto"
	repoErrors "github.com/olgoncharov/otbook/internal/repository/errors"
	"github.com/olgoncharov/otbook/internal/usecase/profile/command/create/mocks"
	"github.com/stretchr/testify/assert"
)

type testHandler struct {
	*Handler

	repoMock   *mocks.ProfileRepoMock
	hasherMock *mocks.PasswordHasherMock
}

func newTestHandler(mc *minimock.Controller) *testHandler {
	repoMock := mocks.NewProfileRepoMock(mc)
	hasherMock := mocks.NewPasswordHasherMock(mc)

	return &testHandler{
		Handler: NewHandler(repoMock, hasherMock),

		repoMock:   repoMock,
		hasherMock: hasherMock,
	}
}

func TestHandler_Handle(t *testing.T) {
	var (
		ctx = context.Background()

		command = Command{
			Username:  "test_user",
			Password:  "raw-password",
			FirstName: "Ivan",
			LastName:  "Ivanov",
			Birthdate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			City:      "Moscow",
			Sex:       "Male",
			Hobby:     "Guitar",
		}

		hashedPassword = "hashed-password"
	)

	t.Run("normal case", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		h := newTestHandler(mc)
		h.hasherMock.HashMock.Expect(command.Password).Return(hashedPassword, nil)
		h.repoMock.CreateProfileMock.Inspect(func(_ context.Context, regInfo repoDTO.RegistrationInfo) {
			assert.Equal(t, command.Username, regInfo.Username)
			assert.Equal(t, command.FirstName, regInfo.FirstName)
			assert.Equal(t, command.LastName, regInfo.LastName)
			assert.Equal(t, command.Birthdate, regInfo.Birthdate)
			assert.Equal(t, command.City, regInfo.City)
			assert.Equal(t, command.Sex, regInfo.Sex)
			assert.Equal(t, command.Hobby, regInfo.Hobby)
			assert.Equal(t, hashedPassword, regInfo.Password)
		}).Return(nil)

		err := h.Handle(ctx, command)
		assert.NoError(t, err)
	})

	t.Run("password hasher returns error", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		h := newTestHandler(mc)
		h.hasherMock.HashMock.Expect(command.Password).Return("", assert.AnError)

		err := h.Handle(ctx, command)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("user with given username already exists", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		h := newTestHandler(mc)
		h.hasherMock.HashMock.Expect(command.Password).Return(hashedPassword, nil)
		h.repoMock.CreateProfileMock.Inspect(func(_ context.Context, regInfo repoDTO.RegistrationInfo) {
			assert.Equal(t, command.Username, regInfo.Username)
			assert.Equal(t, command.FirstName, regInfo.FirstName)
			assert.Equal(t, command.LastName, regInfo.LastName)
			assert.Equal(t, command.Birthdate, regInfo.Birthdate)
			assert.Equal(t, command.City, regInfo.City)
			assert.Equal(t, command.Sex, regInfo.Sex)
			assert.Equal(t, command.Hobby, regInfo.Hobby)
			assert.Equal(t, hashedPassword, regInfo.Password)
		}).Return(repoErrors.ErrUniqueConstraintViolated)

		err := h.Handle(ctx, command)
		assert.ErrorIs(t, err, ErrUserAlreadyExists)
	})

	t.Run("repository returns error", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		h := newTestHandler(mc)
		h.hasherMock.HashMock.Expect(command.Password).Return(hashedPassword, nil)
		h.repoMock.CreateProfileMock.Inspect(func(_ context.Context, regInfo repoDTO.RegistrationInfo) {
			assert.Equal(t, command.Username, regInfo.Username)
			assert.Equal(t, command.FirstName, regInfo.FirstName)
			assert.Equal(t, command.LastName, regInfo.LastName)
			assert.Equal(t, command.Birthdate, regInfo.Birthdate)
			assert.Equal(t, command.City, regInfo.City)
			assert.Equal(t, command.Sex, regInfo.Sex)
			assert.Equal(t, command.Hobby, regInfo.Hobby)
			assert.Equal(t, hashedPassword, regInfo.Password)
		}).Return(assert.AnError)

		err := h.Handle(ctx, command)
		assert.ErrorIs(t, err, assert.AnError)
	})
}
