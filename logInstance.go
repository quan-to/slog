package slog

import "io"

// Instance is a interface to a compatible SLog Logging Instance
type Instance interface {
	// Scope creates a new instance of slog with the specified scope
	Scope(string) Instance
	// SubScope creates a new instance of slog with the specified scope appened to the parent scope
	SubScope(string) Instance
	// WithCustomWriter creates a new instance of slog with the specified writer as default output
	WithCustomWriter(io.Writer) Instance
	// WithFields creates a new instance of slog with the specified fields plus the parent fields
	WithFields(map[string]interface{}) Instance
	// Tag creates a new instance of slog with the specified tag
	Tag(string) Instance
	// Operation creates a new instance of slog with the specified operation
	Operation(LogOperation) Instance

	LogNoFormat(interface{}, ...interface{}) Instance

	// Info prints a INFO log entry (same as Log(INFO, ...)
	Info(str interface{}, v ...interface{}) Instance
	// Debug prints a DEBUG log entry (same as Log(DEBUG, ...)
	Debug(str interface{}, v ...interface{}) Instance
	// Warn prints a WARN log entry (same as Log(WARN, ...)
	Warn(str interface{}, v ...interface{}) Instance
	// Error prints a ERROR log entry (same as Log(ERROR, ...)
	Error(str interface{}, v ...interface{}) Instance
	// Fatal prints a FATAL log entry (same as Log(ERROR, ...) and crashes the program
	Fatal(str interface{}, v ...interface{})

	// Note prints a INFO log entry with operation NOTE (same as Operation(NOTE).Log(INFO, ...)
	Note(interface{}, ...interface{}) Instance
	// Await prints a INFO log entry with operation AWAIT (same as Operation(AWAIT).Log(INFO, ...)
	Await(interface{}, ...interface{}) Instance
	// Done prints a INFO log entry with operation DONE (same as Operation(DONE).Log(INFO, ...)
	Done(interface{}, ...interface{}) Instance
	// Success prints a INFO log entry with operation DONE (same as Done(...))
	Success(interface{}, ...interface{}) Instance
	// IO prints a INFO log entry with operation IO (same as Operation(IO).Log(INFO, ...)
	IO(interface{}, ...interface{}) Instance
	// Log same as Info
	Log(interface{}, ...interface{}) Instance
}
