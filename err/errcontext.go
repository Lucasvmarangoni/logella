package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

type CustomError struct {
	err        error
	added      bool
	operations error
}

var e *CustomError

func Ctx(err error, operationValue string) error {
	key := "Operation"
	operation := fmt.Errorf("%s: %s", key, operationValue)

	if e == nil {
		e = &CustomError{
			err:        err,
			added:      true,
			operations: operation,
		}
	} else {
		e.operations = fmt.Errorf("%w %s: %s", e.operations, key, operationValue)
		e.err = errors.WithStack(err)
	}
	return fmt.Errorf("%w", e.err)
}

func GetOperations() error {
	if e != nil {
		return e.operations
	}
	return nil
}

func Stack() error {
	if e != nil {
		return fmt.Errorf("%w | %s", e.err, e.operations)
	}
	return nil
}
