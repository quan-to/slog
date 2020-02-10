package slog

import (
	"io"
	"os"
	"strings"
)

// TODO: Syslog Output

// FieldRepresentationType specifies which log instance fields formatting should be used
type FieldRepresentationType int

const (
	// NoFields disables the representation of the log instance fields
	NoFields FieldRepresentationType = iota
	// JSONFields enables the representation of log instance fields and formats them as a json string
	JSONFields
	// KeyValueFields enables the representation of log instance fields and formats them as a comma separated key=value fields
	KeyValueFields
)

// Format specifies the logging format (could be pipe separated, JSON, ...)
type Format string

const (
	// JSON specifies to log in JSON format
	JSON Format = "json"
	// PIPE specifies to log in Pipe Delimited Text format
	PIPE Format = "pipe"
)

// ToFormat converts a string to  its corresponding Format type
func ToFormat(s string) Format {
	switch strings.ToLower(s) {
	case string(JSON):
		return JSON
	case string(PIPE):
		return PIPE
	default:
		return PIPE
	}
}

// region Global
var enabledLevels = map[LogLevel]bool{
	DEBUG: true,
	WARN:  true,
	ERROR: true,
	INFO:  true,
	FATAL: true,
}

var fieldRepresentation = JSONFields
var logFormat = PIPE
var defaultOut io.Writer = os.Stdout

var showLines = false

var glog *slogInstance

func init() {
	glog = Scope("Global").(*slogInstance)
	glog.stackOffset += 1 // This will be called from global context, so the stack has one more level
}

// LogNoFormat prints a log string without any ANSI formatting
func LogNoFormat(str interface{}, v ...interface{}) Instance {
	return glog.LogNoFormat(str, v...)
}

// Log is equivalent of calling Info. It logs out a message in INFO level
func Log(str interface{}, v ...interface{}) Instance {
	return glog.Log(str, v...)
}

// Info logs out a message in INFO level
func Info(str interface{}, v ...interface{}) Instance {
	return glog.Info(str, v...)
}

// Debug logs out a message in DEBUG level
func Debug(str interface{}, v ...interface{}) Instance {
	return glog.Debug(str, v...)
}

// Warn logs out a message in WARN level
func Warn(str interface{}, v ...interface{}) Instance {
	return glog.Warn(str, v...)
}

// Error logs out a message in ERROR level
func Error(str interface{}, v ...interface{}) Instance {
	return glog.Error(str, v...)
}

// Fatal logs out a message in ERROR level and closes the program
func Fatal(str interface{}, v ...interface{}) {
	glog.Fatal(str, v)
}

// Scope creates a new slog Instance with the specified root scope
func Scope(scope string) Instance {
	return &slogInstance{
		scope:       []string{scope},
		customOut:   defaultOut,
		stackOffset: 5,
		tag:         "NONE",
		op:          MSG,
	}
}

// SetDefaultOutput sets the Global Default Output I/O and for every new instance created by Scope function
func SetDefaultOutput(o io.Writer) {
	defaultOut = o
	glog.customOut = o
}

// SetDebug globally sets if the DEBUG level messages will be shown. Affects all instances
func SetDebug(enabled bool) {
	enabledLevels[DEBUG] = enabled
}

// SetWarning globally sets if the WARN level messages will be shown. Affects all instances
func SetWarning(enabled bool) {
	enabledLevels[WARN] = enabled
}

// SetInfo globally sets if the INFO level messages will be shown. Affects all instances
func SetInfo(enabled bool) {
	enabledLevels[INFO] = enabled
}

// SetError globally sets if the ERROR level messages will be shown. Affects all instances
func SetError(enabled bool) {
	enabledLevels[ERROR] = enabled
}

// SetShowLines globally sets if the filename and line of the caller function will be shown. Affects all instances
func SetShowLines(enabled bool) {
	showLines = enabled
}

// SetFieldRepresentation globally sets if the representation of log fields. Affects all instances
func SetFieldRepresentation(representationType FieldRepresentationType) {
	fieldRepresentation = representationType
}

// SetLogFormat globally sets the logging format. Affects all instances
func SetLogFormat(f Format) {
	logFormat = f
}

// SetTestMode sets the SLog Instances to test mode a.k.a. all logs disabled. Equivalent to set all levels visibility to false
func SetTestMode() {
	SetDebug(false)
	SetWarning(false)
	SetInfo(false)
	SetError(false)
}

// UnsetTestMode sets the SLog Instances to default mode a.k.a. all logs enabled. Equivalent to set all levels visibility to true
func UnsetTestMode() {
	SetDebug(true)
	SetWarning(true)
	SetInfo(true)
	SetError(true)
}

// DebugEnabled returns if the DEBUG level messages are currently enabled
func DebugEnabled() bool {
	return enabledLevels[DEBUG]
}

// WarningEnabled returns if the WARN level messages are currently enabled
func WarningEnabled() bool {
	return enabledLevels[WARN]
}

// InfoEnabled returns if the INFO level messages are currently enabled
func InfoEnabled() bool {
	return enabledLevels[INFO]
}

// ErrorEnabled returns if the ERROR level messages are currently enabled
func ErrorEnabled() bool {
	return enabledLevels[ERROR]
}

// ShowLinesEnabled returns if the show filename and line from called function is currently enabled
func ShowLinesEnabled() bool {
	return showLines
}

// endregion
