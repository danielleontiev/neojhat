package threaddump_test

import (
	_ "embed"

	"github.com/danielleontiev/neojhat/objects"
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

//go:embed case1.txt
var out1 string
