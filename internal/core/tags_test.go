package core

import "testing"

func TestTag_String(t *testing.T) {
	tests := []struct {
		name string
		tag  Tag
		want string
	}{
		{
			name: "HprofUtf8Tag",
			tag:  HprofUtf8Tag,
			want: "HPROF_UTF8",
		},
		{
			name: "HprofLoadClassTag",
			tag:  HprofLoadClassTag,
			want: "HPROF_LOAD_CLASS",
		},
		{
			name: "HprofFrameTag",
			tag:  HprofFrameTag,
			want: "HPROF_FRAME",
		},
		{
			name: "HprofTraceTag",
			tag:  HprofTraceTag,
			want: "HPROF_TRACE",
		},
		{
			name: "HprofHeapDumpSegmentTag",
			tag:  HprofHeapDumpSegmentTag,
			want: "HPROF_HEAP_DUMP_SEGMENT",
		},
		{
			name: "HprofHeapDumpEndTag",
			tag:  HprofHeapDumpEndTag,
			want: "HPROF_HEAP_DUMP_END",
		},
		{
			name: "UnknownTag",
			tag:  142,
			want: "UNKNOWN_TAG (142)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tag.String(); got != tt.want {
				t.Errorf("Tag.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubRecordType_String(t *testing.T) {
	tests := []struct {
		name string
		s    SubRecordType
		want string
	}{
		{
			name: "HprofGcRootJniGlobalType",
			s:    HprofGcRootJniGlobalType,
			want: "HPROF_GC_ROOT_JNI_GLOBAL",
		},
		{
			name: "HprofGcRootJniLocalType",
			s:    HprofGcRootJniLocalType,
			want: "HPROF_GC_ROOT_JNI_LOCAL",
		},
		{
			name: "HprofGcRootJavaFrameType",
			s:    HprofGcRootJavaFrameType,
			want: "HPROF_GC_ROOT_JAVA_FRAME",
		},
		{
			name: "HprofGcRootStickyClassType",
			s:    HprofGcRootStickyClassType,
			want: "HPROF_GC_ROOT_STICKY_CLASS",
		},
		{
			name: "HprofGcRootThreadObjType",
			s:    HprofGcRootThreadObjType,
			want: "HPROF_GC_ROOT_THREAD_OBJ",
		},
		{
			name: "HprofGcClassDumpType",
			s:    HprofGcClassDumpType,
			want: "HPROF_GC_CLASS_DUMP",
		},
		{
			name: "HprofGcInstanceDumpType",
			s:    HprofGcInstanceDumpType,
			want: "HPROF_GC_INSTANCE_DUMP",
		},
		{
			name: "HprofGcObjArrayDumpType",
			s:    HprofGcObjArrayDumpType,
			want: "HPROF_GC_OBJ_ARRAY_DUMP",
		},
		{
			name: "HprofGcPrimArrayDumpType",
			s:    HprofGcPrimArrayDumpType,
			want: "HPROF_GC_PRIM_ARRAY_DUMP",
		},
		{
			name: "Unknown",
			s:    42,
			want: "UNKNOWN_SUBRECORD_TYPE (42)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("SubRecordType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
