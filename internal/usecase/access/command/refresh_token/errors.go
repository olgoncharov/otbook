package refreshtoken

import "errors"

var (
	ErrInvalidTokenGiven = errors.New("invalid token given")
	ErrExpiredTokenGiven = errors.New("expired token given")
)
