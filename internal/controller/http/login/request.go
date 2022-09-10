package login

import (
	"errors"
)

type requestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *requestBody) validate() error {
	if r.Username == "" {
		return errors.New("empty username")
	}

	if r.Password == "" {
		return errors.New("empty password")
	}

	return nil
}
