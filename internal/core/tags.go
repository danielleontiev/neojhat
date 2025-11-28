package core

import "fmt"

type Tag byte

const (
	HprofUtf8Tag            Tag = 0x01
	HprofLoadClassTag       Tag = 0x02
	HprofFrameTag           Tag = 0x04
	HprofTraceTag           Tag = 0x05
	HprofHeapDumpSegmentTag Tag = 0x1c
	HprofHeapDumpEndTag     Tag = 0x2c
)

var tagMap = map[Tag]string{
	HprofUtf8Tag:            "HPROF_UTF8",
	HprofLoadClassTag:       "HPROF_LOAD_CLASS",
	HprofFrameTag:           "HPROF_FRAME",
	HprofTraceTag:           "HPROF_TRACE",
	HprofHeapDumpSegmentTag: "HPROF_HEAP_DUMP_SEGMENT",
	HprofHeapDumpEndTag:     "HPROF_HEAP_DUMP_END",
}

func (t Tag) String() string {
	str, ok := tagMap[t]
	if ok {
		return str
	}
	return fmt.Sprintf("UNKNOWN_TAG (%v)", byte(t))
}

type SubRecordType byte

const (
	HprofGcRootJniGlobalType   SubRecordType = 0x01
	HprofGcRootJniLocalType    SubRecordType = 0x02
	HprofGcRootJavaFrameType   SubRecordType = 0x03
	HprofGcRootStickyClassType SubRecordType = 0x05
	HprofGcRootThreadObjType   SubRecordType = 0x08
	HprofGcClassDumpType       SubRecordType = 0x20
	HprofGcInstanceDumpType    SubRecordType = 0x21
	HprofGcObjArrayDumpType    SubRecordType = 0x22
	HprofGcPrimArrayDumpType   SubRecordType = 0x23
)

// these duplicated tags from Tags to be able to exit parse subrecords loop.
const (
	HprofHeapDumpEndSubRecord     SubRecordType = 0x2c
	HprofHeapDumpSegmentSubRecord SubRecordType = 0x1c
)

var subRecordTypeMap = map[SubRecordType]string{
	HprofGcRootJniGlobalType:   "HPROF_GC_ROOT_JNI_GLOBAL",
	HprofGcRootJniLocalType:    "HPROF_GC_ROOT_JNI_LOCAL",
	HprofGcRootJavaFrameType:   "HPROF_GC_ROOT_JAVA_FRAME",
	HprofGcRootStickyClassType: "HPROF_GC_ROOT_STICKY_CLASS",
	HprofGcRootThreadObjType:   "HPROF_GC_ROOT_THREAD_OBJ",
	HprofGcClassDumpType:       "HPROF_GC_CLASS_DUMP",
	HprofGcInstanceDumpType:    "HPROF_GC_INSTANCE_DUMP",
	HprofGcObjArrayDumpType:    "HPROF_GC_OBJ_ARRAY_DUMP",
	HprofGcPrimArrayDumpType:   "HPROF_GC_PRIM_ARRAY_DUMP",
}

func (s SubRecordType) String() string {
	str, ok := subRecordTypeMap[s]
	if ok {
		return str
	}
	return fmt.Sprintf("UNKNOWN_SUBRECORD_TYPE (%v)", byte(s))
}
