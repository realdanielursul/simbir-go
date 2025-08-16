package service

import "errors"

var (
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrInvalidUsername       = errors.New("invalid username")
	ErrInvalidPassword       = errors.New("invalid password")
	ErrCannotSignToken       = errors.New("cannot sign token")
	ErrCannotParseToken      = errors.New("cannot parse token")
	ErrAccountNotFound       = errors.New("account not found")
	ErrTransportNotFound     = errors.New("transport not found")
	ErrAccessDenied          = errors.New("access denied")
)
