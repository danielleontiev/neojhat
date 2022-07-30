package threaddump_test

import (
	_ "embed"

	"github.com/danielleontiev/neojhat/summary"
)

var summary1 = summary.Summary{
	Properties: map[string]string{
		"awt.toolkit": "sun.lwawt.macosx.LWCToolkit",
	},
	Env: map[string]string{
		"Architecture": "x86_64",
		"JavaHome":     "/Library/Java/JavaVirtualMachines/temurin-11.jdk/Contents/Home",
		"JavaName":     "OpenJDK 64-Bit Server VM (11.0.12+7, mixed mode)",
		"JavaVendor":   "Eclipse Foundation",
		"JavaVersion":  "11.0.12",
		"System":       "Mac OS X",
	},
	Heap: map[string]string{
		"Classes":   "42",
		"GC Roots":  "43",
		"Heap Size": "44",
		"Instances": "45",
	},
	System: map[string]string{
		"JVM Uptime": "40s",
	},
}

//go:embed case1.txt
var out1 string
