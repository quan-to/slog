package slog

import "github.com/logrusorgru/aurora"

type LogLevel string

const (
	INFO  LogLevel = "I"
	WARN           = "W"
	ERROR          = "E"
	DEBUG          = "D"
)

var levelColors = map[LogLevel]colorFunc{
	INFO:  aurora.Cyan,
	WARN:  aurora.Yellow,
	ERROR: aurora.Red,
	DEBUG: aurora.Magenta,
}
