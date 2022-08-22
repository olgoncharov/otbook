package update

import (
	"context"
	"fmt"
	"time"

	"github.com/olgoncharov/otbook/internal/entity"
)

type (
	profileRepo interface {
		UpdateProfile(ctx context.Context, profile entity.Profile) error
		CheckUsersExistence(ctx context.Context, usernames ...string) (map[string]bool, error)
	}

	Handler struct {
		repo profileRepo
	}

	Command struct {
		Username  string
		FirstName string
		LastName  string
		Birthdate time.Time
		City      string
		Sex       string
		Hobby     string
	}
)

func NewHandler(repo profileRepo) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) Handle(ctx context.Context, command Command) error {
	userExistence, err := h.repo.CheckUsersExistence(ctx, command.Username)
	if err != nil {
		return fmt.Errorf("can't check users existence: %w", err)
	}

	if !userExistence[command.Username] {
		return ErrUserNotFound
	}

	err = h.repo.UpdateProfile(ctx, convertCommandToProfile(command))
	if err != nil {
		return err
	}

	return nil
}

func convertCommandToProfile(command Command) entity.Profile {
	return entity.Profile{
		Username:  command.Username,
		FirstName: command.FirstName,
		LastName:  command.LastName,
		Birthdate: command.Birthdate,
		City:      command.City,
		Sex:       command.Sex,
		Hobby:     command.Hobby,
	}
}
