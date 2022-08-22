package user_profile

import "errors"

var (
	ErrUserNotFound = errors.New("user with given username not found")
)
