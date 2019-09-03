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

func (i *slogInstance) buildText(str string, level LogLevel, v ...interface{}) string {
	logDate := aurora.Gray(formatTime(time.Now()))
	levelColor := levelColors[level]
	scope := padRight(strings.Join(i.scope, " > "), scopeLength)
	stringifiedFields := "{}"

	if i.fields != nil {
		stringifiedFields = buildFieldString(i.fields)
	}

	op := operationColors[i.op](padRight(string(i.op), maxOperationStringLength))
	tag := aurora.Gray(i.tag)

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

func (i *slogInstance) Write(p []byte) (n int, err error) {
	if i.customOut != nil {
		return i.customOut.Write(p)
	}

	fmt.Printf(string(p))
	return len(p), nil
}

func (i *slogInstance) LogNoFormat(str interface{}, v ...interface{}) Instance {
	if enabledLevels[INFO] {
		txt := stripColors(i.buildText(asString(str), INFO, v...))
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
	i2.scope = []string{scope}
	return i2
}

func (i *slogInstance) SubScope(scope string) Instance {
	i2 := i.clone()
	i2.scope = append(i2.scope, scope)
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
		op:          i.op,
		tag:         i.tag,
	}
}
