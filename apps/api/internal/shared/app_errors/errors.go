package app_errors

import "errors"

var (
	ErrConcurrentModification = errors.New("concurrent modification")
	ErrDuplicateResource      = errors.New("duplicate resource")
)
