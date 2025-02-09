package logger

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"os"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false, TimeFormat: time.RFC3339,
		FormatMessage: func(i interface{}) string {
			value := failOnError(i)
			return Format(instanceColorsConfig.Message, strings.ToUpper(value))
		},
		FormatLevel: func(i interface{}) string {
			level := failOnError(i)
			switch level {
			case "info":
				return Format(instanceColorsConfig.Info, strings.ToUpper(level)+" ⇝")
			case "error":
				return Format(instanceColorsConfig.Error, strings.ToUpper(level)+" ⇝")
			case "warn":
				return Format(instanceColorsConfig.Warn, strings.ToUpper(level)+" ⇝")
			case "debug":
				return Format(instanceColorsConfig.Debug, strings.ToUpper(level)+" ⇝")
			case "fatal":
				return Format(instanceColorsConfig.Fatal, strings.ToUpper(level)+" ⇝")
			default:
				return level
			}
		},
		FormatErrFieldValue: func(i interface{}) string {
			value := failOnError(i)
			formattedTrace := Format(instanceColorsConfig.Trace, "trace")
			formattedError := Format(Red, "Error")
			Str := strings.ReplaceAll(value, "trace", formattedTrace)
			Str = strings.ReplaceAll(Str, "Error", formattedError)
			return Str
		},
		FormatErrFieldName: func(i interface{}) string {
			value := failOnError(i)
			return value
		},
	})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return fmt.Sprintf("[\n%s:%d\n]", file, line)
	}
}

func failOnError(i interface{}) string {
	value, ok := i.(string)
	if !ok {
		return "unknown"
	}
	return value
}
