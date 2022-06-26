package objects

import (
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/danielleontiev/neojhat/format"
	"github.com/danielleontiev/neojhat/printing"
)

func PrettyPrint(objects Objects) {
	id := func(s string) string { return s }
	print(objects, id, id, id, id)
}

func PrettyPrintColor(objects Objects) {
	print(objects, printing.Bold, printing.Cyan, printing.Yellow, printing.Blue)
}

func print(objects Objects, headerColor, summaryColor, classNameColor, numColor func(s string) string) {
	var printItems []printItem
	switch objects.SortBy {
	case Size:
		printItems = append(printItems, printItem{
			Name: "Class Name", InstancesCount: "Count", TotalSize: "Size ↓",
		})
		sort.Slice(objects.Items, func(i, j int) bool { return objects.Items[i].TotalSize > objects.Items[j].TotalSize })
	case Count:
		printItems = append(printItems, printItem{
			Name: "Class Name", InstancesCount: "Count ↓", TotalSize: "Size",
		})
		sort.Slice(objects.Items, func(i, j int) bool { return objects.Items[i].InstancesCount > objects.Items[j].InstancesCount })
	}
	for _, item := range objects.Items {
		printItems = append(printItems, printItem{
			Name:           format.ClassName(item.Name),
			InstancesCount: fmt.Sprintf("%v (%v%%)", item.InstancesCount, 100*item.InstancesCount/objects.TotalCount),
			TotalSize:      fmt.Sprintf("%v (%v%%)", format.Size(item.TotalSize), 100*item.TotalSize/objects.TotalSize),
		})
	}
	const gap = 10
	var maxName, maxCount, maxSize int
	for _, item := range printItems {
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

	instances := fmt.Sprintf("Instances: %v", objects.TotalCount)
	size := fmt.Sprintf("Total Size: %v", format.Size(objects.TotalSize))
	fmt.Println(summaryColor(instances))
	fmt.Println(summaryColor(size))
	fmt.Println()

	headerItem := printItems[0]
	header := alignLeft(headerItem.Name, maxName) + " |" + alignRight(headerItem.InstancesCount, maxCount) +
		" |" + alignRight(headerItem.TotalSize, maxSize) + " |"
	fmt.Println(headerColor(header))
	fmt.Println(strings.Repeat("-", gap*3+maxName+maxCount+maxSize+6))
	for _, item := range printItems[1:] {
		fmt.Println(stringifyItem(item))
	}
}
