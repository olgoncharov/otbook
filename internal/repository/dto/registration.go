package dto

import (
	"github.com/olgoncharov/otbook/internal/entity"
)

type RegistrationInfo struct {
	entity.Profile

	Password string
}
