package slog

import (
	"github.com/logrusorgru/aurora"
	"testing"
)

func TestStripColors(t *testing.T) {
	const originalString = "COLOR" + "BOLD" + "BOLD+COLOR"
	var ansiString = aurora.BgGreen("COLOR").Yellow().String() + aurora.Bold("BOLD").String() + aurora.BgGreen("BOLD+COLOR").Bold().String()

	if stripColors(ansiString) != originalString {
		t.Errorf("Expected %s == %s", stripColors(ansiString), originalString)
	}
}

func TestStringSliceIndexOf(t *testing.T) {
	s := []string{"A", "B", "HUEBR", "PPPPPPPP", "AISUEho13h19h39 h"}
	// Found

	for i, v := range s {
		if stringSliceIndexOf(v, s) != i {
			t.Errorf("Expected %s to be at index %d but found %d", v, i, stringSliceIndexOf(v, s))
		}
	}

	// Not Found
	v := "193019237h9xn12j0ge182sk29x3mh02w9k21hw"
	if stringSliceIndexOf(v, s) != -1 {
		t.Errorf("Expected %s not found but got %d", v, stringSliceIndexOf(v, s))
	}
}
