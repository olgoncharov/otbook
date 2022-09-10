package create

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user with given username already exists")
)
