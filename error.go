package clickup

import "errors"

// Custom errors.
var (
	ErrSignatureMismatch = errors.New("Signature mismatch")
	ErrStatusNotUpdated  = errors.New("Task status was not updated")
)
