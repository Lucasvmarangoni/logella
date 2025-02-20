package test_errs

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	// "net/http"
	// "github.com/jackc/pgconn"
	// "github.com/pkg/errors"
)

type Test_Error struct {
	Cause   error
	Code    int
	Message string
	added   bool
	trace   error
}

func (e *Test_Error) Error() string {
	if e.Code == http.StatusInternalServerError {
		return e.Message
	}
	return e.Message + " : " + e.Cause.Error()
}

func Test_Wrap(cause error, code int) error {
	key := "Trace"
	var e *Test_Error
	if pc, _, _, ok := runtime.Caller(1); ok {
		if fn := runtime.FuncForPC(pc); fn != nil {
			traceValue := strings.Split(fn.Name(), ".")

			trace := fmt.Errorf("%s: %s", key, traceValue[len(traceValue)-1])

			if cause != nil {
				var ok bool
				if e, ok = cause.(*Test_Error); !ok {
					e = &Test_Error{
						Cause: cause,
						Code:  code,
						added: true,
						trace: trace,
					}
				}
			}
		}
	}
	return e
}

func Test_Trace(err error) *Test_Error {
	key := "Trace"

	e := err.(*Test_Error)
	if pc, _, _, ok := runtime.Caller(1); ok {
		if fn := runtime.FuncForPC(pc); fn != nil {
			traceValue := strings.Split(fn.Name(), ".")
			e.trace = fmt.Errorf("%w; %s: %s", e.trace, key, traceValue[len(traceValue)-1])
		}
	}
	return e
}

func (e *Test_Error) Test_Stack() error {
	if e != nil {
		return fmt.Errorf("%w | %s", e.Cause, e.trace)
	}
	return nil
}
