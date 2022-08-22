package create

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/olgoncharov/otbook/internal/entity"
	repoDTO "github.com/olgoncharov/otbook/internal/repository/dto"
	repoErrors "github.com/olgoncharov/otbook/internal/repository/errors"
)

type (
	profileRepo interface {
		CreateProfile(ctx context.Context, registrationInfo repoDTO.RegistrationInfo) error
	}

	passwordHasher interface {
		Hash(password string) (string, error)
	}

	Handler struct {
		repo   profileRepo
		hasher passwordHasher
	}

	Command struct {
		Username  string
		Password  string
		FirstName string
		LastName  string
		Birthdate time.Time
		City      string
		Sex       string
		Hobby     string
	}
)

func NewHandler(repo profileRepo, hasher passwordHasher) *Handler {
	return &Handler{
		repo:   repo,
		hasher: hasher,
	}
}

func (h *Handler) Handle(ctx context.Context, command Command) error {
	passwordHash, err := h.hasher.Hash(command.Password)
	if err != nil {
		return fmt.Errorf("can't calculate password hash: %w", err)
	}

	err = h.repo.CreateProfile(ctx, repoDTO.RegistrationInfo{
		Profile: entity.Profile{
			Username:  command.Username,
			FirstName: command.FirstName,
			LastName:  command.LastName,
			Birthdate: command.Birthdate,
			City:      command.City,
			Sex:       command.Sex,
			Hobby:     command.Hobby,
		},
		Password: passwordHash,
	})

	if err == nil {
		return nil
	}

	if errors.Is(err, repoErrors.ErrUniqueConstraintViolated) {
		return ErrUserAlreadyExists
	}

	return fmt.Errorf("can't save profile to storage: %w", err)
}
