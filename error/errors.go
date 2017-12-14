package error

import "errors"

// ErrInternalServer Internal Server Error
var (
	ErrInternalServer = errors.New("Internal Server Error")
	ErrConflict       = errors.New("Conflict Error")
	ErrNotFound       = errors.New("Not Found Error")
)
