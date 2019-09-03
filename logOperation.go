package slog

import "github.com/logrusorgru/aurora"

type LogOperation string

const (
	MSG   LogOperation = "MSG"
	IO                 = "IO"
	AWAIT              = "AWAIT"
	DONE               = "DONE"
	NOTE               = "NOTE"
)

var operationColors = map[LogOperation]colorFunc{
	MSG:   aurora.BgBlack,
	IO:    aurora.BgMagenta,
	AWAIT: aurora.BgCyan,
	DONE:  aurora.BgGreen,
	NOTE:  aurora.BgBlack,
}

var maxOperationStringLength = int(0)

func init() {
	// Compute max length between LogOperation strings
	operations := []LogOperation{MSG, IO, AWAIT, DONE, NOTE}
	for _, v := range operations {
		if len(v) > maxOperationStringLength {
			maxOperationStringLength = len(v)
		}
	}
}
