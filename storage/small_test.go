package storage

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/danielleontiev/neojhat/core"
)

func TestSmallStorage(t *testing.T) {
	var idSize uint32 = 1
	hprofUtf8 := core.HprofUtf8{Identifier: 1, Characters: "hello"}
	hprofLoadClass := core.HprofLoadClass{
		ClassSerialNumber:      1,
		ClassObjectId:          1,
		StackTraceSerialNumber: 1,
		ClassNameId:            1,
	}
	hprofFrame := core.HprofFrame{
		StackFrameId:      1,
		MethodNameId:      1,
		MethodSignatureId: 1,
		SourceFileNameId:  1,
		ClassSerialNumber: 1,
		LineNumber:        1,
	}
	hprofTrace := core.HprofTrace{
		StackTraceSerialNumber: 1,
		ThreadSerialNumber:     1,
		NumberOfFrames:         1,
		StackFrameIds:          nil,
	}
	hprofGcRootJniGlobal := core.HprofGcRootJniGlobal{
		ObjectId:       1,
		JniGlobalRefId: 1,
	}
	hprofGcRootJniLocal := core.HprofGcRootJniLocal{
		ObjectId:                1,
		ThreadSerialNumber:      1,
		FrameNumberInStackTrace: 1,
	}
	hprofGcRootJavaFrame := core.HprofGcRootJavaFrame{
		ObjectId:                1,
		ThreadSerialNumber:      1,
		FrameNumberInStackTrace: 1,
	}
	hprofGcRootStickyClass := core.HprofGcRootStickyClass{ObjectId: 1}
	hprofGcRootThreadObj := core.HprofGcRootThreadObj{
		ThreadObjectId:           1,
		ThreadSequenceNumber:     1,
		StackTraceSequenceNumber: 1,
	}
	hprofGcClassDump := core.HprofGcClassDump{
		ClassObjectId:            1,
		StackTraceSerialNumber:   1,
		SuperclassObjectId:       1,
		ClassloaderObjectId:      1,
		SignersObjectId:          1,
		ProtectionDomainObjectId: 1,
		InstanceSize:             1,
		SizeOfConstantPool:       1,
		ConstantPoolRecords:      nil,
		NumberOfStaticFields:     1,
		StaticFieldRecords:       nil,
		NumberOfInstanceFields:   1,
		InstanceFieldRecords:     nil,
	}

	writeStorage := NewSmallRecordsWriteStorage()

	writeStorage.PutIdSize(idSize)
	writeStorage.PutHprofUtf8(hprofUtf8)
	writeStorage.PutHprofLoadClass(hprofLoadClass)
	writeStorage.PutHprofFrame(hprofFrame)
	writeStorage.PutHprofTrace(hprofTrace)
	writeStorage.PutHprofGcRootJniGlobal(hprofGcRootJniGlobal)
	writeStorage.PutHprofGcRootJniLocal(hprofGcRootJniLocal)
	writeStorage.PutHprofGcRootJavaFrame(hprofGcRootJavaFrame)
	writeStorage.PutHprofGcRootStickyClass(hprofGcRootStickyClass)
	writeStorage.PutHprofGcRootThreadObj(hprofGcRootThreadObj)
	writeStorage.PutHprofGcClassDump(hprofGcClassDump)

	buffer := bytes.NewBuffer(nil)

	if err := writeStorage.SerializeTo(buffer); err != nil {
		t.Errorf("Serialize() err = %v", err)
	}

	readStorage := new(SmallRecordsReadStorage)

	if err := readStorage.RestoreFrom(buffer); err != nil {
		t.Errorf("Restore() err = %v", err)
	}

	gotHprofUtf8, err := readStorage.GetHprofUtf8(1)
	if err != nil {
		t.Errorf("GetHprofUtf8() error = %v", err)
	}
	if !reflect.DeepEqual(gotHprofUtf8, hprofUtf8) {
		t.Errorf("GetHprofUtf8() = %v, expected %v", gotHprofUtf8, hprofUtf8)
	}
	if _, err := readStorage.GetHprofUtf8(2); err == nil {
		t.Errorf("GetHprofUtf8 err = nil")
	}

	gotHprofLoadClassByObjectId, err := readStorage.GetHprofLoadClassByClassObjectId(1)
	if err != nil {
		t.Errorf("GetHprofLoadClassByClassObjectId() error = %v", err)
	}
	if !reflect.DeepEqual(gotHprofLoadClassByObjectId, hprofLoadClass) {
		t.Errorf("GetHprofLoadClassByClassObjectId() = %v, expected %v", gotHprofLoadClassByObjectId, hprofLoadClass)
	}
	if _, err := readStorage.GetHprofLoadClassByClassObjectId(2); err == nil {
		t.Errorf("GetHprofLoadClassByClassObjectId err = nil")
	}

	gotHprofLoadClassByClassSerialNumber, err := readStorage.GetHprofLoadClassByClassSerialNumer(1)
	if err != nil {
		t.Errorf("GetHprofLoadClassByClassSerialNumer() error = %v", err)
	}
	if !reflect.DeepEqual(gotHprofLoadClassByClassSerialNumber, hprofLoadClass) {
		t.Errorf("GetHprofLoadClassByClassSerialNumer() = %v, expected %v", gotHprofLoadClassByClassSerialNumber, hprofLoadClass)
	}
	if _, err := readStorage.GetHprofLoadClassByClassSerialNumer(2); err == nil {
		t.Errorf("GetHprofLoadClassByClassSerialNumer err = nil")
	}

	gotHprofFrame, err := readStorage.GetHprofFrame(1)
	if err != nil {
		t.Errorf("GetHprofFrame() error = %v", err)
	}
	if !reflect.DeepEqual(gotHprofFrame, hprofFrame) {
		t.Errorf("GetHprofFrame() = %v, expected %v", gotHprofFrame, hprofFrame)
	}
	if _, err := readStorage.GetHprofFrame(2); err == nil {
		t.Errorf("GetHprofFrame err = nil")
	}

	gotHprofTrace, err := readStorage.GetHprofTrace(1)
	if err != nil {
		t.Errorf("GetHprofTrace() error = %v", err)
	}
	if !reflect.DeepEqual(gotHprofTrace, hprofTrace) {
		t.Errorf("GetHprofTrace() = %v, expected %v", gotHprofTrace, hprofTrace)
	}
	if _, err := readStorage.GetHprofTrace(2); err == nil {
		t.Errorf("GetHprofTrace err = nil")
	}

	gotHprofGcRootJniGlobals := readStorage.ListHprofGcRootJniGlobal()
	if !reflect.DeepEqual(gotHprofGcRootJniGlobals, []core.HprofGcRootJniGlobal{hprofGcRootJniGlobal}) {
		t.Errorf("ListHprofGcRootJniGlobal() = %v, expected [%v]", gotHprofGcRootJniGlobals, hprofGcRootJniGlobal)
	}

	gotHprofGcRootJniLocals := readStorage.ListHprofGcRootJniLocal()
	if !reflect.DeepEqual(gotHprofGcRootJniLocals, []core.HprofGcRootJniLocal{hprofGcRootJniLocal}) {
		t.Errorf("ListHprofGcRootJniLocal() = %v, expected [%v]", gotHprofGcRootJniLocals, hprofGcRootJniLocal)
	}

	gotHprofGcRootJavaFrames := readStorage.ListHprofGcRootJavaFrame()
	if !reflect.DeepEqual(gotHprofGcRootJavaFrames, []core.HprofGcRootJavaFrame{hprofGcRootJavaFrame}) {
		t.Errorf("ListHprofGcRootJavaFrame() = %v, expected [%v]", gotHprofGcRootJavaFrames, hprofGcRootJavaFrame)
	}

	gotHprofGcRootStickyClasses := readStorage.ListHprofGcRootStickyClass()
	if !reflect.DeepEqual(gotHprofGcRootStickyClasses, []core.HprofGcRootStickyClass{hprofGcRootStickyClass}) {
		t.Errorf("ListHprofGcRootStickyClass() = %v, expected [%v]", gotHprofGcRootStickyClasses, hprofGcRootStickyClass)
	}

	gotHprofGcRootThreadObjs := readStorage.ListHprofGcRootThreadObj()
	if !reflect.DeepEqual(gotHprofGcRootThreadObjs, []core.HprofGcRootThreadObj{hprofGcRootThreadObj}) {
		t.Errorf("ListHprofGcRootThreadObj() = %v, expected [%v]", gotHprofGcRootThreadObjs, hprofGcRootThreadObj)
	}

	gotHprofGcClassDump, err := readStorage.GetHprofGcClassDump(1)
	if err != nil {
		t.Errorf("GetHprofGcClassDump() error = %v", err)
	}
	if !reflect.DeepEqual(gotHprofGcClassDump, hprofGcClassDump) {
		t.Errorf("GetHprofGcClassDump() = %v, expected %v", gotHprofGcClassDump, hprofGcClassDump)
	}
	if _, err := readStorage.GetHprofGcClassDump(2); err == nil {
		t.Errorf("GetHprofGcClassDump err = nil")
	}
}
