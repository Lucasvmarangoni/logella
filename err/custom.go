package errors

import (
	"fmt"
	"github.com/rs/zerolog/log"
)

func IsRequiredError(fieldName, msg string) error {
	return fmt.Errorf("%s is required. %s", fieldName, msg)
}

func IsInvalidError(fieldName, msg string) error {
	return fmt.Errorf("%s is invalid. %s", fieldName, msg)
}

func FailOnErrPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func FailOnErrLog(err error, msg string) {
	if err != nil {
		log.Fatal().Err(ErrCtx(err, msg))
	}
}
