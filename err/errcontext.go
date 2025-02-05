package errs

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type Error struct {
	Cause      error
	Code       int
	Message    string
	added      bool
	operations error
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

func Wrap(cause error, operationValue string, code int) error {
	key := "Operation"
	operation := fmt.Errorf("%s: %s", key, operationValue)

	var e *Error
	if cause != nil {
		var ok bool
		if e, ok = cause.(*Error); !ok {
			e = &Error{
				Cause:      cause,
				Code:       code,
				added:      true,
				operations: operation,
			}
		} else {
			e.operations = fmt.Errorf("%w; %s: %s", e.operations, key, operationValue)
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

func (e *Error) GetOperations() error {
	if e != nil {
		return e.operations
	}
	return nil
}

func (e *Error) Stack() error {
	if e != nil {
		return fmt.Errorf("%w | %s", e.Cause, e.operations)
	}
	return nil
}
