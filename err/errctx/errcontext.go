package errctx

import (
	"fmt"

	"github.com/pkg/errors"
)

type CustomError struct {
	err   error
	added bool
}

var e *CustomError

func ErrCtx(err error, value string) error {
	key := "Operation"

	if e == nil || !e.added {
		e = &CustomError{
			err:   errors.Wrap(err, "Error"),
			added: true,
		}
	} else {
		e = &CustomError{
			err: errors.WithStack(err),
		}
	}
	return fmt.Errorf("%w %s: %s", e.err, key, value)
}
