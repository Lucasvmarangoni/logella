package errs

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

type Error struct {
	Cause   error
	Code    int
	Message string
	trace   error
}

func (e *Error) Error() string {
	if e.Code == http.StatusInternalServerError {
		return e.Message
	}
	return e.Message + " : " + e.Cause.Error()
}

func Wrap(cause error, code int) *Error {
	var e *Error
	if pc, _, _, ok := runtime.Caller(1); ok {
		if fn := runtime.FuncForPC(pc); fn != nil {
			traceValue := strings.Split(fn.Name(), ".")

			trace := fmt.Errorf("trace %s", traceValue[len(traceValue)-1])

			if cause != nil {
				if e, ok = cause.(*Error); !ok {
					e = &Error{
						Cause: cause,
						Code:  code,
						trace: trace,
					}
				}
			}
		}
	}
	return e
}

func Trace(err error) *Error {	
	e := err.(*Error)
	if pc, _, _, ok := runtime.Caller(1); ok {
		if fn := runtime.FuncForPC(pc); fn != nil {
			traceValue := strings.Split(fn.Name(), ".")
			e.trace = fmt.Errorf("%w ➤ %s", e.trace, traceValue[len(traceValue)-1])
		}
	}
	return e
}

func (e *Error) Msg(message string) *Error {
	e.Message = message
	return e
}

func (e *Error) ToClient() error {
	if e != nil {
		if e.Code == http.StatusInternalServerError {
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

func Unwrap(err error) *Error {
	return err.(*Error)
}
