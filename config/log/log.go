package logger

import (
	"fmt"
	"io"
	"sync"

	"github.com/fatih/color"
)

type colors color.Attribute

type LogColorsConfig struct {
	Info    colors
	Error   colors
	Warn    colors
	Debug   colors
	Fatal   colors
	Message colors
	Trace   colors
}

var instanceColorsConfig *LogColorsConfig
var onceConfig sync.Once

const (
	Black = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

func ConfigDefault(out io.Writer) {
	onceConfig.Do(func() {
		Init(out)
		instanceColorsConfig = &LogColorsConfig{
			Info:    Green,
			Error:   Red,
			Warn:    Yellow,
			Debug:   Cyan,
			Fatal:   Red,
			Message: Magenta,
			Trace:   Blue,
		}
	})
}

func ConfigCustom(info, err, warn, debug, fatal, message, trace colors, out io.Writer) {
	onceConfig.Do(func() {
		Init(out)
		instanceColorsConfig = &LogColorsConfig{
			Info:    info,
			Error:   err,
			Warn:    warn,
			Debug:   debug,
			Fatal:   fatal,
			Message: message,
			Trace:   trace,
		}
	})
}

func Format(col colors, value string) string {
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", col, value)
}
