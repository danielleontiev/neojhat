package core

import (
	"reflect"
	"testing"
)

func TestRecordParser_ParseRecordHeader(t *testing.T) {
	hprofUtf8 := []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0d}
	tests := []struct {
		name          string
		parser        RecordParser
		wantHeader    RecordHeader
		wantRemaining uint32
		wantErr       bool
	}{
		{
			name:          "success parse record header",
			parser:        createRecordParser(hprofUtf8, CreateOpts{idSize: 8}),
			wantHeader:    RecordHeader{Tag: HprofUtf8Tag, Remaining: 13},
			wantRemaining: 13,
		},
		{
			name:          "error parse record header",
			parser:        createRecordParser(empty, CreateOpts{idSize: 8}),
			wantRemaining: 0,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHeader, err := tt.parser.ParseRecordHeader()
			if (err != nil) != tt.wantErr {
				t.Errorf("RecordParser.ParseRecordHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotHeader != tt.wantHeader {
				t.Errorf("RecordParser.ParseRecordHeader() = %v, want %v", gotHeader, tt.wantHeader)
			}
		})
	}
}

func TestRecordParser_ParseHprofUtf8(t *testing.T) {
	javaStr := []byte{0x4a, 0x41, 0x56, 0x41}
	tests := []struct {
		name      string
		parser    RecordParser
		remaining uint32
		want      HprofUtf8
		wantErr   bool
	}{
		{
			name:      "success",
			parser:    createRecordParser(concat(one4, javaStr), CreateOpts{idSize: 4}),
			remaining: 8,
			want:      HprofUtf8{Identifier: 1, Characters: "JAVA"},
		},
		{
			name:    "error",
			parser:  createRecordParser(empty, CreateOpts{idSize: 4}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.parser.ParseHprofUtf8(tt.remaining)
			if (err != nil) != tt.wantErr {
				t.Errorf("RecordParser.ParseHprofUtf8() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RecordParser.ParseHprofUtf8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecordParser_ParseHprofLoadClass(t *testing.T) {
	tests := []struct {
		name    string
		parser  RecordParser
		want    HprofLoadClass
		wantErr bool
	}{
		{
			name:   "success",
			parser: createRecordParser(concat(one4, one4, one4, one4), CreateOpts{idSize: 4}),
			want:   HprofLoadClass{ClassSerialNumber: 1, ClassObjectId: 1, StackTraceSerialNumber: 1, ClassNameId: 1},
		},
		{
			name:    "error",
			parser:  createRecordParser(empty, CreateOpts{}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.parser.ParseHprofLoadClass()
			if (err != nil) != tt.wantErr {
				t.Errorf("RecordParser.ParseHprofLoadClass() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RecordParser.ParseHprofLoadClass() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecordParser_ParseHprofFrame(t *testing.T) {
	tests := []struct {
		name    string
		parser  RecordParser
		want    HprofFrame
		wantErr bool
	}{
		{
			name:   "success",
			parser: createRecordParser(concat(one4, one4, one4, one4, one4, one4), CreateOpts{idSize: 4}),
			want: HprofFrame{
				StackFrameId:      1,
				MethodNameId:      1,
				MethodSignatureId: 1,
				SourceFileNameId:  1,
				ClassSerialNumber: 1,
				LineNumber:        1,
			},
		},
		{
			name:    "error",
			parser:  createRecordParser(empty, CreateOpts{idSize: 4}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.parser.ParseHprofFrame()
			if (err != nil) != tt.wantErr {
				t.Errorf("RecordParser.ParseHprofFrame() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RecordParser.ParseHprofFrame() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecordParser_ParseHprofTrace(t *testing.T) {
	tests := []struct {
		name    string
		parser  RecordParser
		want    HprofTrace
		wantErr bool
	}{
		{
			name:   "success",
			parser: createRecordParser(concat(one4, one4, one4, one4), CreateOpts{idSize: 4}),
			want: HprofTrace{
				StackTraceSerialNumber: 1,
				ThreadSerialNumber:     1,
				NumberOfFrames:         1,
				StackFrameIds:          []Identifier{1},
			},
		},
		{
			name:    "error",
			parser:  createRecordParser(empty, CreateOpts{}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.parser.ParseHprofTrace()
			if (err != nil) != tt.wantErr {
				t.Errorf("RecordParser.ParseHprofTrace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RecordParser.ParseHprofTrace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecordParser_parseHprofSubRecordHeader(t *testing.T) {
	tests := []struct {
		name    string
		parser  RecordParser
		want    SubRecordHeader
		wantErr bool
	}{
		{
			name:   "success",
			parser: createRecordParser([]byte{byte(HprofGcRootJniGlobalType)}, CreateOpts{}),
			want:   SubRecordHeader{SubRecordType: HprofGcRootJniGlobalType},
		},
		{
			name:    "error",
			parser:  createRecordParser(empty, CreateOpts{}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.parser.ParseSubRecordHeader()
			if (err != nil) != tt.wantErr {
				t.Errorf("RecordParser.ParseSubRecordHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RecordParser.ParseSubRecordHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecordParser_ParseHprofGcRootThreadObj(t *testing.T) {
	tests := []struct {
		name    string
		parser  RecordParser
		want    HprofGcRootThreadObj
		wantErr bool
	}{
		{
			name:   "success",
			parser: createRecordParser(concat(one4, one4, one4), CreateOpts{idSize: 4}),
			want: HprofGcRootThreadObj{
				ThreadObjectId:           1,
				ThreadSequenceNumber:     1,
				StackTraceSequenceNumber: 1,
			},
		},
		{
			name:    "error",
			parser:  createRecordParser(empty, CreateOpts{idSize: 4}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.parser.ParseHprofGcRootThreadObj()
			if (err != nil) != tt.wantErr {
				t.Errorf("RecordParser.ParseHprofGcRootThreadObj() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RecordParser.ParseHprofGcRootThreadObj() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecordParser_ParseHprofGcRootJniGlobal(t *testing.T) {
	tests := []struct {
		name    string
		parser  RecordParser
		want    HprofGcRootJniGlobal
		wantErr bool
	}{
		{
			name:   "success",
			parser: createRecordParser(concat(one4, one4), CreateOpts{idSize: 4}),
			want: HprofGcRootJniGlobal{
				ObjectId:       1,
				JniGlobalRefId: 1,
			},
		},
		{
			name:    "error",
			parser:  createRecordParser(empty, CreateOpts{idSize: 4}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.parser.ParseHprofGcRootJniGlobal()
			if (err != nil) != tt.wantErr {
				t.Errorf("RecordParser.ParseHprofGcRootJniGlobal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RecordParser.ParseHprofGcRootJniGlobal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecordParser_ParseHprofGcRootJniLocal(t *testing.T) {
	tests := []struct {
		name    string
		parser  RecordParser
		want    HprofGcRootJniLocal
		wantErr bool
	}{
		{
			name:   "success",
			parser: createRecordParser(concat(one4, one4, one4), CreateOpts{idSize: 4}),
			want: HprofGcRootJniLocal{
				ObjectId:                1,
				ThreadSerialNumber:      1,
				FrameNumberInStackTrace: 1,
			},
		},
		{
			name:    "error",
			parser:  createRecordParser(empty, CreateOpts{}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.parser.ParseHprofGcRootJniLocal()
			if (err != nil) != tt.wantErr {
				t.Errorf("RecordParser.ParseHprofGcRootJniLocal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RecordParser.ParseHprofGcRootJniLocal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecordParser_ParseHprofGcRootJavaFrame(t *testing.T) {
	tests := []struct {
		name    string
		parser  RecordParser
		want    HprofGcRootJavaFrame
		wantErr bool
	}{
		{
			name:   "success",
			parser: createRecordParser(concat(one4, one4, one4), CreateOpts{idSize: 4}),
			want: HprofGcRootJavaFrame{
				ObjectId:                1,
				ThreadSerialNumber:      1,
				FrameNumberInStackTrace: 1,
			},
		},
		{
			name:    "error",
			parser:  createRecordParser(empty, CreateOpts{}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.parser.ParseHprofGcRootJavaFrame()
			if (err != nil) != tt.wantErr {
				t.Errorf("RecordParser.ParseHprofGcRootJavaFrame() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RecordParser.ParseHprofGcRootJavaFrame() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecordParser_ParseHprofGcRootStickyClass(t *testing.T) {
	tests := []struct {
		name    string
		parser  RecordParser
		want    HprofGcRootStickyClass
		wantErr bool
	}{
		{
			name:   "success",
			parser: createRecordParser(one4, CreateOpts{idSize: 4}),
			want: HprofGcRootStickyClass{
				ObjectId: 1,
			},
		},
		{
			name:    "error",
			parser:  createRecordParser(empty, CreateOpts{}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.parser.ParseHprofGcRootStickyClass()
			if (err != nil) != tt.wantErr {
				t.Errorf("RecordParser.ParseHprofGcRootStickyClass() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RecordParser.ParseHprofGcRootStickyClass() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecordParser_ParseHprofGcClassDumpHeader(t *testing.T) {
	record := concat(
		one4, one4, one4, one4, one4, one4, one4, one4, one4, // header
		one2,                     // constant pool size
		one2, []byte{0x04}, one1, // constant pool
		one2,                     // number of static fields
		one4, []byte{0x04}, one1, // static fields
		one2,                       // inst fields number
		one4, []byte{byte(Object)}, // inst fields
	)
	tests := []struct {
		name    string
		parser  RecordParser
		want    HprofGcClassDump
		wantErr bool
	}{
		{
			name:   "success",
			parser: createRecordParser(record, CreateOpts{idSize: 4}),
			want: HprofGcClassDump{
				ClassObjectId:            1,
				StackTraceSerialNumber:   1,
				SuperclassObjectId:       1,
				ClassloaderObjectId:      1,
				SignersObjectId:          1,
				ProtectionDomainObjectId: 1,
				InstanceSize:             1,
				SizeOfConstantPool:       1,
				ConstantPoolRecords: []HprofGcClassDumpConstantPoolRecord{
					{
						ConstantPoolIndex: 1,
						Ty:                Boolean,
						Value:             JavaValue{Type: Boolean, Value: true},
					},
				},
				NumberOfStaticFields: 1,
				StaticFieldRecords: []HprofGcClassDumpStaticFieldsRecord{
					{
						StaticFieldName: 1,
						Ty:              Boolean,
						Value:           JavaValue{Type: Boolean, Value: true},
					},
				},
				NumberOfInstanceFields: 1,
				InstanceFieldRecords: []HprofGcClassDumpInstanceFieldsRecord{
					{InstanceFieldName: 1, Ty: Object},
				},
			},
		},
		{
			name:    "error",
			parser:  createRecordParser(empty, CreateOpts{idSize: 4}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.parser.ParseHprofGcClassDump()
			if (err != nil) != tt.wantErr {
				t.Errorf("RecordParser.ParseHprofGcClassDumpHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RecordParser.ParseHprofGcClassDumpHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecordParser_ParseHprofGcClassDumpConstantPoolRecord(t *testing.T) {
	tests := []struct {
		name    string
		parser  RecordParser
		want    HprofGcClassDumpConstantPoolRecord
		wantErr bool
	}{
		{
			name:   "success",
			parser: createRecordParser(concat(one2, []byte{0x04}, one1), CreateOpts{}),
			want: HprofGcClassDumpConstantPoolRecord{
				ConstantPoolIndex: 1,
				Ty:                Boolean,
				Value:             JavaValue{Type: Boolean, Value: true},
			},
		},
		{
			name:    "error",
			parser:  createRecordParser(empty, CreateOpts{}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.parser.parseHprofGcClassDumpConstantPoolRecord()
			if (err != nil) != tt.wantErr {
				t.Errorf("RecordParser.ParseHprofGcClassDumpConstantPoolRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RecordParser.ParseHprofGcClassDumpConstantPoolRecord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecordParser_ParseHprofGcClassDumpStaticFieldsRecord(t *testing.T) {
	tests := []struct {
		name    string
		parser  RecordParser
		want    HprofGcClassDumpStaticFieldsRecord
		wantErr bool
	}{
		{
			name:   "success",
			parser: createRecordParser(concat(one4, []byte{0x04}, one1), CreateOpts{idSize: 4}),
			want: HprofGcClassDumpStaticFieldsRecord{
				StaticFieldName: 1,
				Ty:              Boolean,
				Value:           JavaValue{Type: Boolean, Value: true},
			},
		},
		{
			name:    "error",
			parser:  createRecordParser(empty, CreateOpts{idSize: 4}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.parser.parseHprofGcClassDumpStaticFieldsRecord()
			if (err != nil) != tt.wantErr {
				t.Errorf("RecordParser.ParseHprofGcClassDumpStaticFieldsRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RecordParser.ParseHprofGcClassDumpStaticFieldsRecord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecordParser_ParseHprofGcClassDumpInstanceFieldsRecord(t *testing.T) {
	tests := []struct {
		name    string
		parser  RecordParser
		want    HprofGcClassDumpInstanceFieldsRecord
		wantErr bool
	}{
		{
			name:   "success",
			parser: createRecordParser(concat(one4, []byte{byte(Object)}), CreateOpts{idSize: 4}),
			want:   HprofGcClassDumpInstanceFieldsRecord{InstanceFieldName: 1, Ty: Object},
		},
		{
			name:    "error",
			parser:  createRecordParser(empty, CreateOpts{idSize: 4}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.parser.parseHprofGcClassDumpInstanceFieldsRecord()
			if (err != nil) != tt.wantErr {
				t.Errorf("RecordParser.ParseHprofGcClassDumpInstanceFieldsRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RecordParser.ParseHprofGcClassDumpInstanceFieldsRecord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecordParser_ParseHprofGcClassDumpInstanceDumpHeader(t *testing.T) {
	tests := []struct {
		name    string
		parser  RecordParser
		want    HprofGcClassDumpInstanceDumpHeader
		wantErr bool
	}{
		{
			name:   "success",
			parser: createRecordParser(concat(one4, one4, one4, one4), CreateOpts{idSize: 4}),
			want: HprofGcClassDumpInstanceDumpHeader{
				ObjectId:                1,
				StackTraceSerialNumber:  1,
				ClassObjectId:           1,
				NumberOfBytesThatFollow: 1,
			},
		},
		{
			name:    "error",
			parser:  createRecordParser(empty, CreateOpts{idSize: 4}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.parser.ParseHprofGcClassDumpInstanceDumpHeader()
			if (err != nil) != tt.wantErr {
				t.Errorf("RecordParser.ParseHprofGcClassDumpInstanceDumpHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RecordParser.ParseHprofGcClassDumpInstanceDumpHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecordParser_ParseHprofGcObjArrayDumpHeader(t *testing.T) {
	tests := []struct {
		name    string
		parser  RecordParser
		want    HprofGcObjArrayDumpHeader
		wantErr bool
	}{
		{
			name:   "success",
			parser: createRecordParser(concat(one4, one4, one4, one4), CreateOpts{idSize: 4}),
			want: HprofGcObjArrayDumpHeader{
				ArrayObjectId:          1,
				StackTraceSerialNumber: 1,
				NumberOfElements:       1,
				ArrayClassId:           1,
			},
		},
		{
			name:    "error",
			parser:  createRecordParser(empty, CreateOpts{idSize: 4}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.parser.ParseHprofGcObjArrayDumpHeader()
			if (err != nil) != tt.wantErr {
				t.Errorf("RecordParser.ParseHprofGcObjArrayDumpHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RecordParser.ParseHprofGcObjArrayDumpHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecordParser_ParseHprofGcPrimArrayDumpHeader(t *testing.T) {
	tests := []struct {
		name    string
		parser  RecordParser
		want    HprofGcPrimArrayDumpHeader
		wantErr bool
	}{
		{
			name:   "success",
			parser: createRecordParser(concat(one4, one4, one4, []byte{byte(Object)}), CreateOpts{idSize: 4}),
			want: HprofGcPrimArrayDumpHeader{
				ArrayObjectId:          1,
				StackTraceSerialNumber: 1,
				NumberOfElements:       1,
				ElementType:            Object,
			},
		},
		{
			name:    "success",
			parser:  createRecordParser(empty, CreateOpts{idSize: 4}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.parser.ParseHprofGcPrimArrayDumpHeader()
			if (err != nil) != tt.wantErr {
				t.Errorf("RecordParser.ParseHprofGcPrimArrayDumpHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RecordParser.ParseHprofGcPrimArrayDumpHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecordParser_ReadBytes(t *testing.T) {
	tests := []struct {
		name    string
		parser  RecordParser
		n       int
		want    []byte
		wantErr bool
	}{
		{
			name:   "success",
			parser: createRecordParser([]byte{0x00, 0x00, 0x00, 0x00}, CreateOpts{}),
			n:      2,
			want:   []byte{0x00, 0x00},
		},
		{
			name:    "error",
			parser:  createRecordParser([]byte{0x00, 0x00, 0x00, 0x00}, CreateOpts{}),
			n:       5,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.parser.ReadBytes(tt.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("RecordParser.ReadBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RecordParser.ReadBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
