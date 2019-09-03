package slog

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

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
