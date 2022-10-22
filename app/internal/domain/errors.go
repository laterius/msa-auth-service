package domain

import "errors"

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrInvalidUserId = errors.New("invalid user id")
	InternalError    = errors.New("internal error")
)
