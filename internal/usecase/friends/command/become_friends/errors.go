package becomefriends

import "errors"

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrSelfFriendship = errors.New("you cannot add oneself to friends")
	ErrAlreadyFriends = errors.New("users are friends already")
)
