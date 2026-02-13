package domain

import "errors"

var (
	ErrAlreadyExists = errors.New("count value already exists")
	ErrNotFound      = errors.New("count value not found")
)
