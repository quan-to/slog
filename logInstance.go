package slog

import "io"

// Instance is a interface to a compatible SLog Logging Instance
type Instance interface {
	// Scope returns a new instance with the specified root scope (parent scope is discarded)
	Scope(string) Instance
	// SubScope returns a new instance with the specified scope appended to parent scope
	SubScope(string) Instance
	// WithCustomWriter returns a new instance with the specified custom output
	WithCustomWriter(io.Writer) Instance
	// WithFields returns a new instance with the parent fields plus the current fields. If key collision happens, the value specified in fields argument will be used.
	WithFields(map[string]interface{}) Instance
	// Tag returns a new instance with the specified tag.
	Tag(string) Instance
	// Operation returns a new instance with the specified operation.
	Operation(LogOperation) Instance

	// LogNoFormat prints a log string without any ANSI formatting
	LogNoFormat(interface{}, ...interface{}) Instance

	// Info logs out a message in INFO level
	Info(str interface{}, v ...interface{}) Instance
	// Debug logs out a message in DEBUG level
	Debug(str interface{}, v ...interface{}) Instance
	// Warn logs out a message in WARN level
	Warn(str interface{}, v ...interface{}) Instance
	// Error logs out a message in ERROR level
	Error(str interface{}, v ...interface{}) Instance
	// Fatal logs out a message in ERROR level and closes the program
	Fatal(str interface{}, v ...interface{})

	// Note logs out a message in INFO level and with Operation NOTE. Returns an instance of operation NOTE
	Note(interface{}, ...interface{}) Instance
	// Await logs out a message in INFO level and with Operation AWAIT. Returns an instance of operation AWAIT
	Await(interface{}, ...interface{}) Instance
	// Done logs out a message in INFO level and with Operation DONE. Returns an instance of operation DONE
	Done(interface{}, ...interface{}) Instance
	// Success logs out a message in INFO level and with Operation DONE. Returns an instance of operation DONE
	Success(interface{}, ...interface{}) Instance
	// IO logs out a message in INFO level and with Operation IO. Returns an instance of operation IO
	IO(interface{}, ...interface{}) Instance
	// Log is equivalent of calling Info. It logs out a message in INFO level
	Log(interface{}, ...interface{}) Instance
}
