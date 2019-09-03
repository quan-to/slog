package slog

import "github.com/logrusorgru/aurora"

type LogLevel string

const (
	INFO     LogLevel = "I"
	WARN              = "W"
	ERROR             = "E"
	DEBUG             = "D"
	LogInfo           = "I"
	LogWarn           = "W"
	LogDebug          = "D"
	LogError          = "E"
)

var levelColors = map[LogLevel]colorFunc{
	INFO:  aurora.Cyan,
	WARN:  aurora.Brown,
	ERROR: aurora.Red,
	DEBUG: aurora.Magenta,
}
