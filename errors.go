package cache

import (
	"errors"
	"fmt"
)

// Sentinel errors wrapped by the errors returned from cache operations, for
// use with errors.Is. The Error() text is unchanged from previous versions.
var (
	ErrNotFound      = errors.New("item not found")
	ErrAlreadyExists = errors.New("item already exists")
	ErrOverflow      = errors.New("overflow would occur")
	ErrUnderflow     = errors.New("underflow would occur")
	ErrTypeMismatch  = errors.New("value does not have the expected type")
)

// keyedError keeps the historical message format while exposing a sentinel
// via Unwrap for errors.Is.
type keyedError struct {
	msg string
	err error
}

func (e keyedError) Error() string { return e.msg }
func (e keyedError) Unwrap() error { return e.err }

// keyErrorf returns the legacy formatted message for key k wrapping sentinel.
func keyErrorf(sentinel error, format, k string) error {
	return keyedError{msg: fmt.Sprintf(format, k), err: sentinel}
}
