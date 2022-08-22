package hash

import (
	"golang.org/x/crypto/bcrypt"
)

type HashGenerator struct {
	cost int
}

func NewHashGenerator(cost int) *HashGenerator {
	return &HashGenerator{
		cost: cost,
	}
}

func (h *HashGenerator) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
