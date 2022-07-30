package threaddump_test

import (
	"strings"
	"testing"

	"github.com/danielleontiev/neojhat/objects"
)

func TestPlainCase1(t *testing.T) {
	builder := &strings.Builder{}
	objects.PrettyPrint(objects1, builder)
	result := builder.String()
	if result != out1 {
		t.Errorf("PrettyPrint = \n---\n%s\n---\n, expected \n---\n%s\n---\n", result, out1)
	}
}
