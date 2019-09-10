package slog

import "github.com/logrusorgru/aurora"

// LogLevel type specifies the level of log to be used
type LogLevel string

const (
	// INFO represents a Information Log Level (or verbose)
	INFO LogLevel = "I"

	// WARN represents a Warning Log Level
	WARN = "W"

	// ERROR represents an error message
	ERROR = "E"

	// FATAL represents an fatal message
	FATAL = "F"

	// DEBUG represents an debug message
	DEBUG = "D"
)

var levelColors = map[LogLevel]colorFunc{
	INFO:  aurora.Cyan,
	WARN:  aurora.Yellow,
	ERROR: aurora.Red,
	DEBUG: aurora.Magenta,
	FATAL: aurora.Red,
}
