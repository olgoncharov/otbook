package errors

import "errors"

var (
	ErrNoRowsFound              = errors.New("no rows found")
	ErrUniqueConstraintViolated = errors.New("unique constraint violated")
	ErrCacheMiss                = errors.New("cache miss")
)
