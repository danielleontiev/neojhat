package storage

import (
	"encoding/gob"
	"fmt"
	"io"

	"github.com/danielleontiev/neojhat/internal/core"
)

func init() {
	gob.Register(core.Identifier(0))
	gob.Register(core.JavaType(0))
}

type MetaWriteStorage struct {
	MetaStorage
}

func NewMetaWriteStorage() *MetaWriteStorage {
	underlyingMetaStorage := MetaStorage{
		Counters: Counters{
			InstancesCount:         make(map[core.Identifier]int),
			PrimArraysCount:        make(map[core.JavaType]int),
			PrimArrayElementsCount: make(map[core.JavaType]int),
			ObjArraysCount:         make(map[core.Identifier]int),
			ObjArrayElementsCount:  make(map[core.Identifier]int),
		},
	}
	return &MetaWriteStorage{underlyingMetaStorage}
}

func (s *MetaWriteStorage) SerializeTo(destination io.Writer) error {
	encoder := gob.NewEncoder(destination)
	if err := encoder.Encode(s.MetaStorage); err != nil {
		return fmt.Errorf("cannot serialize: %w", err)
	}
	return nil
}

func (s *MetaWriteStorage) AddInstance(obj any) {
	switch o := obj.(type) {
	case core.HprofGcClassDumpInstanceDumpHeader:
		s.MetaStorage.Counters.InstancesCount[o.ClassObjectId]++
	case core.HprofGcObjArrayDumpHeader:
		s.MetaStorage.Counters.ObjArraysCount[o.ArrayClassId]++
		s.MetaStorage.Counters.ObjArrayElementsCount[o.ArrayClassId] += int(o.NumberOfElements)
	case core.HprofGcPrimArrayDumpHeader:
		s.MetaStorage.Counters.PrimArraysCount[o.ElementType]++
		s.MetaStorage.Counters.PrimArrayElementsCount[o.ElementType] += int(o.NumberOfElements)
	}
}

type MetaReadStorage struct {
	MetaStorage
}

func NewMetaReadStorage() *MetaReadStorage {
	return new(MetaReadStorage)
}

func (s *MetaReadStorage) RestoreFrom(source io.Reader) error {
	var underlyingMetaStorage MetaStorage
	decoder := gob.NewDecoder(source)
	if err := decoder.Decode(&underlyingMetaStorage); err != nil {
		return fmt.Errorf("cannot deserialize: %w", err)
	}
	s.MetaStorage = underlyingMetaStorage
	return nil
}

type MetaStorage struct {
	Counters Counters
}

type Counters struct {
	InstancesCount         map[core.Identifier]int
	PrimArraysCount        map[core.JavaType]int
	PrimArrayElementsCount map[core.JavaType]int
	ObjArraysCount         map[core.Identifier]int
	ObjArrayElementsCount  map[core.Identifier]int
}
