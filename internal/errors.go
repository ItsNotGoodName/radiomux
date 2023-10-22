package internal

import (
	"errors"
)

var (
	// ErrNotFound means the resource was not found.
	ErrNotFound = errors.New("not found")
)
