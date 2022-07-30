package summary

import (
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/danielleontiev/neojhat/printing"
)

// PrettyPrint prints given Properties
// in beautiful manner
func PrettyPrint(s Summary, destination io.Writer) {
	identity := func(s string) string { return s }
	printSummary(s, identity, identity, identity, destination)
}

// PrettyPrintColor prints given Properties
// in beautiful manner with colors
func PrettyPrintColor(s Summary) {
	printSummary(s, printing.Blue, printing.Bold, printing.Cyan, os.Stdout)
}

func printSummary(s Summary, titleColor, keyColor, valColor func(string) string, destination io.Writer) {
	maxKeyLen := maxKeyLength(s)
	env := titleColor("- Environment")
	heap := titleColor("- Heap")
	system := titleColor("- System")
	props := titleColor("- Properties")
	fmt.Fprintln(destination, env)
	printProperties(s.Env, maxKeyLen, keyColor, valColor, destination)
	fmt.Fprintln(destination)
	fmt.Fprintln(destination, heap)
	printProperties(s.Heap, maxKeyLen, keyColor, valColor, destination)
	fmt.Fprintln(destination)
	fmt.Fprintln(destination, system)
	printProperties(s.System, maxKeyLen, keyColor, valColor, destination)
	if s.Properties != nil {
		fmt.Fprintln(destination)
		fmt.Fprintln(destination, props)
		printProperties(s.Properties, maxKeyLen, keyColor, valColor, destination)
	}
}

func printProperties(properties Properties, maxKeyLen int, keyColor, valColor func(string) string, destination io.Writer) {
	const spaceCount = 10
	var keys []string
	for k := range properties {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := strings.TrimSpace(properties[k])
		spaceSize := maxKeyLen + spaceCount - len(k)
		s := keyColor(k) + ":" + strings.Repeat(" ", spaceSize) + valColor(v) + "\n"
		fmt.Fprint(destination, s)
	}
}

func maxKeyLength(s Summary) int {
	maxFromProps := func(properties Properties) float64 {
		var maxKeyLen int
		for k := range properties {
			l := len(k)
			if l > maxKeyLen {
				maxKeyLen = l
			}
		}
		return float64(maxKeyLen)
	}
	return int(math.Max(maxFromProps(s.Env), math.Max(maxFromProps(s.Heap), maxFromProps(s.Properties))))
}
