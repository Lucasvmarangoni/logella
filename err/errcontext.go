package errs

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type Error struct {
	Cause   error
	Code    int
	Message string
	added   bool
	context error
}

func (e *Error) Error() string {
	if e.Code == http.StatusInternalServerError {
		return e.Message
	}
	return e.Message + " : " + e.Cause.Error()
}

var Status = map[int]string{
	http.StatusBadRequest:          "BadRequest",
	http.StatusUnauthorized:        "Unauthorized",
	http.StatusForbidden:           "Forbidden",
	http.StatusInternalServerError: "InternalServerError",
}

func Wrap(cause error, contextValue string, code int) error {
	key := "Context"
	context := fmt.Errorf("%s: %s", key, contextValue)

	var e *Error
	if cause != nil {
		var ok bool
		if e, ok = cause.(*Error); !ok {
			e = &Error{
				Cause:   cause,
				Code:    code,
				added:   true,
				context: context,
			}
		} else {
			e.context = fmt.Errorf("%w %s: %s", e.context, key, contextValue)
			e.Cause = errors.Unwrap(cause)
		}
	}
	return e
}

func Assertion(err error) *Error {
	return err.(*Error)
}

func (e *Error) Msg(message string) {
	e.Message = message
}

func (e *Error) Context(operationValue string) error {
	key := "Context"

	e.context = fmt.Errorf("%w; %s: %s", e.context, key, operationValue)
	return e.Cause
}

func (e *Error) GetContext() error {
	if e != nil {
		return e.context
	}
	return nil
}

func (e *Error) Stack() error {
	if e != nil {
		if e.Code == 500 {
			return fmt.Errorf("%s | %s", "Internal Server Error", e.context)
		}
		return fmt.Errorf("%w | %s", e.Cause, e.context)
	}
	return nil
}
