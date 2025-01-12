package errs

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

func PanicErr(err error, ctx string) {
	if err != nil {
		panic(fmt.Errorf("%w  - ctx: %s", err, ctx))
	}
}

func PanicBool(boolean bool, msg string) {
	if !boolean {
		panic(msg)
	}
}

func FailOnErrLog(err error, ctx, msg string) {
	if err != nil {
		log.Fatal().Stack().Err(Ctx(err, ctx)).Msg(msg)
	}
}
