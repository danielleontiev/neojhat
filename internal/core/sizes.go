package core

import "fmt"

type SizeInfo struct {
	idSize int
}

func NewSizeInfo(idSize uint32) *SizeInfo {
	return &SizeInfo{idSize: int(idSize)}
}

func (s SizeInfo) Of(record any) int {
	switch v := record.(type) {
	case HprofGcRootJniGlobal:
		return s.idSize * 2
	case HprofGcRootJniLocal:
		return s.idSize + 2*4
	case HprofGcRootJavaFrame:
		return s.idSize + 2*4
	case HprofGcRootStickyClass:
		return s.idSize
	case HprofGcRootThreadObj:
		return s.idSize + 2*4
	case HprofGcClassDump:
		fieldsSize := 7*s.idSize + 2*4 + 3*2
		var constantPoolSize int
		for _, cpr := range v.ConstantPoolRecords {
			constantPoolRecordSize := 2 + 1 + s.OfType(cpr.Ty)
			constantPoolSize += constantPoolRecordSize
		}
		var staticFieldsSize int
		for _, sfr := range v.StaticFieldRecords {
			staticFieldRecordSize := s.idSize + 1 + s.OfType(sfr.Ty)
			staticFieldsSize += staticFieldRecordSize
		}
		instanceFieldsSize := int(v.NumberOfInstanceFields) * (s.idSize + 1)
		return fieldsSize + constantPoolSize + staticFieldsSize + instanceFieldsSize
	default:
		panic(fmt.Sprintf("unexpected call of SizeInfo.Of(%v)", v))
	}
}

func (s SizeInfo) OfType(javaType JavaType) int {
	switch javaType {
	case Object:
		return s.idSize
	case Byte:
		return 1
	case Boolean:
		return 1
	case Char:
		return 2
	case Short:
		return 2
	case Float:
		return 4
	case Int:
		return 4
	case Double:
		return 8
	case Long:
		return 8
	}
	panic(fmt.Sprintf("unknown type %v", javaType))
}

func (s SizeInfo) OfObject(record any) (size, recordsSize int) {
	switch v := record.(type) {
	case HprofGcClassDumpInstanceDumpHeader:
		recordsSize = int(v.NumberOfBytesThatFollow)
		size = 2*s.idSize + 2*4 + recordsSize
	case HprofGcObjArrayDumpHeader:
		recordsSize = int(v.NumberOfElements) * s.idSize
		size = 2*s.idSize + 2*4 + recordsSize
	case HprofGcPrimArrayDumpHeader:
		recordsSize = int(v.NumberOfElements) * s.OfType(v.ElementType)
		size = s.idSize + 2*4 + 1 + recordsSize
	default:
		panic(fmt.Sprintf("unexpected call of SizeInfo.OfObject(%v)", v))
	}
	return
}
