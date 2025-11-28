package core

import "testing"

const idSize = 8

var size = NewSizeInfo(idSize)

func TestHprofGcRootThreadObj_Size(t *testing.T) {
	r := HprofGcRootThreadObj{
		ThreadObjectId:           1,
		ThreadSequenceNumber:     1,
		StackTraceSequenceNumber: 1,
	}
	var want int = 16
	if got := size.Of(r); got != want {
		t.Errorf("HprofGcRootThreadObj.Size() = %v, want %v", got, want)
	}
}

func TestHprofGcRootJniGlobal_Size(t *testing.T) {
	r := HprofGcRootJniGlobal{
		ObjectId:       1,
		JniGlobalRefId: 1,
	}
	var want int = 16
	if got := size.Of(r); got != want {
		t.Errorf("HprofGcRootJniGlobal.Size() = %v, want %v", got, want)
	}
}

func TestHprofGcRootJniLocal_Size(t *testing.T) {
	r := HprofGcRootJniLocal{
		ObjectId:                1,
		ThreadSerialNumber:      1,
		FrameNumberInStackTrace: 1,
	}
	var want int = 16
	if got := size.Of(r); got != want {
		t.Errorf("HprofGcRootJniLocal.Size() = %v, want %v", got, want)
	}
}

func TestHprofGcRootJavaFrame_Size(t *testing.T) {
	r := HprofGcRootJavaFrame{
		ObjectId:                1,
		ThreadSerialNumber:      1,
		FrameNumberInStackTrace: 1,
	}
	var want int = 16
	if got := size.Of(r); got != want {
		t.Errorf("HprofGcRootJavaFrame.Size() = %v, want %v", got, want)
	}
}

func TestHprofGcRootStickyClass_Size(t *testing.T) {
	r := HprofGcRootStickyClass{
		ObjectId: 1,
	}
	var want int = 8
	if got := size.Of(r); got != want {
		t.Errorf("HprofGcRootStickyClass.Size() = %v, want %v", got, want)
	}
}

func TestHprofGcClassDump_Size(t *testing.T) {
	r := HprofGcClassDump{
		ClassObjectId:            1,
		StackTraceSerialNumber:   1,
		SuperclassObjectId:       1,
		ClassloaderObjectId:      1,
		SignersObjectId:          1,
		ProtectionDomainObjectId: 1,
		InstanceSize:             1,
		SizeOfConstantPool:       1,
		ConstantPoolRecords: []HprofGcClassDumpConstantPoolRecord{
			{ConstantPoolIndex: 1, Ty: Boolean, Value: JavaValue{Type: Boolean, Value: true}},
		},
		NumberOfStaticFields: 1,
		StaticFieldRecords: []HprofGcClassDumpStaticFieldsRecord{
			{StaticFieldName: 1, Ty: Boolean, Value: JavaValue{Type: Boolean, Value: true}},
		},
		NumberOfInstanceFields: 1,
		InstanceFieldRecords: []HprofGcClassDumpInstanceFieldsRecord{
			{InstanceFieldName: 1, Ty: Boolean},
		},
	}
	var want int = 93
	if got := size.Of(r); got != want {
		t.Errorf("HprofGcClassDump.Size() = %v, want %v", got, want)
	}
}

func TestHprofGcClassDumpInstanceDumpHeader_Size(t *testing.T) {
	r := HprofGcClassDumpInstanceDumpHeader{
		ObjectId:                1,
		StackTraceSerialNumber:  1,
		ClassObjectId:           1,
		NumberOfBytesThatFollow: 1,
	}
	var wantFullSize int = 25
	var wantRecordsSize int = 1
	if gotFullSize, gotRecordsSize := size.OfObject(r); (gotFullSize != wantFullSize) || (gotRecordsSize != wantRecordsSize) {
		t.Errorf("HprofGcClassDumpInstanceDumpHeader.Size() = (%v, %v), want (%v, %v)", gotFullSize, gotRecordsSize, wantFullSize, wantRecordsSize)
	}
}

func TestHprofGcObjArrayDumpHeader_Size(t *testing.T) {
	r := HprofGcObjArrayDumpHeader{
		ArrayObjectId:          1,
		StackTraceSerialNumber: 1,
		NumberOfElements:       1,
		ArrayClassId:           1,
	}
	var wantFullSize int = 32
	var wantRecordsSize int = 8
	if gotFullSize, gotRecordsSize := size.OfObject(r); (gotFullSize != wantFullSize) || (gotRecordsSize != wantRecordsSize) {
		t.Errorf("HprofGcObjArrayDumpHeader.Size() = (%v, %v), want (%v, %v)", gotFullSize, gotRecordsSize, wantFullSize, wantRecordsSize)
	}
}

func TestHprofGcPrimArrayDumpHeader_Size(t *testing.T) {
	r := HprofGcPrimArrayDumpHeader{
		ArrayObjectId:          1,
		StackTraceSerialNumber: 1,
		NumberOfElements:       1,
		ElementType:            Boolean,
	}
	var wantFullSize int = 18
	var wantRecordsSize int = 1
	if gotFullSize, gotRecordsSize := size.OfObject(r); (gotFullSize != wantFullSize) || (gotRecordsSize != wantRecordsSize) {
		t.Errorf("HprofGcPrimArrayDumpHeader.Size() = (%v, %v), want (%v, %v)", gotFullSize, gotRecordsSize, wantFullSize, wantRecordsSize)
	}
}

func TestJavaType_Size(t *testing.T) {
	tests := []struct {
		name string
		j    JavaType
		want int
	}{
		{
			name: "Object",
			j:    Object,
			want: 8,
		},
		{
			name: "Byte",
			j:    Byte,
			want: 1,
		},
		{
			name: "Boolean",
			j:    Boolean,
			want: 1,
		},
		{
			name: "Char",
			j:    Char,
			want: 2,
		},
		{
			name: "Short",
			j:    Short,
			want: 2,
		},
		{
			name: "Float",
			j:    Float,
			want: 4,
		},
		{
			name: "Int",
			j:    Int,
			want: 4,
		},
		{
			name: "Double",
			j:    Double,
			want: 8,
		},
		{
			name: "Long",
			j:    Long,
			want: 8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := size.OfType(tt.j); got != tt.want {
				t.Errorf("JavaType.Size() = %v, want %v", got, tt.want)
			}
		})
	}
}
