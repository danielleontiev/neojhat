package dump

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/danielleontiev/neojhat/core"
	"github.com/danielleontiev/neojhat/storage"
)

func TestParser_ParseHeapDump(t *testing.T) {
	heapDump := bytes.NewReader(testHeapDump)
	smallWriter := storage.NewSmallRecordsWriteStorage()
	instanceDumpWriteVolume := storage.NewRamWriteVolume()
	objArrayDumpWriteVolume := storage.NewRamWriteVolume()
	primArrayDumpWriteVolume := storage.NewRamWriteVolume()
	bigWriter := storage.NewBigRecordsWriteStorage(
		instanceDumpWriteVolume, objArrayDumpWriteVolume, primArrayDumpWriteVolume)
	metaWriter := storage.NewMetaWriteStorage()
	creator := NewParser(heapDump, smallWriter, bigWriter, metaWriter)

	if err := creator.ParseHeapDump(); err != nil {
		t.Errorf("ParseHeapDump() error = %v", err)
	}
	bigReader, err := storage.NewBigRecordsReadStorage(
		storage.NewRamReadVolume(instanceDumpWriteVolume.Bytes()), instanceDumpWriteVolume.Len(),
		storage.NewRamReadVolume(objArrayDumpWriteVolume.Bytes()), objArrayDumpWriteVolume.Len(),
		storage.NewRamReadVolume(primArrayDumpWriteVolume.Bytes()), primArrayDumpWriteVolume.Len(),
	)
	smallReader := storage.NewSmallRecordsReadStorage()
	metaReader := storage.NewMetaReadStorage()
	smallBuf := bytes.NewBuffer(nil)
	metaBuf := bytes.NewBuffer(nil)
	if err := smallWriter.SerializeTo(smallBuf); err != nil {
		t.Errorf("error serializing smallwriter: %v", err)
	}
	if err := metaWriter.SerializeTo(metaBuf); err != nil {
		t.Errorf("error serializing meta writer: %v", err)
	}
	if err := smallReader.RestoreFrom(smallBuf); err != nil {
		t.Errorf("error creating smallreader: %v", err)
	}
	if err := metaReader.RestoreFrom(metaBuf); err != nil {
		t.Errorf("error creating meta reader: %v", err)
	}

	newHeapDump := bytes.NewReader(testHeapDump)
	newParser := core.NewRecordParser(newHeapDump, 8)
	// utf8
	utf8, err := smallReader.GetHprofUtf8(1)
	if err != nil {
		t.Errorf("cannot get: %v", err)
	}
	expectedUtf8 := core.HprofUtf8{Identifier: 1, Characters: "JAVA"}
	if !reflect.DeepEqual(utf8, expectedUtf8) {
		t.Errorf("GetHprofUtf8 = %v, want %v", utf8, expectedUtf8)
	}
	// loadClass
	loadClass, err := smallReader.GetHprofLoadClassByClassSerialNumer(1)
	if err != nil {
		t.Errorf("cannot get: %v", err)
	}
	expectedLoadClass := core.HprofLoadClass{
		ClassSerialNumber: 1, ClassObjectId: 1,
		StackTraceSerialNumber: 1, ClassNameId: 1,
	}
	if !reflect.DeepEqual(loadClass, expectedLoadClass) {
		t.Errorf("GetHprofLoadClassByClassSerialNumer = %v, want %v", loadClass, expectedLoadClass)
	}
	// trace
	trace, err := smallReader.GetHprofTrace(1)
	if err != nil {
		t.Errorf("cannot get: %v", err)
	}
	expectedTrace := core.HprofTrace{
		StackTraceSerialNumber: 1, ThreadSerialNumber: 1, NumberOfFrames: 3,
		StackFrameIds: []core.Identifier{1, 1, 1}}
	if !reflect.DeepEqual(trace, expectedTrace) {
		t.Errorf("GetHprofTrace = %v, want %v", trace, expectedTrace)
	}
	// frame
	frame, err := smallReader.GetHprofFrame(1)
	if err != nil {
		t.Errorf("cannot get: %v", err)
	}
	expectedFrame := core.HprofFrame{StackFrameId: 1, MethodNameId: 1, MethodSignatureId: 1, SourceFileNameId: 1, ClassSerialNumber: 1, LineNumber: 1}
	if !reflect.DeepEqual(frame, expectedFrame) {
		t.Errorf("GetHprofFrame = %v, want %v", frame, expectedFrame)
	}
	// thread obj
	threadObjects := smallReader.ListHprofGcRootThreadObj()
	expectedThreadObjects := []core.HprofGcRootThreadObj{
		{ThreadObjectId: 1, ThreadSequenceNumber: 1, StackTraceSequenceNumber: 1},
	}
	if !reflect.DeepEqual(threadObjects, expectedThreadObjects) {
		t.Errorf("ListHprofGcRootThreadObj = %v, want %v", threadObjects, expectedThreadObjects)
	}
	// jni local
	jniLocals := smallReader.ListHprofGcRootJniLocal()
	expectedJniLocals := []core.HprofGcRootJniLocal{
		{ObjectId: 1, ThreadSerialNumber: 1, FrameNumberInStackTrace: 1},
	}
	if !reflect.DeepEqual(jniLocals, expectedJniLocals) {
		t.Errorf("ListHprofGcRootJniLocal = %v, want %v", jniLocals, expectedJniLocals)
	}
	// jni global
	jniGlobals := smallReader.ListHprofGcRootJniGlobal()
	expectedJniGlobals := []core.HprofGcRootJniGlobal{
		{ObjectId: 1, JniGlobalRefId: 1},
	}
	if !reflect.DeepEqual(jniGlobals, expectedJniGlobals) {
		t.Errorf("ListHprofGcRootJniGlobal = %v, want %v", jniGlobals, expectedJniGlobals)
	}
	// java frame
	javaFrames := smallReader.ListHprofGcRootJavaFrame()
	expectedJavaFrames := []core.HprofGcRootJavaFrame{
		{ObjectId: 1, ThreadSerialNumber: 1, FrameNumberInStackTrace: 1},
	}
	if !reflect.DeepEqual(javaFrames, expectedJavaFrames) {
		t.Errorf("ListHprofGcRootJavaFrame = %v, want %v", javaFrames, expectedJavaFrames)
	}
	// sticky class
	stickyClasses := smallReader.ListHprofGcRootStickyClass()
	expectedStickyClasses := []core.HprofGcRootStickyClass{
		{ObjectId: 1},
		{ObjectId: 2},
	}
	if !reflect.DeepEqual(stickyClasses, expectedStickyClasses) {
		t.Errorf("ListHprofGcRootStickyClass = %v, want %v", stickyClasses, expectedStickyClasses)
	}
	// class dump
	classDump, err := smallReader.GetHprofGcClassDump(1)
	if err != nil {
		t.Errorf("cannot get: %v", err)
	}
	expectedClassDump := core.HprofGcClassDump{
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
	if !reflect.DeepEqual(classDump, expectedClassDump) {
		t.Errorf("GetHprofGcClassDump = %v, want %v", classDump, expectedClassDump)
	}
	// instance dump
	instanceDumpOffset, err := bigReader.HprofGcInstanceDumpGetOffset(1)
	if err != nil {
		t.Errorf("HprofGcInstanceDumpGetOffset() error = %v", err)
	}
	if instanceDumpOffset != 383 {
		t.Errorf("instanceDumpOffset = %v, want %v", instanceDumpOffset, 383)
	}
	seekTo(instanceDumpOffset, newHeapDump, t)
	instanceDump, err := newParser.ParseHprofGcClassDumpInstanceDumpHeader()
	if err != nil {
		t.Errorf("cannot parse: %v", err)
	}
	expectedInstanceDump := core.HprofGcClassDumpInstanceDumpHeader{
		ObjectId: 1, StackTraceSerialNumber: 1, ClassObjectId: 1, NumberOfBytesThatFollow: 1}
	if !reflect.DeepEqual(instanceDump, expectedInstanceDump) {
		t.Errorf("ParseHprofGcClassDumpInstanceDumpHeader = %v, want %v", instanceDump, expectedInstanceDump)
	}
	// obj array
	objArrayOffset, err := bigReader.HprofGcObjArrayDumpGetOffset(1)
	if err != nil {
		t.Errorf("HprofGcObjArrayDumpGetOffset() error = %v", err)
	}
	if objArrayOffset != 409 {
		t.Errorf("objArrayOffset = %v, want %v", objArrayOffset, 409)
	}
	seekTo(objArrayOffset, newHeapDump, t)
	objArray, err := newParser.ParseHprofGcObjArrayDumpHeader()
	if err != nil {
		t.Errorf("cannot parse: %v", err)
	}
	expectedObjArray := core.HprofGcObjArrayDumpHeader{
		ArrayObjectId: 1, StackTraceSerialNumber: 1, NumberOfElements: 1, ArrayClassId: 1}
	if !reflect.DeepEqual(objArray, expectedObjArray) {
		t.Errorf("ParseHprofGcObjArrayDumpHeader = %v, want %v", objArray, expectedObjArray)
	}
	// prim array
	primArrayOffset, err := bigReader.HprofGcPrimArrayDumpGetOffset(1)
	if err != nil {
		t.Errorf("HprofGcPrimArrayDumpGetOffset() error = %v", err)
	}
	if primArrayOffset != 442 {
		t.Errorf("primArrayOffset = %v, want %v", primArrayOffset, 442)
	}
	seekTo(primArrayOffset, newHeapDump, t)
	primArray, err := newParser.ParseHprofGcPrimArrayDumpHeader()
	if err != nil {
		t.Errorf("cannot parse: %v", err)
	}
	expectedPrimArray := core.HprofGcPrimArrayDumpHeader{
		ArrayObjectId: 1, StackTraceSerialNumber: 1, NumberOfElements: 1, ElementType: core.Boolean}
	if !reflect.DeepEqual(primArray, expectedPrimArray) {
		t.Errorf("ParseHprofGcPrimArrayDumpHeader = %v, want %v", primArray, expectedPrimArray)
	}
	// final position
	pos := creator.GetPosition()
	if (pos != 487) || pos != len(testHeapDump) {
		t.Errorf("wrong position = %v, want %v", pos, 487)
	}
}

func seekTo(offset int, readSeeker io.ReadSeeker, t *testing.T) {
	_, err := readSeeker.Seek(int64(offset), io.SeekStart)
	if err != nil {
		t.Errorf("cannot seek to %v: %v", offset, err)
	}
}
