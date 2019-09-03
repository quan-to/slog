package slog

import (
	"io"
	"os"
)

// TODO: Syslog Output

type FieldRepresentationType int

const (
	NoFields FieldRepresentationType = iota
	JSONFields
	KeyValueFields
)

// region Global
var enabledLevels = map[LogLevel]bool{
	DEBUG: true,
	WARN:  true,
	ERROR: true,
	INFO:  true,
}

var fieldRepresentation = JSONFields
var defaultOut io.Writer = os.Stdout

var showLines = false

var glog *slogInstance

func init() {
	glog = Scope("Global").(*slogInstance)
	glog.stackOffset += 1 // This will be called from global context, so the stack has one more level
}

func LogNoFormat(str interface{}, v ...interface{}) Instance {
	return glog.LogNoFormat(str, v...)
}

func Log(str interface{}, v ...interface{}) Instance {
	return glog.Log(str, v...)
}

func Info(str interface{}, v ...interface{}) Instance {
	return glog.Info(str, v...)
}

func Debug(str interface{}, v ...interface{}) Instance {
	return glog.Debug(str, v...)
}

func Warn(str interface{}, v ...interface{}) Instance {
	return glog.Warn(str, v...)
}

func Error(str interface{}, v ...interface{}) Instance {
	return glog.Error(str, v...)
}

func Fatal(str interface{}, v ...interface{}) {
	glog.Fatal(str, v)
}

func Scope(scope string) Instance {
	return &slogInstance{
		scope:       scope,
		customOut:   defaultOut,
		stackOffset: 4,
	}
}

func SetDefaultOutput(o io.Writer) {
	defaultOut = o
	glog.customOut = o
}

func SetDebug(enabled bool) {
	enabledLevels[DEBUG] = enabled
}
func SetWarning(enabled bool) {
	enabledLevels[WARN] = enabled
}
func SetInfo(enabled bool) {
	enabledLevels[INFO] = enabled
}
func SetError(enabled bool) {
	enabledLevels[ERROR] = enabled
}
func SetShowLines(enabled bool) {
	showLines = enabled
}

func SetFieldRepresentation(representationType FieldRepresentationType) {
	fieldRepresentation = representationType
}

func SetTestMode() {
	SetDebug(false)
	SetWarning(false)
	SetInfo(false)
	SetError(false)
}

func UnsetTestMode() {
	SetDebug(true)
	SetWarning(true)
	SetInfo(true)
	SetError(true)
}

func DebugEnabled() bool {
	return enabledLevels[DEBUG]
}
func WarningEnabled() bool {
	return enabledLevels[WARN]
}
func InfoEnabled() bool {
	return enabledLevels[INFO]
}
func ErrorEnabled() bool {
	return enabledLevels[ERROR]
}
func ShowLinesEnabled() bool {
	return showLines
}

// endregion
