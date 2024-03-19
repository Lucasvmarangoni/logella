package logger

import (
	"fmt"
	"sync"

	"github.com/fatih/color"
)

type colors color.Attribute

type LogColorsConfig struct {
	Info      colors
	Error     colors
	Warn      colors
	Debug     colors
	Fatal     colors
	Message   colors
	Operation colors
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

func GetLogColorsConfig() {
	onceConfig.Do(func() {
		instanceColorsConfig = &LogColorsConfig{
			Info:      Green,
			Error:     Red,
			Warn:      Yellow,
			Debug:     Cyan,
			Fatal:     Red,
			Message:   Magenta,
			Operation: Blue,
		}
	})

}

func ConfigDefault() {
	GetLogColorsConfig()
}

func ConfigCustom(info, err, warn, debug, fatal, message, operation colors) {
	onceConfig.Do(func() {
		instanceColorsConfig = &LogColorsConfig{
			Info:      info,
			Error:     err,
			Warn:      warn,
			Debug:     debug,
			Fatal:     fatal,
			Message:   message,
			Operation: operation,
		}
	})
}

func Format(col colors, value string) string {
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", col, value)
}
