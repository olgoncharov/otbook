package hash

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type HashChecker struct{}

func NewHashChecker() *HashChecker {
	return &HashChecker{}
}

func (c *HashChecker) Check(input, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(input))

	if err == nil {
		return true, nil
	}

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}

	return false, err
}
