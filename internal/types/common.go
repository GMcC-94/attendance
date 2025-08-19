package types

import "errors"

var (
	ErrNotFound = errors.New("resource not found")
	ErrConflict = errors.New("resource conflict")
)
