package output_test

import (
	"strings"
	"testing"
)

func compareLineByLine(t *testing.T, left, right string) {
	leftStrings := strings.Split(left, "\n")
	rightStrings := strings.Split(right, "\n")

	if len(leftStrings) != len(rightStrings) {
		t.Errorf("len(left) != len(right), %v != %v", len(leftStrings), len(rightStrings))
		t.Errorf("\n---\n%s\n---\n, expected \n---\n%s\n---\n", left, right)
		return
	}

	for i, l := range leftStrings {
		r := rightStrings[i]
		lTrim := strings.TrimSpace(l)
		rTrim := strings.TrimSpace(r)
		if lTrim != rTrim {
			t.Errorf("Lines not equal:\n-\n%s\n+\n%s\n", lTrim, rTrim)
			return
		}
	}
}
