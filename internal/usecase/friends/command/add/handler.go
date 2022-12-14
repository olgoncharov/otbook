package add

import (
	"context"
	"errors"
	"fmt"

	repoErrors "github.com/olgoncharov/otbook/internal/repository/errors"
)

type (
	usersRepo interface {
		CheckUsersExistence(ctx context.Context, usernames ...string) (map[string]bool, error)
	}

	friendsRepo interface {
		AddFriend(ctx context.Context, user, newFriend string) error
	}

	cacheUpdater interface {
		AddFriend(ctx context.Context, user string, friend string)
	}

	Handler struct {
		usersRepo    usersRepo
		friendsRepo  friendsRepo
		cacheUpdater cacheUpdater
	}

	Command struct {
		User      string
		NewFriend string
	}
)

func NewHandler(uRepo usersRepo, fRepo friendsRepo, cUpdater cacheUpdater) *Handler {
	return &Handler{
		usersRepo:    uRepo,
		friendsRepo:  fRepo,
		cacheUpdater: cUpdater,
	}
}

func (h *Handler) Handle(ctx context.Context, command Command) error {
	if command.User == command.NewFriend {
		return ErrSelfFriendship
	}

	usersExistence, err := h.usersRepo.CheckUsersExistence(ctx, command.User, command.NewFriend)
	if err != nil {
		return fmt.Errorf("can't check users existence: %w", err)
	}

	if !usersExistence[command.User] {
		return fmt.Errorf("%w: %s", ErrUserNotFound, command.User)
	}

	if !usersExistence[command.NewFriend] {
		return fmt.Errorf("%w: %s", ErrUserNotFound, command.NewFriend)
	}

	err = h.friendsRepo.AddFriend(ctx, command.User, command.NewFriend)
	if errors.Is(err, repoErrors.ErrUniqueConstraintViolated) {
		return ErrAlreadyFriends
	}

	if err != nil {
		return err
	}

	h.cacheUpdater.AddFriend(ctx, command.User, command.NewFriend)

	return nil
}
