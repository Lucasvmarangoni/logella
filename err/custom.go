package errs

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/rs/zerolog/log"
)

func New(cause error) *Error {
	var e *Error
	if pc, file, line, ok := runtime.Caller(1); ok {
		projectRoot, _ := os.Getwd()
		relativePath, _ := filepath.Rel(projectRoot, file)
		if fn := runtime.FuncForPC(pc); fn != nil {
			traceValue := strings.Split(fn.Name(), ".")

			trace := fmt.Errorf("path %s:%d trace %s", relativePath, line, traceValue[len(traceValue)-1])

			if cause != nil {
				if e, ok = cause.(*Error); !ok {
					e = &Error{
						Cause: cause,
						trace: trace,
					}
				}
			}
		}
	}
	return e
}

func PanicErr(err error, msg string) {
	if err != nil {
		panic(err)
	}
}

func PanicBool(boolean bool, msg string) {
	if !boolean {
		panic(msg)
	}
}

func FailOnErrLog(err error, msg string) {
	if pc, _, _, ok := runtime.Caller(1); ok {
		if fn := runtime.FuncForPC(pc); fn != nil {
			traceValue := strings.Split(fn.Name(), ".")
			log.Fatal().Stack().Err(
				fmt.Errorf("%w ➤ %s", err.(*Error).trace, traceValue[len(traceValue)-1])).
				Msg(msg)
		}
	}
}

func IsRequiredError(fieldName, msg string) error {
	return fmt.Errorf("%s is required. %s", fieldName, msg)
}

func IsInvalidError(fieldName, msg string) error {
	return fmt.Errorf("%s is invalid. %s", fieldName, msg)
}
