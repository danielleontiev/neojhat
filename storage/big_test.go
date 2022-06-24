package storage

import (
	"testing"

	"github.com/danielleontiev/neojhat/core"
)

func Test_HprofGcInstanceDumpGetOffset(t *testing.T) {
	instanceDumpWriteVolume := NewRamWriteVolume()
	objArrayDumpWriteVolume := NewRamWriteVolume()
	primArrayDumpWriteVolume := NewRamWriteVolume()
	writer := NewBigRecordsWriteStorage(instanceDumpWriteVolume, objArrayDumpWriteVolume, primArrayDumpWriteVolume)

	arg := core.Identifier(1)
	want := 1

	if err := writer.HprofGcInstanceDumpPutOffset(arg, want); err != nil {
		t.Errorf("HprofGcInstanceDumpPutOffset() error = %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Errorf("cannot close writer: %v", err)
	}

	reader, err := NewBigRecordsReadStorage(
		NewRamReadVolume(instanceDumpWriteVolume.Bytes()), instanceDumpWriteVolume.Len(),
		NewRamReadVolume(objArrayDumpWriteVolume.Bytes()), objArrayDumpWriteVolume.Len(),
		NewRamReadVolume(primArrayDumpWriteVolume.Bytes()), primArrayDumpWriteVolume.Len(),
	)
	if err != nil {
		t.Errorf("Cannot create Reader: %v", err)
	}
	got, err := reader.HprofGcInstanceDumpGetOffset(arg)
	if err != nil {
		t.Errorf("HprofGcInstanceDumpGetOffset() error = %v", err)
	}
	if got != want {
		t.Errorf("HprofGcInstanceDumpGetOffset() = %v, want %v", got, want)
	}
}

func Test_HprofGcObjArrayDumpGetOffset(t *testing.T) {
	instanceDumpWriteVolume := NewRamWriteVolume()
	objArrayDumpWriteVolume := NewRamWriteVolume()
	primArrayDumpWriteVolume := NewRamWriteVolume()
	writer := NewBigRecordsWriteStorage(instanceDumpWriteVolume, objArrayDumpWriteVolume, primArrayDumpWriteVolume)

	arg := core.Identifier(1)
	want := 1

	if err := writer.HprofGcObjArrayDumpPutOffset(arg, want); err != nil {
		t.Errorf("HprofGcObjArrayDumpPutOffset() error = %v", err)
	}

	if err := writer.Close(); err != nil {
		t.Errorf("cannot close writer: %v", err)
	}

	instanceDumpReadVolume := NewRamReadVolume(instanceDumpWriteVolume.Bytes())
	objArrayDumpReadVolume := NewRamReadVolume(objArrayDumpWriteVolume.Bytes())
	primArrayDumpReadVolume := NewRamReadVolume(primArrayDumpWriteVolume.Bytes())
	reader, err := NewBigRecordsReadStorage(
		instanceDumpReadVolume, instanceDumpWriteVolume.Len(),
		objArrayDumpReadVolume, objArrayDumpWriteVolume.Len(),
		primArrayDumpReadVolume, primArrayDumpWriteVolume.Len(),
	)
	if err != nil {
		t.Errorf("Cannot create Reader: %v", err)
	}
	got, err := reader.HprofGcObjArrayDumpGetOffset(arg)
	if err != nil {
		t.Errorf("HprofGcObjArrayDumpGetOffset() error = %v", err)
	}
	if got != want {
		t.Errorf("HprofGcObjArrayDumpGetOffset() = %v, want %v", got, want)
	}
}

func Test_HprofGcPrimArrayDumpGetOffset(t *testing.T) {
	instanceDumpWriteVolume := NewRamWriteVolume()
	objArrayDumpWriteVolume := NewRamWriteVolume()
	primArrayDumpWriteVolume := NewRamWriteVolume()
	writer := NewBigRecordsWriteStorage(instanceDumpWriteVolume, objArrayDumpWriteVolume, primArrayDumpWriteVolume)

	arg := core.Identifier(1)
	want := 1

	if err := writer.HprofGcPrimArrayDumpPutOffset(arg, want); err != nil {
		t.Errorf("Persistent.HprofGcPrimArrayDumpPutOffset() error = %v", err)
	}

	if err := writer.Close(); err != nil {
		t.Errorf("cannot close writer: %v", err)
	}

	instanceDumpReadVolume := NewRamReadVolume(instanceDumpWriteVolume.Bytes())
	objArrayDumpReadVolume := NewRamReadVolume(objArrayDumpWriteVolume.Bytes())
	primArrayDumpReadVolume := NewRamReadVolume(primArrayDumpWriteVolume.Bytes())
	reader, err := NewBigRecordsReadStorage(
		instanceDumpReadVolume, instanceDumpWriteVolume.Len(),
		objArrayDumpReadVolume, objArrayDumpWriteVolume.Len(),
		primArrayDumpReadVolume, primArrayDumpWriteVolume.Len(),
	)
	if err != nil {
		t.Errorf("Cannot create Reader: %v", err)
	}
	got, err := reader.HprofGcPrimArrayDumpGetOffset(arg)
	if err != nil {
		t.Errorf("Persistent.HprofGcPrimArrayDumpGetOffset() error = %v", err)
	}
	if got != want {
		t.Errorf("Persistent.HprofGcPrimArrayDumpGetOffset() = %v, want %v", got, want)
	}
}
