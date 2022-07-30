package threaddump_test

import (
	"strings"
	"testing"

	"github.com/danielleontiev/neojhat/threaddump"
)

func TestPlainCase1(t *testing.T) {
	builder := &strings.Builder{}
	threaddump.PrettyPrint(td1, true, builder)
	result := builder.String()
	if result != out1 {
		t.Errorf("PrettyPrint = \n---\n%s\n---\n, expected \n---\n%s\n---\n", result, out1)
	}
}

func TestPlainCase1NoLocal(t *testing.T) {
	builder := &strings.Builder{}
	threaddump.PrettyPrint(td1, false, builder)
	result := builder.String()
	if result != out1noLocal {
		t.Errorf("PrettyPrint = \n---\n%s\n---\n, expected \n---\n%s\n---\n", result, out1noLocal)
	}
}
