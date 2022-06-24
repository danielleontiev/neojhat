package dump

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/danielleontiev/neojhat/core"
	"github.com/danielleontiev/neojhat/storage"
)

var (
	readerTestFileHeader = []byte{
		0x4a, 0x41, 0x56, 0x41, 0x20, 0x50, 0x52, 0x4f, 0x46, 0x49, 0x4c, 0x45, 0x20, 0x31, 0x2e, 0x30, 0x2e, 0x32, 0x00, // header, 0-terminated
		0x00, 0x00, 0x00, 0x08, // identifier size
		0x00, 0x00, 0x01, 0x7b, // timestamp, low word
		0x7f, 0x28, 0xa8, 0x27, // timestamp, high word
	}
)

var (
	testHeapDump = concat(
		// 0
		readerTestFileHeader,
		// 31
		createRecordHeader(core.HprofUtf8Tag, 8+4),
		// 40
		one8,                           // identifier
		[]byte{0x4a, 0x41, 0x56, 0x41}, // "JAVA" string
		// 52
		createRecordHeader(core.HprofLoadClassTag, 4+8+4+8),
		// 61
		one4, // class serial number
		one8, // class object id
		one4, // stack trace serial number
		one8, // class name id
		// 85
		createRecordHeader(core.HprofTraceTag, 4+4+4+3*8),
		// 94
		one4,   // stack trace serial number
		one4,   // thread serial number
		three4, // number of frames
		one8,   // frame id 1
		one8,   // frame id 2
		one8,   // frame id 3
		// 130
		createRecordHeader(core.HprofFrameTag, 8+8+8+8+4+4),
		// 139
		one8, // stack frame id
		one8, // method name id
		one8, // method signature id
		one8, // source file name id
		one4, // class serial number
		one4, // line number
		// 179
		createRecordHeader(core.HprofHeapDumpSegmentTag, 0), // remaining is not applicable here, hence 0
		// 188
		createSubRecordHeader(core.HprofGcRootThreadObjType),
		// 189
		one8, // thread object id
		one4, // thread sequence number
		one4, // stack trace sequence number
		// 205
		createSubRecordHeader(core.HprofGcRootJniLocalType),
		// 206
		one8, // object id
		one4, // thread serial number
		one4, // frame number in stack trace
		// 222
		createSubRecordHeader(core.HprofGcRootJniGlobalType),
		// 223
		one8, // object id
		one8, // jni global ref id
		// 239
		createSubRecordHeader(core.HprofGcRootJavaFrameType),
		// 240
		one8, // object id
		one4, // thread serial number
		one4, // frame number in stack trace
		// 256
		createSubRecordHeader(core.HprofGcRootStickyClassType),
		// 257
		one8, // object id
		// 265
		createSubRecordHeader(core.HprofGcClassDumpType),
		// 266
		one8, // class object id
		one4, // stack trace serial number
		one8, // super class object id
		one8, // class loader object id
		one8, // signers object id
		one8, // protection domain object id
		one8, // reserved
		one8, // reserved
		one4, // instance size
		two2, // size of constant pool
		// first constant pool record
		one2,                       // constant pool index
		[]byte{byte(core.Boolean)}, // type
		[]byte{0x00},               // value
		// second constant pool record
		one2,                       // constant pool index
		[]byte{byte(core.Boolean)}, // type
		[]byte{0x01},               // value
		two2,                       // number of static fields
		// first static field
		one8,                       // field name
		[]byte{byte(core.Boolean)}, // type
		[]byte{0x00},               // value
		// second static field
		one8,                       // field name
		[]byte{byte(core.Boolean)}, // type
		[]byte{0x01},               // value
		two2,                       // number of instance fields
		// first instance field
		one8,                      // field name
		[]byte{byte(core.Object)}, // type
		// second instance field
		one8,                      // field name
		[]byte{byte(core.Object)}, // type
		// 381
		createSubRecordHeader(core.HprofGcInstanceDumpType),
		// 383
		one8, // object id
		one4, // stack trace serial number
		one8, // class object id
		one4, // number of bytes that follow
		one1, // single byte
		// 408
		createSubRecordHeader(core.HprofGcObjArrayDumpType),
		// 409
		one8, // array object id
		one4, // stack trace serial number
		one4, // number of elements
		one8, // array class id
		one8, // single element
		// 441
		createSubRecordHeader(core.HprofGcPrimArrayDumpType),
		// 442
		one8,                       // array object id
		one4,                       // stack trace serial number
		one4,                       // number of elements
		[]byte{byte(core.Boolean)}, // element type
		one1,                       // single element
		// 460
		// another sub record begins here
		createRecordHeader(core.HprofHeapDumpSegmentTag, 0),
		// 469
		createSubRecordHeader(core.HprofGcRootStickyClassType),
		// 470
		two8, // object id
		// 478
		createRecordHeader(core.HprofHeapDumpEndTag, 0),
		// 487
	)
)

var reader = createReader(testHeapDump)

func TestReader_GetHprofUtf8(t *testing.T) {
	got, err := reader.GetHprofUtf8(1)
	want := core.HprofUtf8{
		Identifier: 1,
		Characters: "JAVA",
	}
	if err != nil {
		t.Errorf("GetHprofUtf8 error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetHprofUtf8 = %v, want %v", got, want)
	}
}

func TestReader_GetHprofLoadClassByClassSerialNumber(t *testing.T) {
	got, err := reader.GetHprofLoadClassByClassSerialNumer(1)
	want := core.HprofLoadClass{
		ClassSerialNumber:      1,
		ClassObjectId:          1,
		StackTraceSerialNumber: 1,
		ClassNameId:            1,
	}
	if err != nil {
		t.Errorf("GetHprofLoadClassByClassSerialNumber error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetHprofLoadClassByClassSerialNumber = %v, want %v", got, want)
	}
}

func TestReader_GetHprofLoadClassByClassObjectId(t *testing.T) {
	got, err := reader.GetHprofLoadClassByClassObjectId(1)
	want := core.HprofLoadClass{
		ClassSerialNumber:      1,
		ClassObjectId:          1,
		StackTraceSerialNumber: 1,
		ClassNameId:            1,
	}
	if err != nil {
		t.Errorf("GetHprofLoadClassByClassObjectId error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetHprofLoadClassByClassObjectId = %v, want %v", got, want)
	}
}

func TestReader_GetHprofFrame(t *testing.T) {
	got, err := reader.GetHprofFrame(1)
	want := core.HprofFrame{
		StackFrameId:      1,
		MethodNameId:      1,
		MethodSignatureId: 1,
		SourceFileNameId:  1,
		ClassSerialNumber: 1,
		LineNumber:        1,
	}
	if err != nil {
		t.Errorf("GetHprofFrame error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetHprofFrame = %v, want %v", got, want)
	}
}

func TestReader_GetHprofTrace(t *testing.T) {
	got, err := reader.GetHprofTrace(1)
	want := core.HprofTrace{
		StackTraceSerialNumber: 1,
		ThreadSerialNumber:     1,
		NumberOfFrames:         3,
		StackFrameIds:          []core.Identifier{1, 1, 1},
	}
	if err != nil {
		t.Errorf("GetHprofTrace error = %v", err)

	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetHprofTrace = %v, want %v", got, want)
	}
}

func TestReader_GetHprofGcClassDump(t *testing.T) {
	got, err := reader.GetHprofGcClassDump(1)
	want := core.HprofGcClassDump{
		ClassObjectId:            1,
		StackTraceSerialNumber:   1,
		SuperclassObjectId:       1,
		ClassloaderObjectId:      1,
		SignersObjectId:          1,
		ProtectionDomainObjectId: 1,
		InstanceSize:             1,
		SizeOfConstantPool:       2,
		ConstantPoolRecords: []core.HprofGcClassDumpConstantPoolRecord{
			{ConstantPoolIndex: 1, Ty: core.Boolean, Value: core.JavaValue{Type: core.Boolean, Value: false}},
			{ConstantPoolIndex: 1, Ty: core.Boolean, Value: core.JavaValue{Type: core.Boolean, Value: true}},
		},
		NumberOfStaticFields: 2,
		StaticFieldRecords: []core.HprofGcClassDumpStaticFieldsRecord{
			{StaticFieldName: 1, Ty: core.Boolean, Value: core.JavaValue{Type: core.Boolean, Value: false}},
			{StaticFieldName: 1, Ty: core.Boolean, Value: core.JavaValue{Type: core.Boolean, Value: true}},
		},
		NumberOfInstanceFields: 2,
		InstanceFieldRecords: []core.HprofGcClassDumpInstanceFieldsRecord{
			{InstanceFieldName: 1, Ty: core.Object},
			{InstanceFieldName: 1, Ty: core.Object},
		},
	}
	if err != nil {
		t.Errorf("GetHprofGcClassDump error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetHprofGcClassDump = %v, want %v", got, want)
	}
}

func TestReader_GetHprofGcInstanceDump(t *testing.T) {
	got, err := reader.GetHprofGcInstanceDump(1)
	want := core.HprofGcClassDumpInstanceDumpHeader{
		ObjectId:                1,
		StackTraceSerialNumber:  1,
		ClassObjectId:           1,
		NumberOfBytesThatFollow: 1,
	}
	if err != nil {
		t.Errorf("GetHprofGcInstanceDump() error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetHprofGcInstanceDump() = %v, want %v", got, want)
	}
}

func TestReader_GetHprofGcObjArray(t *testing.T) {
	got, err := reader.GetHprofGcObjArray(1)
	want := core.HprofGcObjArrayDumpHeader{
		ArrayObjectId:          1,
		StackTraceSerialNumber: 1,
		NumberOfElements:       1,
		ArrayClassId:           1,
	}
	if err != nil {
		t.Errorf("GetHprofGcObjArray() error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetHprofGcObjArray() = %v, want %v", got, want)
	}
}

func TestReader_GetHprofGcPrimArray(t *testing.T) {
	got, err := reader.GetHprofGcPrimArray(1)
	want := core.HprofGcPrimArrayDumpHeader{
		ArrayObjectId:          1,
		StackTraceSerialNumber: 1,
		NumberOfElements:       1,
		ElementType:            core.Boolean,
	}
	if err != nil {
		t.Errorf("GetHprofGcPrimArray() error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetHprofGcPrimArray() = %v, want %v", got, want)
	}
}

func TestReader_GetBytesFromCurrent(t *testing.T) {
	_, err := reader.heapDump.Seek(0, io.SeekStart)
	if err != nil {
		t.Errorf("error seeking to 0: %v", err)
	}
	got, err := reader.GetBytesFromCurrent(2)
	want := []byte{0x4a, 0x41}
	if err != nil {
		t.Errorf("GetBytesFromCurrent() error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetBytesFromCurrent() = %v, want %v", got, want)
	}
}

func TestReader_GetMeta(t *testing.T) {
	expectedMeta := storage.MetaStorage{
		Counters: storage.Counters{
			InstancesCount: map[core.Identifier]int{
				1: 1,
			},
			PrimArraysCount: map[core.JavaType]int{
				core.Boolean: 1,
			},
			PrimArrayElementsCount: map[core.JavaType]int{
				core.Boolean: 1,
			},
			ObjArraysCount: map[core.Identifier]int{
				1: 1,
			},
			ObjArrayElementsCount: map[core.Identifier]int{
				1: 1,
			},
		},
	}

	if !reflect.DeepEqual(reader.MetaStorage, expectedMeta) {
		t.Errorf("MetaStorage = %+v, expected %+v", reader.MetaStorage, expectedMeta)
	}
}

func concat(head []byte, tail ...[]byte) []byte {
	res := head[:]
	for _, b := range tail {
		res = append(res, b...)
	}
	return res
}

func createReader(in []byte) *ParsedAccessor {
	heapDump := bytes.NewReader(in)
	smallWriter := storage.NewSmallRecordsWriteStorage()
	instanceDumpWriteVolume := storage.NewRamWriteVolume()
	objArrayDumpWriteVolume := storage.NewRamWriteVolume()
	primArrayDumpWriteVolume := storage.NewRamWriteVolume()
	bigWriter := storage.NewBigRecordsWriteStorage(
		instanceDumpWriteVolume, objArrayDumpWriteVolume, primArrayDumpWriteVolume)
	metaWriter := storage.NewMetaWriteStorage()
	parser := NewParser(heapDump, smallWriter, bigWriter, metaWriter)

	if err := parser.ParseHeapDump(); err != nil {
		panic(fmt.Errorf("error indexing sample input: %v", err))
	}
	bigReader, err := storage.NewBigRecordsReadStorage(
		storage.NewRamReadVolume(instanceDumpWriteVolume.Bytes()), instanceDumpWriteVolume.Len(),
		storage.NewRamReadVolume(objArrayDumpWriteVolume.Bytes()), objArrayDumpWriteVolume.Len(),
		storage.NewRamReadVolume(primArrayDumpWriteVolume.Bytes()), primArrayDumpWriteVolume.Len(),
	)
	if err != nil {
		panic(fmt.Errorf("error creating bigreader: %v", err))
	}
	smallReader := storage.NewSmallRecordsReadStorage()
	metaReader := storage.NewMetaReadStorage()
	smallBuf := bytes.NewBuffer(nil)
	metaBuf := bytes.NewBuffer(nil)
	if err := smallWriter.SerializeTo(smallBuf); err != nil {
		panic(fmt.Errorf("error serializing smallwriter: %v", err))
	}
	if err := metaWriter.SerializeTo(metaBuf); err != nil {
		panic(fmt.Errorf("error serializing meta writer: %v", err))
	}
	if err := smallReader.RestoreFrom(smallBuf); err != nil {
		panic(fmt.Errorf("error creating smallreader: %v", err))
	}
	if err := metaReader.RestoreFrom(metaBuf); err != nil {
		panic(fmt.Errorf("error creating meta reader: %v", err))
	}
	reader := NewParsedAccessor(heapDump, bigReader, smallReader, metaReader)
	return reader
}

// record header length = 9
func createRecordHeader(tag core.Tag, remaining uint32) []byte {
	start := []byte{
		byte(tag),              // record tag
		0x00, 0x00, 0x00, 0x00, // time, always zero
	}
	end := make([]byte, 4) // remaining bytes in record
	binary.BigEndian.PutUint32(end, remaining)
	return concat(start, end)
}

func createSubRecordHeader(ty core.SubRecordType) []byte {
	return []byte{byte(ty)}
}

var (
	one1   = []byte{0x01}
	one2   = []byte{0x00, 0x01}
	one4   = []byte{0x00, 0x00, 0x00, 0x01}
	one8   = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}
	three4 = []byte{0x00, 0x00, 0x00, 0x03}
	two2   = []byte{0x00, 0x02}
	two8   = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02}
)
