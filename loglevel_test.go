package slog

import (
	"testing"
)

func TestGetDescription(t *testing.T) {
	var LA LogLevel = "LA"

	testCases := []struct {
		name                string
		level               LogLevel
		expectedDescription string
	}{
		{
			name:                "input is INFO",
			level:               INFO,
			expectedDescription: "info",
		},
		{
			name:                "input is WARN",
			level:               WARN,
			expectedDescription: "warn",
		},
		{
			name:                "input is ERROR",
			level:               ERROR,
			expectedDescription: "error",
		},
		{
			name:                "input is FATAL",
			level:               FATAL,
			expectedDescription: "fatal",
		},
		{
			name:                "input is DEBUG",
			level:               DEBUG,
			expectedDescription: "debug",
		},
		{
			name:                "input with no description",
			level:               LA,
			expectedDescription: "LA",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := getDescription(tc.level)
			if result != tc.expectedDescription {
				t.Errorf("Got %q want %q.", result, tc.expectedDescription)
			}
		})
	}
}
