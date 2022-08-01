package output_test

import (
	_ "embed"

	"strings"
	"testing"

	"github.com/danielleontiev/neojhat/objects"
	"github.com/danielleontiev/neojhat/output"
)

var objects1 = objects.Objects{
	TotalSize:  2100000,
	TotalCount: 73224,
	SortBy:     objects.Size,
	Items: []objects.ObjectItem{
		{
			Name:           "byte[]",
			TotalSize:      571000,
			InstancesCount: 16558,
		},
		{
			Name:           "java.lang.Object[]",
			TotalSize:      340000,
			InstancesCount: 3141,
		},
		{
			Name:           "java.lang.String",
			TotalSize:      203000,
			InstancesCount: 16004,
		},
		{
			Name:           "java.lang.reflect.Method",
			TotalSize:      142000,
			InstancesCount: 1119,
		},
		{
			Name:           "java.util.HashMap$Node",
			TotalSize:      132000,
			InstancesCount: 4851,
		},
	},
}

var (
	//go:embed test-data/objects1.txt
	objects1txt string
	//go:embed test-data/objects1.html
	objects1html string
)

func TestObjectsPlain1(t *testing.T) {
	builder := &strings.Builder{}
	output.ObjectsPlain(objects1, builder)
	result := builder.String()
	if result != objects1txt {
		compareLineByLine(t, result, objects1txt)
	}
}

func TestObjectsHtml1(t *testing.T) {
	builder := &strings.Builder{}
	output.ObjectsHtml(objects1, builder)
	result := builder.String()
	if result != objects1html {
		compareLineByLine(t, result, objects1html)
	}
}
