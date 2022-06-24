package summary

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/danielleontiev/neojhat/printing"
)

// PrettyPrint prints given Properties
// in beautiful manner
func PrettyPrint(s Summary) {
	identity := func(s string) string { return s }
	printSummary(s, identity, identity, identity)
}

// PrettyPrintColor prints given Properties
// in beautiful manner with colors
func PrettyPrintColor(s Summary) {
	printSummary(s, printing.Blue, printing.Bold, printing.Cyan)
}

func printSummary(s Summary, titleColor, keyColor, valColor func(string) string) {
	maxKeyLen := maxKeyLength(s)
	env := titleColor("- Environment")
	heap := titleColor("- Heap")
	system := titleColor("- System")
	props := titleColor("- Properties")
	fmt.Println(env)
	printProperties(s.env, maxKeyLen, keyColor, valColor)
	fmt.Println()
	fmt.Println(heap)
	printProperties(s.heap, maxKeyLen, keyColor, valColor)
	fmt.Println()
	fmt.Println(system)
	printProperties(s.system, maxKeyLen, keyColor, valColor)
	if s.properties != nil {
		fmt.Println()
		fmt.Println(props)
		printProperties(s.properties, maxKeyLen, keyColor, valColor)
	}
}

func printProperties(properties Properties, maxKeyLen int, keyColor, valColor func(string) string) {
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
		fmt.Print(s)
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
	return int(math.Max(maxFromProps(s.env), math.Max(maxFromProps(s.heap), maxFromProps(s.properties))))
}
