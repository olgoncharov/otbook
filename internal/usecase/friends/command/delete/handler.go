package delete

import (
	"context"
	"fmt"
)

type (
	usersRepo interface {
		CheckUsersExistence(ctx context.Context, usernames ...string) (map[string]bool, error)
	}

	friendsRepo interface {
		DeleteFriend(ctx context.Context, user, friend string) error
	}

	Handler struct {
		usersRepo   usersRepo
		friendsRepo friendsRepo
	}

	Command struct {
		User   string
		Friend string
	}
)

func NewHandler(uRepo usersRepo, fRepo friendsRepo) *Handler {
	return &Handler{
		usersRepo:   uRepo,
		friendsRepo: fRepo,
	}
}

func (h *Handler) Handle(ctx context.Context, command Command) error {
	usersExistence, err := h.usersRepo.CheckUsersExistence(ctx, command.User, command.Friend)
	if err != nil {
		return fmt.Errorf("can't check users existence: %w", err)
	}

	if !usersExistence[command.User] {
		return fmt.Errorf("%w: %s", ErrUserNotFound, command.User)
	}

	if !usersExistence[command.Friend] {
		return fmt.Errorf("%w: %s", ErrUserNotFound, command.Friend)
	}

	return h.friendsRepo.DeleteFriend(ctx, command.User, command.Friend)
}
