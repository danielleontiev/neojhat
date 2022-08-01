package output

import (
	_ "embed"

	"fmt"
	"html/template"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/danielleontiev/neojhat/format"
	"github.com/danielleontiev/neojhat/summary"
)

// SummaryPlain prints given Properties
// in beautiful manner
func SummaryPlain(s summary.Summary, destination io.Writer) {
	identity := func(s string) string { return s }
	printSummary(s, identity, identity, identity, destination)
}

// SummaryPlainColor prints given Properties
// in beautiful manner with colors
func SummaryPlainColor(s summary.Summary) {
	printSummary(s, Blue, Bold, Cyan, os.Stdout)
}

func printSummary(s summary.Summary, titleColor, keyColor, valColor func(string) string, destination io.Writer) {
	props := getPropertiesList(s)
	maxKeyLen := maxKeyLength(props)
	for _, p := range props {
		fmt.Fprintln(destination, "- "+p.Name)
		printProperties(p.Kv, maxKeyLen, keyColor, valColor, destination)
		fmt.Fprintln(destination)
	}
}

type Properties struct {
	Name string
	Kv   []summary.Kv
}

func getPropertiesList(s summary.Summary) []Properties {
	env := []summary.Kv{
		{Key: "Architecture", Val: s.Env.Architecture},
		{Key: "JavaHome", Val: s.Env.JavaHome},
		{Key: "JavaName", Val: s.Env.JavaName},
		{Key: "JavaVendor", Val: s.Env.JavaVendor},
		{Key: "JavaVersion", Val: s.Env.JavaVersion},
		{Key: "System", Val: s.Env.System},
	}
	heap := []summary.Kv{
		{Key: "Classes", Val: strconv.Itoa(s.Heap.Classes)},
		{Key: "GC Roots", Val: strconv.Itoa(s.Heap.GcRoots)},
		{Key: "Instances", Val: strconv.Itoa(s.Heap.Instances)},
		{Key: "Heap Size", Val: format.Size(s.Heap.HeapSize)},
	}
	system := []summary.Kv{
		{Key: "JVM Uptime", Val: s.System.JvmUptime},
	}
	properties := []Properties{
		{Name: "Environment", Kv: env},
		{Name: "Heap", Kv: heap},
		{Name: "System", Kv: system},
	}
	if s.Properties != nil {
		properties = append(properties, Properties{Name: "Properties", Kv: s.Properties})
	}
	return properties
}

func printProperties(properties []summary.Kv, maxKeyLen int, keyColor, valColor func(string) string, destination io.Writer) {
	const spaceCount = 10
	for _, kv := range properties {
		v := strings.TrimSpace(kv.Val)
		spaceSize := maxKeyLen + spaceCount - len(kv.Key)
		s := keyColor(kv.Key) + ":" + strings.Repeat(" ", spaceSize) + valColor(v) + "\n"
		fmt.Fprint(destination, s)
	}
}

func maxKeyLength(props []Properties) int {
	var properties []summary.Kv
	for _, p := range props {
		for _, v := range p.Kv {
			properties = append(properties, v)
		}
	}
	var maxKeyLen int
	for _, kv := range properties {
		l := len(kv.Key)
		if l > maxKeyLen {
			maxKeyLen = l
		}
	}
	return maxKeyLen
}

var (
	//go:embed templates/summary.html
	summaryHtml string
)

// SummaryHtml prints the output of summary command in nice
// beautifully-formatted HTML
func SummaryHtml(s summary.Summary, destination io.Writer) error {
	coreTemplate, err := template.New("core").Parse(coreHtml)
	if err != nil {
		return err
	}
	summaryTemplate, err := coreTemplate.Parse(summaryHtml)
	if err != nil {
		return err
	}
	props := getPropertiesList(s)
	return summaryTemplate.Execute(destination, data{
		Title:   "Heap Summary",
		Favicon: faviconBase64,
		Payload: props,
	})
}
