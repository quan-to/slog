package slog

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"io"
	"path"
	"reflect"
	"runtime"
	"time"
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

type slogInstance struct {
	scope       string
	fields      map[string]interface{}
	customOut   io.Writer
	stackOffset int
	tag         string
	op          LogOperation
}

func getCallerString(stackOffset int) string {
	_, fn, line, _ := runtime.Caller(stackOffset)
	return fmt.Sprintf("%s:%d", path.Base(fn), line)
}

func (i *slogInstance) commonLog(str string, level LogLevel, v ...interface{}) {
	txt := ""
	c := levelColors[level]
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

func (i *slogInstance) argsOnlyLog(str interface{}, level LogLevel, v ...interface{}) {
	txt := ""

	c := levelColors[level]

	args := append([]interface{}{str}, v...)

	baseFormat := ""

	for range args {
		baseFormat += "%v "
	}
	fmt.Printf("%v - %s\n", c, level)
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

func (i *slogInstance) Write(p []byte) (n int, err error) {
	if i.customOut != nil {
		return i.customOut.Write(p)
	}

	fmt.Printf(string(p))
	return len(p), nil
}

func (i *slogInstance) LogNoFormat(str interface{}, v ...interface{}) Instance {
	if enabledLevels[INFO] {
		ColorInfo := levelColors[INFO]
		txt := ""
		if i.fields != nil {
			fieldsTxt := buildFieldString(i.fields)
			txt = fmt.Sprintf(logBaseWithFieldsFormat, ColorInfo(formatTime(time.Now())), ColorInfo(INFO), ColorInfo(aurora.Bold(i.scope)).String(), fmt.Sprintf(asString(str), v...), fieldsTxt)
		} else {
			txt = fmt.Sprintf(logBaseFormat, ColorInfo(formatTime(time.Now())), ColorInfo(INFO), ColorInfo(aurora.Bold(i.scope)).String(), fmt.Sprintf(asString(str), v...))
		}

		_, _ = i.Write([]byte(txt))
	}
	return i
}

func (i *slogInstance) Log(str interface{}, v ...interface{}) Instance {
	// Do not call i.Info, to not change the stack and break filename:line
	if enabledLevels[INFO] {
		i.log(str, INFO, v...)
	}
	return i
}

func (i *slogInstance) Info(str interface{}, v ...interface{}) Instance {
	if enabledLevels[INFO] {
		i.log(str, INFO, v...)
	}
	return i
}

func (i *slogInstance) Debug(str interface{}, v ...interface{}) Instance {
	if enabledLevels[DEBUG] {
		i.log(str, DEBUG, v...)
	}
	return i
}

func (i *slogInstance) Warn(str interface{}, v ...interface{}) Instance {
	if enabledLevels[WARN] {
		i.log(str, WARN, v...)
	}
	return i
}

func (i *slogInstance) Error(str interface{}, v ...interface{}) Instance {
	if enabledLevels[ERROR] {
		i.log(str, ERROR, v...)
	}
	return i
}

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

func (i *slogInstance) Note(str interface{}, v ...interface{}) Instance {
	return i.Operation(NOTE).Info(str, v...)
}

func (i *slogInstance) Await(str interface{}, v ...interface{}) Instance {
	return i.Operation(AWAIT).Info(str, v...)
}

func (i *slogInstance) Done(str interface{}, v ...interface{}) Instance {
	return i.Operation(DONE).Info(str, v...)
}

func (i *slogInstance) Success(str interface{}, v ...interface{}) Instance {
	return i.Done(str, v...)
}

func (i *slogInstance) IO(str interface{}, v ...interface{}) Instance {
	return i.Operation(IO).Info(str, v...)
}

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

func (i *slogInstance) WithCustomWriter(w io.Writer) Instance {
	i2 := i.clone()
	i2.customOut = w
	return i2
}

func (i *slogInstance) Scope(scope string) Instance {
	i2 := i.clone()
	i2.scope = scope
	return i2
}

func (i *slogInstance) SubScope(scope string) Instance {
	i2 := i.clone()
	i2.scope = fmt.Sprintf("%s > %s", i.scope, scope)
	return i2
}

func (i *slogInstance) Tag(tag string) Instance {
	i2 := i.clone()
	i2.tag = tag
	return i2
}

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
	}
}
