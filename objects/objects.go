package objects

import (
	"github.com/danielleontiev/neojhat/core"
	"github.com/danielleontiev/neojhat/dump"
	"github.com/danielleontiev/neojhat/format"
)

func GetObjects(parserAccessor *dump.ParsedAccessor, sortBy SortBy) (Objects, error) {
	meta := parserAccessor.MetaStorage
	sizeInfo := core.NewSizeInfo(parserAccessor.IdSize)

	var totalSize, totalCount int
	var items []ObjectItem
	for arrType, instancesCount := range meta.Counters.PrimArraysCount {
		elementsCount := meta.Counters.PrimArrayElementsCount[arrType]
		name := arrType.String() + "[]"
		size := sizeInfo.OfType(arrType) * elementsCount
		totalSize += size
		totalCount += instancesCount
		items = append(items,
			ObjectItem{Name: name, InstancesCount: instancesCount, TotalSize: size})
	}
	for classId, instancesCount := range meta.Counters.ObjArraysCount {
		elementsCount := meta.Counters.ObjArrayElementsCount[classId]
		loadClass, err := parserAccessor.GetHprofLoadClassByClassObjectId(classId)
		if err != nil {
			return Objects{}, err
		}
		className, err := parserAccessor.GetHprofUtf8(loadClass.ClassNameId)
		if err != nil {
			return Objects{}, err
		}
		name, _ := format.Signature(className.Characters)
		size := sizeInfo.OfType(core.Object) * elementsCount
		totalSize += size
		totalCount += instancesCount
		items = append(items,
			ObjectItem{Name: name, InstancesCount: instancesCount, TotalSize: size})
	}
	for classId, instancesCount := range meta.Counters.InstancesCount {
		loadClass, err := parserAccessor.GetHprofLoadClassByClassObjectId(classId)
		if err != nil {
			return Objects{}, err
		}
		className, err := parserAccessor.GetHprofUtf8(loadClass.ClassNameId)
		if err != nil {
			return Objects{}, err
		}
		classDump, err := parserAccessor.GetHprofGcClassDump(classId)
		if err != nil {
			return Objects{}, err
		}
		name := className.Characters
		size := int(classDump.InstanceSize) * instancesCount
		totalSize += size
		totalCount += instancesCount
		items = append(items,
			ObjectItem{Name: name, InstancesCount: instancesCount, TotalSize: size})
	}

	return Objects{
		Items:      items,
		TotalSize:  totalSize,
		TotalCount: totalCount,
		SortBy:     sortBy,
	}, nil
}
