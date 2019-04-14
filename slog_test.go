package slog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestSetDefaultOutput(t *testing.T) {
	od := defaultOut

	SetShowLines(true)
	SetDefaultOutput(os.Stdout)
	// Crash Tests
	Debug("Test %s %d %f %v", "huebr", 1, 10.0, true)
	Warn("Test %s %d %f %v", "huebr", 1, 10.0, true)
	Info("Test %s %d %f %v", "huebr", 1, 10.0, true)
	Log("Test %s %d %f %v", "huebr", 1, 10.0, true)
	LogNoFormat("Test %s %d %f %v", "huebr", 1, 10.0, true)

	SetDefaultOutput(os.Stderr)
	// Crash Tests
	Debug("Test %s %d %f %v", "huebr", 1, 10.0, true)
	Warn("Test %s %d %f %v", "huebr", 1, 10.0, true)
	Info("Test %s %d %f %v", "huebr", 1, 10.0, true)
	Log("Test %s %d %f %v", "huebr", 1, 10.0, true)
	LogNoFormat("Test %s %d %f %v", "huebr", 1, 10.0, true)

	buff := bytes.NewBufferString("")
	SetDefaultOutput(buff)
	// Crash Tests
	Debug("Test %s %d %f %v", "huebr", 1, 10.0, true)
	Warn("Test %s %d %f %v", "huebr", 1, 10.0, true)
	Info("Test %s %d %f %v", "huebr", 1, 10.0, true)
	Log("Test %s %d %f %v", "huebr", 1, 10.0, true)
	LogNoFormat("Test %s %d %f %v", "huebr", 1, 10.0, true)

	buff.Reset()

	// Test output
	Info("Test %s %d %f %v", "huebr", 1, 10.0, true)

	if strings.Index(buff.String(), "Test huebr 1 10.000000 true") == -1 {
		t.Errorf("Expected string %s in %s", "Test huebr 1 10.000000 true", buff.String())
	}

	SetDefaultOutput(od)
}

func TestArgsOnly(t *testing.T) {
	SetFieldRepresentation(NoFields)
	buff := bytes.NewBufferString("")
	i := Scope("ArgsOnly").WithCustomWriter(buff)

	i.Info("huebr", 1, 10.0, true)

	o := buff.String()
	if strings.Index(o, "huebr 1 10 true") == -1 {
		t.Errorf("Expected \"%s\" in output: \"%s\"", "huebr 1 10.0 true", o)
	}

	SetFieldRepresentation(JSONFields)
	buff = bytes.NewBufferString("")
	i = Scope("ArgsOnly").WithFields(map[string]interface{}{
		"a": "b",
	}).WithCustomWriter(buff)

	i.Info(555, 1, 10.0, true)

	o = buff.String()
	if strings.Index(o, "555 1 10 true") == -1 {
		t.Errorf("Expected \"%s\" in output: \"%s\"", "huebr 1 10.0 true", o)
	}

	SetFieldRepresentation(KeyValueFields)
	buff = bytes.NewBufferString("")
	i = Scope("ArgsOnly").WithFields(map[string]interface{}{
		"a": "b",
	}).WithCustomWriter(buff)

	i.Info("huebr", 1, 10.0, true)

	o = buff.String()
	if strings.Index(o, "huebr 1 10 true") == -1 {
		t.Errorf("Expected \"%s\" in output: \"%s\"", "huebr 1 10.0 true", o)
	}
}

func TestDefaultOutput(t *testing.T) {
	i := Scope("DefaultOutput").WithCustomWriter(nil)

	// Crash Tests
	i.Debug("Test %s %d %f %v", "huebr", 1, 10.0, true)
	i.Warn("Test %s %d %f %v", "huebr", 1, 10.0, true)
	i.Info("Test %s %d %f %v", "huebr", 1, 10.0, true)
	i.Log("Test %s %d %f %v", "huebr", 1, 10.0, true)
	i.LogNoFormat("Test %s %d %f %v", "huebr", 1, 10.0, true)
}

func TestSetTestMode(t *testing.T) {
	SetTestMode()
	// Test Mode should be all unset

	if DebugEnabled() {
		t.Fatalf("Debug is set to true! Should be false")
	}

	if WarningEnabled() {
		t.Fatalf("Warn is set to true! Should be false")
	}

	if WarningEnabled() {
		t.Fatalf("Error is set to true! Should be false")
	}

	if InfoEnabled() {
		t.Fatalf("Info is set to true! Should be false")
	}

	UnsetTestMode()
	// Test Mode should be all set

	if !DebugEnabled() {
		t.Fatalf("Debug is set to false! Should be true")
	}

	if !WarningEnabled() {
		t.Fatalf("Warn is set to false! Should be true")
	}

	if !ErrorEnabled() {
		t.Fatalf("Error is set to false! Should be true")
	}

	if !InfoEnabled() {
		t.Fatalf("Info is set to false! Should be true")
	}
}

func TestSubScope(t *testing.T) {
	i := Scope("ABCD").SubScope("EFGH")

	if strings.Index(i.scope, "ABCD") == -1 {
		t.Errorf("Expected ABCD in Scope")
	}
	if strings.Index(i.scope, "EFGH") == -1 {
		t.Errorf("Expected EFGH in Scope")
	}

	// No Crash tests
	i.Debug("Test %s %d %f %v", "huebr", 1, 10.0, true)
	i.Warn("Test %s %d %f %v", "huebr", 1, 10.0, true)
	i.Info("Test %s %d %f %v", "huebr", 1, 10.0, true)
	i.Log("Test %s %d %f %v", "huebr", 1, 10.0, true)
	i.LogNoFormat("Test %s %d %f %v", "huebr", 1, 10.0, true)
}

func TestWithFields(t *testing.T) {
	SetFieldRepresentation(KeyValueFields)
	i := Scope("WithFields").WithFields(map[string]interface{}{
		"a": "b",
		"b": 5,
	})

	if i.fields["a"] != "b" {
		t.Errorf("Expected field \"a\" to be \"b\"")
	}

	if i.fields["b"] != 5 {
		t.Errorf("Expected field \"b\" to be 5")
	}

	// Child should inherit parent fields and replace the existent ones
	i = i.WithFields(map[string]interface{}{
		"c": 3.14,
		"a": 9,
	})

	if i.fields["a"] != 9 {
		t.Errorf("Expected field \"a\" to be 9")
	}

	if i.fields["b"] != 5 {
		t.Errorf("Expected field \"b\" to be 5")
	}

	if i.fields["c"] != 3.14 {
		t.Errorf("Expected field \"b\" to be 5")
	}

	// No Crash tests
	i.Debug("Test %s %d %f %v", "huebr", 1, 10.0, true)
	i.Warn("Test %s %d %f %v", "huebr", 1, 10.0, true)
	i.Info("Test %s %d %f %v", "huebr", 1, 10.0, true)
	i.Log("Test %s %d %f %v", "huebr", 1, 10.0, true)
	i.LogNoFormat("Test %s %d %f %v", "huebr", 1, 10.0, true)
}

func TestWithFieldsJSON(t *testing.T) {
	SetFieldRepresentation(JSONFields)
	buff := bytes.NewBufferString("")
	i := Scope("WithFieldsJSON").WithFields(map[string]interface{}{
		"a": "b",
		"b": 5,
	}).WithCustomWriter(buff)

	jsonDataB, _ := json.Marshal(i.fields)
	jsonData := string(jsonDataB)

	i.Info("Test %s %d %f %v", "huebr", 1, 10.0, true)

	o := buff.String()
	if strings.Index(o, jsonData) == -1 {
		t.Errorf("Expected \"%s\" in output: \"%s\"", jsonData, o)
	}

	buff.Reset()

	i.Warn("Test %s %d %f %v", "huebr", 1, 10.0, true)

	o = buff.String()
	if strings.Index(o, jsonData) == -1 {
		t.Errorf("Expected \"%s\" in output: \"%s\"", jsonData, o)
	}

	buff.Reset()

	i.Debug("Test %s %d %f %v", "huebr", 1, 10.0, true)

	o = buff.String()
	if strings.Index(o, jsonData) == -1 {
		t.Errorf("Expected \"%s\" in output: \"%s\"", jsonData, o)
	}

	buff.Reset()

	i.Error("Test %s %d %f %v", "huebr", 1, 10.0, true)

	o = buff.String()
	if strings.Index(o, jsonData) == -1 {
		t.Errorf("Expected \"%s\" in output: \"%s\"", jsonData, o)
	}

	buff.Reset()
}

func TestWithFieldsKV(t *testing.T) {
	SetFieldRepresentation(KeyValueFields)
	buff := bytes.NewBufferString("")
	i := Scope("WithFieldsKV").WithFields(map[string]interface{}{
		"a": "b",
		"b": 5,
	}).WithCustomWriter(buff)

	kvData := ""

	testFields := func(out string) { // This is nescessary since the range orders can randomly change
		for k, v := range i.fields {
			kvz := fmt.Sprintf("%s=%v,", k, v)
			if strings.Index(out, kvz) == -1 {
				t.Errorf("Expected \"%s\" in output: \"%s\"", kvData, out)
			}
		}
	}

	i.Info("Test %s %d %f %v", "huebr", 1, 10.0, true)

	o := buff.String()
	testFields(o)

	buff.Reset()

	i.Warn("Test %s %d %f %v", "huebr", 1, 10.0, true)

	o = buff.String()
	testFields(o)

	buff.Reset()

	i.Debug("Test %s %d %f %v", "huebr", 1, 10.0, true)

	o = buff.String()
	testFields(o)

	buff.Reset()

	i.Error("Test %s %d %f %v", "huebr", 1, 10.0, true)

	o = buff.String()
	testFields(o)

	buff.Reset()
}

func TestWithFieldsNoFields(t *testing.T) {
	SetFieldRepresentation(NoFields)
	buff := bytes.NewBufferString("")
	i := Scope("WithFieldsKV").WithFields(map[string]interface{}{
		"a": "b",
		"b": 5,
	}).WithCustomWriter(buff)

	kvData := ""

	for k, v := range i.fields {
		kvData += fmt.Sprintf("%s=%v,", k, v)
	}

	jsonDataB, _ := json.Marshal(i.fields)
	jsonData := string(jsonDataB)

	i.Info("Test %s %d %f %v", "huebr", 1, 10.0, true)

	o := buff.String()
	if strings.Index(o, kvData) != -1 {
		t.Errorf("Not expected \"%s\" in output: \"%s\"", kvData, o)
	}

	if strings.Index(o, jsonData) != -1 {
		t.Errorf("Not expected \"%s\" in output: \"%s\"", jsonData, o)
	}

	buff.Reset()

	i.Warn("Test %s %d %f %v", "huebr", 1, 10.0, true)

	o = buff.String()
	if strings.Index(o, kvData) != -1 {
		t.Errorf("Not expected \"%s\" in output: \"%s\"", kvData, o)
	}

	if strings.Index(o, jsonData) != -1 {
		t.Errorf("Not expected \"%s\" in output: \"%s\"", jsonData, o)
	}

	buff.Reset()

	i.Debug("Test %s %d %f %v", "huebr", 1, 10.0, true)

	o = buff.String()
	if strings.Index(o, kvData) != -1 {
		t.Errorf("Not expected \"%s\" in output: \"%s\"", kvData, o)
	}

	if strings.Index(o, jsonData) != -1 {
		t.Errorf("Not expected \"%s\" in output: \"%s\"", jsonData, o)
	}
	buff.Reset()

	i.Error("Test %s %d %f %v", "huebr", 1, 10.0, true)

	o = buff.String()
	if strings.Index(o, kvData) != -1 {
		t.Errorf("Not expected \"%s\" in output: \"%s\"", kvData, o)
	}

	if strings.Index(o, jsonData) != -1 {
		t.Errorf("Not expected \"%s\" in output: \"%s\"", jsonData, o)
	}

	buff.Reset()
}

func TestDebug(t *testing.T) {
	Debug("Test %s %d %f %v", "huebr", 1, 10.0, true) // Shouldn't crash
}

func TestError(t *testing.T) {
	Error("Test %s %d %f %v", "huebr", 1, 10.0, true) // Shouldn't crash
}

func TestWarn(t *testing.T) {
	Warn("Test %s %d %f %v", "huebr", 1, 10.0, true) // Shouldn't crash
}

func TestInfo(t *testing.T) {
	Info("Test %s %d %f %v", "huebr", 1, 10.0, true) // Shouldn't crash
}

func TestLog(t *testing.T) {
	Log("Test %s %d %f %v", "huebr", 1, 10.0, true) // Shouldn't crash
}

func TestLogNoFormat(t *testing.T) {
	LogNoFormat("Test %s %d %f %v", "huebr", 1, 10.0, true)
}

func TestFatal(t *testing.T) {
	assertPanic(t, func() {
		Fatal("Test Fatal")
	}, "Fatal should panic")
	assertPanic(t, func() {
		Fatal("Test %s %d %f %v", "huebr", 1, 10.0, true)
	}, "Fatal should panic")
}

func TestScope(t *testing.T) {
	scoped := Scope("test-scope")
	if scoped.scope != "test-scope" {
		t.Fatalf("Expected test-scope got %s", scoped.scope)
	}
}

func TestSetDebug(t *testing.T) {
	SetDebug(true)
	if !DebugEnabled() {
		t.Fatalf("Debug is set to false! Should be true")
	}
	SetDebug(false)
	if DebugEnabled() {
		t.Fatalf("Debug is set to true! Should be false")
	}
}

func TestSetError(t *testing.T) {
	SetError(true)
	if !ErrorEnabled() {
		t.Fatalf("Error is set to false! Should be true")
	}
	SetError(false)
	if ErrorEnabled() {
		t.Fatalf("Error is set to true! Should be false")
	}
}

func TestSetInfo(t *testing.T) {
	SetInfo(true)
	if !InfoEnabled() {
		t.Fatalf("Info is set to false! Should be true")
	}
	SetInfo(false)
	if InfoEnabled() {
		t.Fatalf("Info is set to true! Should be false")
	}
}

func TestSetWarn(t *testing.T) {
	SetWarning(true)
	if !WarningEnabled() {
		t.Fatalf("Warn is set to false! Should be true")
	}
	SetWarning(false)
	if WarningEnabled() {
		t.Fatalf("Warn is set to true! Should be false")
	}
}

func TestSetShowLines(t *testing.T) {
	SetShowLines(true)
	if !ShowLinesEnabled() {
		t.Fatal("ShowLines is set to false! Should be true")
	}
	SetShowLines(false)
	if ShowLinesEnabled() {
		t.Fatal("ShowLines is set to true! Should be false")
	}
}

type test struct{}

func (test) String() string {
	return "test"
}

func TestAsString(t *testing.T) {
	var tstringcast StringCast

	tstringcast = &test{}

	tests := []interface{}{
		"string",
		123456,
		123456.1,
		true,
		map[string]string{},
		[]int{1, 2, 3, 4, 5},
		complex(float32(1), float32(1)),
		complex(float64(1), float64(1)),
		fmt.Errorf("error format"),
		tstringcast,
	}

	outputs := make([]string, len(tests))

	for i, v := range tests { // Fill tests
		outputs[i] = fmt.Sprint(v) // Should be same output
	}

	for i, v := range tests {
		s := asString(v)
		if s != outputs[i] {
			t.Errorf("#%d expected %s got %s.", i, outputs[i], s)
		}
	}
}

func assertPanic(t *testing.T, f func(), message string) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(message)
		}
	}()
	f()
}
