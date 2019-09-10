package slog

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"io"
	"reflect"
	"strings"
	"time"
)

var scopeLength = 24

// SetScopeLength sets the scope field length (adds left pad when nescessary) - Affects globally all SLog Instances
func SetScopeLength(length int) {
	scopeLength = length
}

type slogInstance struct {
	scope       []string
	fields      map[string]interface{}
	customOut   io.Writer
	stackOffset int
	tag         string
	op          LogOperation
}

func (i *slogInstance) incStackOffset() *slogInstance {
	i.stackOffset++
	return i
}

func (i *slogInstance) buildText(str string, level LogLevel, v ...interface{}) string {
	logDate := aurora.Gray(7, formatTime(time.Now()))
	levelColor := levelColors[level]
	scope := padRight(strings.Join(i.scope, " > "), scopeLength)
	stringifiedFields := "{}"

	if i.fields != nil {
		stringifiedFields = buildFieldString(i.fields)
	}

	op := operationColors[i.op](padRight(string(i.op), maxOperationStringLength)).White()
	tag := aurora.Gray(7, i.tag)

	logHead := logDate.String() + " " + pipeChar + " " + levelColor(aurora.Bold(level)).String() + " " + pipeChar + " " + op.String() + " " + pipeChar + " " + tag.String() + " " + pipeChar + " " + scope + " " + pipeChar + " "

	if showLines {
		cs := getCallerString(i.stackOffset)
		logHead += cs + " " + pipeChar + " "
	}

	logTail := pipeChar + " " + stringifiedFields
	logHeadLength := len(stripColors(logHead))

	baseString := fmt.Sprintf(asString(str), v...)
	baseString = addPadForLines(baseString, logHeadLength)

	return logHead + levelColor(baseString).String() + " " + logTail + LineBreak
}

func (i *slogInstance) commonLog(str string, level LogLevel, v ...interface{}) {
	_, _ = i.Write([]byte(i.buildText(str, level, v...)))
}

func (i *slogInstance) argsOnlyLog(str interface{}, level LogLevel, v ...interface{}) {
	args := append([]interface{}{str}, v...)

	baseFormat := ""

	for range args {
		baseFormat += "%v "
	}

	_, _ = i.Write([]byte(i.buildText(baseFormat, level, args...)))
}

func (i *slogInstance) log(str interface{}, level LogLevel, v ...interface{}) {
	switch ft := str.(type) {
	case string: // Use normal logging
		if hasFormatData(ft) {
			i.commonLog(ft, level, v...)
		} else {
			i.argsOnlyLog(str, level, v...)
		}
	default: // Args only, to enable slog.Info(a,b,c,d,e)
		i.argsOnlyLog(str, level, v...)
	}
}

// Write writes bytes to the instance output
func (i *slogInstance) Write(p []byte) (n int, err error) {
	if i.customOut != nil {
		return i.customOut.Write(p)
	}

	fmt.Printf(string(p))
	return len(p), nil
}

// LogNoFormat prints a log string without any ANSI formatting
func (i *slogInstance) LogNoFormat(str interface{}, v ...interface{}) Instance {
	if enabledLevels[INFO] {
		txt := stripColors(i.buildText(asString(str), INFO, v...))
		_, _ = i.Write([]byte(txt))
	}
	return i
}

// Log is equivalent of calling Info. It logs out a message in INFO level
func (i *slogInstance) Log(str interface{}, v ...interface{}) Instance {
	// Do not call i.Info, to not change the stack and break filename:line
	if enabledLevels[INFO] {
		i.log(str, INFO, v...)
	}
	return i
}

// Info logs out a message in INFO level
func (i *slogInstance) Info(str interface{}, v ...interface{}) Instance {
	if enabledLevels[INFO] {
		i.log(str, INFO, v...)
	}
	return i
}

// Debug logs out a message in DEBUG level
func (i *slogInstance) Debug(str interface{}, v ...interface{}) Instance {
	if enabledLevels[DEBUG] {
		i.log(str, DEBUG, v...)
	}
	return i
}

// Warn logs out a message in WARN level
func (i *slogInstance) Warn(str interface{}, v ...interface{}) Instance {
	if enabledLevels[WARN] {
		i.log(str, WARN, v...)
	}
	return i
}

// Error logs out a message in ERROR level
func (i *slogInstance) Error(str interface{}, v ...interface{}) Instance {
	if enabledLevels[ERROR] {
		i.log(str, ERROR, v...)
	}
	return i
}

// Fatal logs out a message in ERROR level and closes the program
func (i *slogInstance) Fatal(str interface{}, v ...interface{}) {
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

	i.log(msg, ERROR)
	panic(msg)
}

// WithFields returns a new instance with the parent fields plus the current fields. If key collision happens, the value specified in fields argument will be used.
func (i *slogInstance) WithFields(fields map[string]interface{}) Instance {
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

// WithCustomWriter returns a new instance with the specified custom output
func (i *slogInstance) WithCustomWriter(w io.Writer) Instance {
	i2 := i.clone()
	i2.customOut = w
	return i2
}

// Scope returns a new instance with the specified root scope (parent scope is discarded)
func (i *slogInstance) Scope(scope string) Instance {
	i2 := i.clone()
	i2.scope = []string{scope}
	return i2
}

// SubScope returns a new instance with the specified scope appended to parent scope
func (i *slogInstance) SubScope(scope string) Instance {
	i2 := i.clone()
	i2.scope = append(i2.scope, scope)
	return i2
}

// Tag returns a new instance with the specified tag.
func (i *slogInstance) Tag(tag string) Instance {
	i2 := i.clone()
	i2.tag = tag
	return i2
}

// Operation returns a new instance with the specified operation.
func (i *slogInstance) Operation(op LogOperation) Instance {
	i2 := i.clone()
	i2.op = op
	return i2
}

func (i *slogInstance) clone() *slogInstance {
	return &slogInstance{
		fields:      i.fields,
		scope:       i.scope,
		customOut:   i.customOut,
		stackOffset: i.stackOffset,
		op:          i.op,
		tag:         i.tag,
	}
}
