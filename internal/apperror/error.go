package apperror

import "errors"

type Kind string

const (
	Validation Kind = "validation"
	NotFound   Kind = "not_found"
	Internal   Kind = "internal"
)

var (
	ErrInvalidJSON      = New(Validation, "invalid JSON", nil)
	ErrInvalidMonitorID = New(Validation, "monitor id must be a number", nil)
)

type Error struct {
	Kind    Kind
	Message string
	Err     error
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func New(kind Kind, message string, err error) *Error {
	return &Error{Kind: kind, Message: message, Err: err}
}

func KindOf(err error) Kind {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr.Kind
	}
	return Internal
}
