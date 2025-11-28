package output_test

import (
	_ "embed"

	"strings"
	"testing"

	"github.com/danielleontiev/neojhat/internal/output"
	"github.com/danielleontiev/neojhat/internal/summary"
)

var summary1 = summary.Summary{
	Properties: []summary.Kv{
		{Key: "awt.toolkit", Val: "sun.lwawt.macosx.LWCToolkit"},
	},
	Env: summary.EnvProperties{
		Architecture: "x86_64",
		JavaHome:     "/Library/Java/JavaVirtualMachines/temurin-11.jdk/Contents/Home",
		JavaName:     "OpenJDK 64-Bit Server VM (11.0.12+7, mixed mode)",
		JavaVendor:   "Eclipse Foundation",
		JavaVersion:  "11.0.12",
		System:       "Mac OS X",
	},
	Heap: summary.HeapProperties{
		Classes:   42,
		GcRoots:   43,
		HeapSize:  44,
		Instances: 45,
	},
	System: summary.SystemProperties{
		JvmUptime: "40s",
	},
}

var (
	//go:embed test-data/summary1.txt
	summary1txt string
	//go:embed test-data/summary1.html
	summary1html string
)

func TestSummaryPlain1(t *testing.T) {
	builder := &strings.Builder{}
	output.SummaryPlain(summary1, builder)
	result := builder.String()
	if result != summary1txt {
		compareLineByLine(t, result, summary1txt)
	}
}

func TestSummaryHtml1(t *testing.T) {
	builder := &strings.Builder{}
	output.SummaryHtml(summary1, builder)
	result := builder.String()
	compareLineByLine(t, summary1html, result)
}
