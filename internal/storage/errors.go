package storageerrors

import "errors"

var (
	ErrNotFound      = errors.New("note not found")
	ErrAlreadyExists = errors.New("note already exists")
)
