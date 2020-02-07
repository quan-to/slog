package slog

import (
	"encoding/json"
	"fmt"
	"github.com/logrusorgru/aurora"
	"path"
	"regexp"
	"runtime"
	"strings"
	"time"
)

var pipeChar = aurora.Bold("|").White().String()

func buildFieldString(data map[string]interface{}) string {
	retVal := ""
	switch fieldRepresentation {
	case JSONFields:
		retVal = buildJson(data)
	case KeyValueFields:
		for k, v := range data {
			retVal += fmt.Sprintf("%s=%v,", k, v)
		}
	}

	return retVal
}

func buildJson(data map[string]interface{}) string {
	v, _ := json.Marshal(data)
	return string(v)
}

// formatTime returns the specified Date in ISO format (RFC3339)
func formatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

func hasFormatData(str string) bool {
	// TODO: Better testing
	return strings.Index(str, "%") > -1
}

func getCallerString(stackOffset int) string {
	_, fn, line, _ := runtime.Caller(stackOffset)
	return fmt.Sprintf("%s:%d", path.Base(fn), line)
}

func padRight(str string, length int) string {
	for i := len(str); i < length; i++ {
		str = str + " "
	}

	return str
}

func addPadding(str string, length int) string {
	pad := ""
	for i := 0; i < length; i++ {
		pad += " "
	}

	return pad + str
}

func addPadForLines(str string, length int) string {
	lines := strings.Split(str, LineBreak)
	for i, v := range lines {
		if i > 0 {
			lines[i] = addPadding(v, length)
		}
	}

	return strings.Join(lines, LineBreak)
}

var stripColorRgx = regexp.MustCompile("[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))")

func stripColors(str string) string {
	return stripColorRgx.ReplaceAllString(str, "")
}

func stringSliceIndexOf(str string, s []string) int {
	for i, v := range s {
		if v == str {
			return i
		}
	}

	return -1
}
