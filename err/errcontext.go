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
	trace   error
}

func (e *Error) Error() string {
	if e.Code == http.StatusInternalServerError {
		return e.Message
	}
	return e.Message + " : " + e.Cause.Error()
}

var Status = map[int]string{
	http.StatusBadRequest:           "BadRequest",
	http.StatusUnauthorized:         "Unauthorized",
	http.StatusForbidden:            "Forbidden",
	http.StatusNotFound:             "NotFound",
	http.StatusInternalServerError:  "InternalServerError",
	http.StatusUnsupportedMediaType: "UnsupportedMediaType",
	http.StatusNoContent:            "NoContent",
	http.StatusMovedPermanently:     "MovedPermanently",
	http.StatusTemporaryRedirect:    "TemporaryRedirect",
	http.StatusBadGateway:           "BadGateway",
	http.StatusServiceUnavailable:   "ServiceUnavailable",
}


func Wrap(cause error, traceValue string, code int) error {
	key := "Trace"
	trace := fmt.Errorf("%s: %s", key, traceValue)

	var e *Error
	if cause != nil {
		var ok bool
		if e, ok = cause.(*Error); !ok {
			e = &Error{
				Cause: cause,
				Code:  code,
				added: true,
				trace: trace,
			}
		} else {
			e.trace = fmt.Errorf("%w %s: %s", e.trace, key, traceValue)
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

func (e *Error) Trace(operationValue string) *Error {
	key := "Trace"

	e.trace = fmt.Errorf("%w; %s: %s", e.trace, key, operationValue)
	return e
}

func (e *Error) ToClient() error {
	if e != nil {
		if e.Code == 500 {
			return fmt.Errorf("%s", "Internal Server Error")
		}
		return fmt.Errorf("%w", e.Cause)
	}
	return nil
}

func (e *Error) Stack() error {
	if e != nil {
		return fmt.Errorf("%w | %s", e.Cause, e.trace)
	}
	return nil
}
