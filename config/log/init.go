package logger

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

)

func Init(out io.Writer) {
	if out == nil {
		panic("output cannot be nil")
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: out, NoColor: false, TimeFormat: time.RFC3339,
		FormatMessage: func(i any) string {
			value := failOnError(i)
			return Format(instanceColorsConfig.Message, strings.ToUpper(value))
		},
		FormatLevel: func(i any) string {
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
		FormatErrFieldValue: func(i any) string {
			value := failOnError(i)

			formattedFile := Format(instanceColorsConfig.Trace, "path")
			Str := strings.ReplaceAll(value, "path", formattedFile)

			formattedLine := Format(instanceColorsConfig.Trace, ":")
			Str = strings.ReplaceAll(Str, ":", formattedLine)

			formattedTrace := Format(instanceColorsConfig.Trace, "trace")
			Str = strings.ReplaceAll(Str, "trace", formattedTrace)

			formattedTraceArrow := Format(instanceColorsConfig.Trace, "➤")
			Str = strings.ReplaceAll(Str, "➤", formattedTraceArrow)

			return Str
		},
		FormatErrFieldName: func(i any) string {
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

func failOnError(i any) string {
	value, ok := i.(string)
	if !ok {
		return "unknown"
	}
	return value
}
