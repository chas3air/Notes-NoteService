package storageerrors

import "errors"

var (
	ErrNotFound      = errors.New("tip not found")
	ErrAlreadyExists = errors.New("tip already exists")
)
