package errs

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type Error struct {
	Err        error
	Code       int
	Message    string
	added      bool
	operations error
}

var Status = map[int]string{
	http.StatusBadRequest:          "BadRequest",
	http.StatusUnauthorized:        "Unauthorized",
	http.StatusForbidden:           "Forbidden",
	http.StatusInternalServerError: "InternalServerError",
}

var E *Error

func Wrap(err error, operationValue string, code int) error {
	key := "Operation"
	operation := fmt.Errorf("%s: %s", key, operationValue)

	if E == nil {
		E = &Error{
			Err:        fmt.Errorf("%w", err),
			Code:       code,
			added:      true,
			operations: operation,
		}
	} else {
		E.operations = fmt.Errorf("%w %s: %s", E.operations, key, operationValue)
		E.Err = errors.WithStack(err)
	}
	return E.Err
}

func Msg(message string) {
	E.Message = message
}

func GetOperations() error {
	if E != nil {
		return E.operations
	}
	return nil
}

func Stack() error {
	if E != nil {
		return fmt.Errorf("%w | %s", E.Err, E.operations)
	}
	return nil
}
