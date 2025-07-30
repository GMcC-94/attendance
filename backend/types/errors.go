package types

import "errors"

var (
	ErrUsernameTaken = errors.New("types: username is already in use")
)
