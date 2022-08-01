package output

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/danielleontiev/neojhat/format"
	"github.com/danielleontiev/neojhat/objects"
)

type printObject struct {
	NameHeader           string
	InstancesCountHeader string
	TotalSizeHeader      string
	TotalCount           string
	TotalSize            string
	Items                []printItem
}
type printItem struct {
	Name           string
	TotalSize      string
	InstancesCount string
}

// ObjectsPlain print the result of objects command
// without colors
func ObjectsPlain(o objects.Objects, destination io.Writer) {
	id := func(s string) string { return s }
	print(o, id, id, id, id, destination)
}

// ObjectsPlainColor is the same as ObjectsPlain but
// with colorful output
func ObjectsPlainColor(o objects.Objects) {
	print(o, Bold, Cyan, Yellow, Blue, os.Stdout)
}

func print(o objects.Objects, headerColor, summaryColor, classNameColor, numColor func(s string) string, destination io.Writer) {
	printObj := getPrintItems(o)
	const gap = 10
	var maxName, maxCount, maxSize int
	for _, item := range append(printObj.Items, printItem{Name: printObj.NameHeader, TotalSize: printObj.TotalSizeHeader, InstancesCount: printObj.InstancesCountHeader}) {
		if len(item.Name) > maxName {
			maxName = len(item.Name)
		}
		if len(item.InstancesCount) > maxCount {
			maxCount = len(item.InstancesCount)
		}
		if len(item.TotalSize) > maxSize {
			maxSize = len(item.TotalSize)
		}
	}

	alignRight := func(s string, max int) string {
		return strings.Repeat(" ", max+gap-utf8.RuneCountInString(s)) + s
	}
	alignLeft := func(s string, max int) string {
		return s + strings.Repeat(" ", max+gap-utf8.RuneCountInString(s))
	}
	stringifyItem := func(i printItem) string {
		name := alignLeft(i.Name, maxName)
		count := alignRight(i.InstancesCount, maxCount)
		size := alignRight(i.TotalSize, maxSize)
		return classNameColor(name) + " |" + numColor(count) + " |" + numColor(size) + " |"
	}

	fmt.Fprintln(destination, summaryColor(fmt.Sprintf("Instances: %v", printObj.TotalCount)))
	fmt.Fprintln(destination, summaryColor(fmt.Sprintf("Total Size: %v", printObj.TotalSize)))
	fmt.Fprintln(destination)

	header := alignLeft(printObj.NameHeader, maxName) + " |" + alignRight(printObj.InstancesCountHeader, maxCount) +
		" |" + alignRight(printObj.TotalSizeHeader, maxSize) + " |"
	fmt.Fprintln(destination, headerColor(header))
	fmt.Fprintln(destination, strings.Repeat("-", gap*3+maxName+maxCount+maxSize+6))
	for _, item := range printObj.Items {
		fmt.Fprintln(destination, stringifyItem(item))
	}
}

func getPrintItems(o objects.Objects) printObject {
	var printObj printObject
	var printItems []printItem
	switch o.SortBy {
	case objects.Size:
		printObj = printObject{
			NameHeader: "Class Name", InstancesCountHeader: "Count", TotalSizeHeader: "Size ↓",
		}
		sort.Slice(o.Items, func(i, j int) bool {
			if o.Items[i].TotalSize == o.Items[j].TotalSize {
				return o.Items[i].Name < o.Items[j].Name
			}
			return o.Items[i].TotalSize > o.Items[j].TotalSize
		})
	case objects.Count:
		printObj = printObject{
			NameHeader: "Class Name", InstancesCountHeader: "Count ↓", TotalSizeHeader: "Size",
		}
		sort.Slice(o.Items, func(i, j int) bool {
			if o.Items[i].InstancesCount == o.Items[j].InstancesCount {
				return o.Items[i].Name < o.Items[j].Name
			}
			return o.Items[i].InstancesCount > o.Items[j].InstancesCount
		})
	}
	for _, item := range o.Items {
		printItems = append(printItems, printItem{
			Name:           format.ClassName(item.Name),
			InstancesCount: fmt.Sprintf("%v (%v%%)", item.InstancesCount, 100*item.InstancesCount/o.TotalCount),
			TotalSize:      fmt.Sprintf("%v (%v%%)", format.Size(item.TotalSize), 100*item.TotalSize/o.TotalSize),
		})
	}
	printObj.TotalCount = strconv.Itoa(o.TotalCount)
	printObj.TotalSize = format.Size(o.TotalSize)
	printObj.Items = printItems
	return printObj
}
