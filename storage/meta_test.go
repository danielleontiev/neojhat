package storage

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/danielleontiev/neojhat/core"
)

func TestMetaStorage(t *testing.T) {
	writeStorage := NewMetaWriteStorage()

	obj1 := core.HprofGcClassDumpInstanceDumpHeader{
		ClassObjectId: 1,
	}
	obj2 := core.HprofGcClassDumpInstanceDumpHeader{
		ClassObjectId: 2,
	}
	objArr := core.HprofGcObjArrayDumpHeader{
		NumberOfElements: 10,
	}
	primArr := core.HprofGcPrimArrayDumpHeader{
		ElementType:      core.Int,
		NumberOfElements: 5,
	}

	writeStorage.AddInstance(obj1)
	writeStorage.AddInstance(obj2)
	writeStorage.AddInstance(obj2) // two times
	writeStorage.AddInstance(objArr)
	writeStorage.AddInstance(primArr)

	buffer := bytes.NewBuffer(nil)

	if err := writeStorage.SerializeTo(buffer); err != nil {
		t.Errorf("Serialize() err = %v", err)
	}

	readStorage := NewMetaReadStorage()

	if err := readStorage.RestoreFrom(buffer); err != nil {
		t.Errorf("Restore() err = %v", err)
	}

	expected := MetaStorage{
		Counters: Counters{
			InstancesCount: map[core.Identifier]int{
				1: 1,
				2: 2,
			},
			PrimArraysCount: map[core.JavaType]int{
				core.Int: 1,
			},
			PrimArrayElementsCount: map[core.JavaType]int{
				core.Int: 5,
			},
			ObjArraysCount: map[core.Identifier]int{
				0: 1,
			},
			ObjArrayElementsCount: map[core.Identifier]int{
				0: 10,
			},
		},
	}
	if !reflect.DeepEqual(readStorage.MetaStorage, expected) {
		t.Errorf("Instances = %+v, expected %+v", readStorage.MetaStorage, expected)
	}
}
