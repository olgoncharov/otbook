package becomefriends

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
		CreateFriends(ctx context.Context, firstUsername, secondUsername string) error
	}

	Handler struct {
		usersRepo   usersRepo
		friendsRepo friendsRepo
	}

	Command struct {
		FirstUser  string
		SecondUser string
	}
)

func NewHandler(uRepo usersRepo, fRepo friendsRepo) *Handler {
	return &Handler{
		usersRepo:   uRepo,
		friendsRepo: fRepo,
	}
}

func (h *Handler) Handle(ctx context.Context, command Command) error {
	if command.FirstUser == command.SecondUser {
		return ErrSelfFriendship
	}

	usersExistence, err := h.usersRepo.CheckUsersExistence(ctx, command.FirstUser, command.SecondUser)
	if err != nil {
		return fmt.Errorf("can't check users existence: %w", err)
	}

	if !usersExistence[command.FirstUser] {
		return fmt.Errorf("%w: %s", ErrUserNotFound, command.FirstUser)
	}

	if !usersExistence[command.SecondUser] {
		return fmt.Errorf("%w: %s", ErrUserNotFound, command.SecondUser)
	}

	err = h.friendsRepo.CreateFriends(ctx, command.FirstUser, command.SecondUser)
	if errors.Is(err, repoErrors.ErrUniqueConstraintViolated) {
		return ErrAlreadyFriends
	}

	return err
}
