package slog

import (
	"encoding/json"
	"fmt"
	"github.com/logrusorgru/aurora"
	"io"
	"path"
	"reflect"
	"runtime"
	"strings"
	"time"
)

type LogLevel string

const (
	LogInfo  LogLevel = "I"
	LogWarn           = "W"
	LogDebug          = "D"
	LogError          = "E"
)

var (
	ColorInfo  = aurora.Cyan
	ColorWarn  = aurora.Brown
	ColorError = aurora.Red
	ColorDebug = aurora.Magenta
)

var logBaseFormat = "%s|%1s| %30v | %s" + LineBreak
var logBaseWithFieldsFormat = "%s|%5s| %30v | %s | %v" + LineBreak

func SetScopeLength(length int) {
	logBaseFormat = "%s|%1s| %" + fmt.Sprintf("%d", length) + "v | %s" + LineBreak
	logBaseWithFieldsFormat = "%s|%1s| %" + fmt.Sprintf("%d", length) + "v | %s | %v" + LineBreak
}

func init() {
	SetScopeLength(30)
}

func buildFieldString(data map[string]interface{}) string {
	retVal := ""
	switch fieldRepresentation {
	case JSONFields:
		v, _ := json.Marshal(data)
		retVal = string(v)
	case KeyValueFields:
		for k, v := range data {
			retVal += fmt.Sprintf("%s=%v,", k, v)
		}
	}

	return retVal
}

func formatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

func hasFormatData(str string) bool {
	// TODO: Better testing
	return strings.Index(str, "%") > -1
}

type Instance struct {
	scope       string
	fields      map[string]interface{}
	customOut   io.Writer
	stackOffset int
}

func getCallerString(stackOffset int) string {
	_, fn, line, _ := runtime.Caller(stackOffset)
	return fmt.Sprintf("%s:%d", path.Base(fn), line)
}

func (i *Instance) commonLog(str string, level LogLevel, c func(arg interface{}) aurora.Value, v ...interface{}) {
	txt := ""

	baseString := c(fmt.Sprintf(asString(str), v...)).String()

	if showLines {
		cs := getCallerString(i.stackOffset)
		baseString = fmt.Sprintf("%25s | %s", aurora.Blue(cs).String(), baseString)
	}

	if i.fields != nil {
		fieldsTxt := buildFieldString(i.fields)
		txt = fmt.Sprintf(logBaseWithFieldsFormat, c(formatTime(time.Now())), c(level), c(aurora.Bold(i.scope)).String(), baseString, c(fieldsTxt))
	} else {
		txt = fmt.Sprintf(logBaseFormat, c(formatTime(time.Now())), c(level), c(aurora.Bold(i.scope)).String(), baseString)
	}

	_, _ = i.Write([]byte(txt))
}

func (i *Instance) argsOnlyLog(str interface{}, level LogLevel, c func(arg interface{}) aurora.Value, v ...interface{}) {
	txt := ""

	args := append([]interface{}{str}, v...)

	baseFormat := ""

	for range args {
		baseFormat += "%v "
	}

	baseString := c(fmt.Sprintf(baseFormat, args...)).String()

	if showLines {
		cs := getCallerString(i.stackOffset)
		baseString = fmt.Sprintf("%25s | %s", aurora.Blue(cs).String(), baseString)
	}

	if i.fields != nil {
		fieldsTxt := buildFieldString(i.fields)
		txt = fmt.Sprintf(logBaseWithFieldsFormat, c(formatTime(time.Now())), c(level), c(aurora.Bold(i.scope)).String(), baseString, c(fieldsTxt))
	} else {
		txt = fmt.Sprintf(logBaseFormat, c(formatTime(time.Now())), c(level), c(aurora.Bold(i.scope)).String(), baseString)
	}

	_, _ = i.Write([]byte(txt))
}

func (i *Instance) log(str interface{}, level LogLevel, c func(arg interface{}) aurora.Value, v ...interface{}) {
	switch ft := str.(type) {
	case string: // Use normal logging
		if hasFormatData(ft) {
			i.commonLog(ft, level, c, v...)
		} else {
			i.argsOnlyLog(str, level, c, v...)
		}
	default: // Args only, to enable slog.Info(a,b,c,d,e)
		i.argsOnlyLog(str, level, c, v...)
	}
}

func (i *Instance) Write(p []byte) (n int, err error) {
	if i.customOut != nil {
		return i.customOut.Write(p)
	}

	fmt.Printf(string(p))
	return len(p), nil
}

func (i *Instance) LogNoFormat(str interface{}, v ...interface{}) *Instance {
	if infoEnabled {
		txt := ""
		if i.fields != nil {
			fieldsTxt := buildFieldString(i.fields)
			txt = fmt.Sprintf(logBaseWithFieldsFormat, ColorInfo(formatTime(time.Now())), ColorInfo(LogInfo), ColorInfo(aurora.Bold(i.scope)).String(), fmt.Sprintf(asString(str), v...), fieldsTxt)
		} else {
			txt = fmt.Sprintf(logBaseFormat, ColorInfo(formatTime(time.Now())), ColorInfo(LogInfo), ColorInfo(aurora.Bold(i.scope)).String(), fmt.Sprintf(asString(str), v...))
		}

		_, _ = i.Write([]byte(txt))
	}
	return i
}

func (i *Instance) Log(str interface{}, v ...interface{}) *Instance {
	// Do not call i.Info, to not change the stack and break filename:line
	if infoEnabled {
		i.log(str, LogInfo, ColorInfo, v...)
	}
	return i
}

func (i *Instance) Info(str interface{}, v ...interface{}) *Instance {
	if infoEnabled {
		i.log(str, LogInfo, ColorInfo, v...)
	}
	return i
}

func (i *Instance) Debug(str interface{}, v ...interface{}) *Instance {
	if debugEnabled {
		i.log(str, LogDebug, ColorDebug, v...)
	}
	return i
}

func (i *Instance) Warn(str interface{}, v ...interface{}) *Instance {
	if warnEnabled {
		i.log(str, LogWarn, ColorWarn, v...)
	}
	return i
}

func (i *Instance) Error(str interface{}, v ...interface{}) *Instance {
	if errorEnabled {
		i.log(str, LogError, ColorError, v...)
	}
	return i
}

func (i *Instance) Fatal(str interface{}, v ...interface{}) {
	varargs := v
	if len(varargs) == 1 && reflect.TypeOf(v[0]) == reflect.TypeOf([]interface{}{}) {
		varargs = v[0].([]interface{})
	}

	var msg string
	if varargs == nil || len(varargs) == 0 {
		msg = asString(str)
	} else {
		msg = fmt.Sprintf(asString(str), varargs...)
	}

	i.log(msg, LogError, ColorError)
	panic(msg)
}

func (i *Instance) WithFields(fields map[string]interface{}) *Instance {
	if i.fields != nil {
		// Append Parent fields
		for k, v := range i.fields {
			if fields[k] == nil {
				fields[k] = v
			}
		}
	}

	i2 := i.clone()
	i2.fields = fields
	return i2
}

func (i *Instance) WithCustomWriter(w io.Writer) *Instance {
	i2 := i.clone()
	i2.customOut = w
	return i2
}

func (i *Instance) SubScope(scope string) *Instance {
	i2 := i.clone()
	i2.scope = fmt.Sprintf("%s ▶ %s", i.scope, scope)
	return i2
}

func (i *Instance) clone() *Instance {
	return &Instance{
		fields:      i.fields,
		scope:       i.scope,
		customOut:   i.customOut,
		stackOffset: i.stackOffset,
	}
}
