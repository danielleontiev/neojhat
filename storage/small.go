package storage

import (
	"encoding/gob"
	"fmt"
	"io"
	"time"

	"github.com/danielleontiev/neojhat/core"
)

func init() {
	gob.Register(core.Identifier(0))
}

type underlyingStorage struct {
	IdSize                 uint32
	Timestamp              time.Time
	HprofUtf8              map[core.Identifier]core.HprofUtf8
	HprofLoadClass         []core.HprofLoadClass
	HprofFrame             map[core.Identifier]core.HprofFrame
	HprofTrace             map[uint32]core.HprofTrace
	HprofGcRootJniGlobal   []core.HprofGcRootJniGlobal
	HprofGcRootJniLocal    []core.HprofGcRootJniLocal
	HprofGcRootJavaFrame   []core.HprofGcRootJavaFrame
	HprofGcRootStickyClass []core.HprofGcRootStickyClass
	HprofGcRootThreadObj   []core.HprofGcRootThreadObj
	HprofGcClassDump       map[core.Identifier]core.HprofGcClassDump
}

type SmallRecordsWriteStorage struct {
	underlyingStorage
}

func NewSmallRecordsWriteStorage() *SmallRecordsWriteStorage {
	underlyingStorage := underlyingStorage{
		HprofUtf8:        make(map[core.Identifier]core.HprofUtf8),
		HprofFrame:       make(map[core.Identifier]core.HprofFrame),
		HprofTrace:       make(map[uint32]core.HprofTrace),
		HprofGcClassDump: make(map[core.Identifier]core.HprofGcClassDump),
	}
	return &SmallRecordsWriteStorage{underlyingStorage}
}

func (s *SmallRecordsWriteStorage) SerializeTo(destination io.Writer) error {
	encoder := gob.NewEncoder(destination)
	if err := encoder.Encode(s.underlyingStorage); err != nil {
		return fmt.Errorf("cannot serialize: %w", err)
	}
	return nil
}

func (s *SmallRecordsReadStorage) RestoreFrom(source io.Reader) error {
	var underlyingStorage underlyingStorage
	decoder := gob.NewDecoder(source)
	if err := decoder.Decode(&underlyingStorage); err != nil {
		return fmt.Errorf("cannot deserialize: %w", err)
	}
	s.underlyingStorage = underlyingStorage
	return nil
}

type SmallRecordsReadStorage struct {
	underlyingStorage
}

func NewSmallRecordsReadStorage() *SmallRecordsReadStorage {
	return new(SmallRecordsReadStorage)
}

func (s *SmallRecordsWriteStorage) PutIdSize(idSize uint32) {
	s.IdSize = idSize
}

func (s *SmallRecordsWriteStorage) PutTimestamp(timestamp time.Time) {
	s.Timestamp = timestamp
}

func (s *SmallRecordsWriteStorage) PutHprofUtf8(record core.HprofUtf8) {
	s.HprofUtf8[record.Identifier] = record
}

func (s *SmallRecordsWriteStorage) PutHprofLoadClass(record core.HprofLoadClass) {
	s.HprofLoadClass = append(s.HprofLoadClass, record)
}

func (s *SmallRecordsWriteStorage) PutHprofFrame(record core.HprofFrame) {
	s.HprofFrame[record.StackFrameId] = record
}

func (s *SmallRecordsWriteStorage) PutHprofTrace(record core.HprofTrace) {
	s.HprofTrace[record.ThreadSerialNumber] = record
}

func (s *SmallRecordsWriteStorage) PutHprofGcRootJniGlobal(record core.HprofGcRootJniGlobal) {
	s.HprofGcRootJniGlobal = append(s.HprofGcRootJniGlobal, record)
}

func (s *SmallRecordsWriteStorage) PutHprofGcRootJniLocal(record core.HprofGcRootJniLocal) {
	s.HprofGcRootJniLocal = append(s.HprofGcRootJniLocal, record)
}

func (s *SmallRecordsWriteStorage) PutHprofGcRootJavaFrame(record core.HprofGcRootJavaFrame) {
	s.HprofGcRootJavaFrame = append(s.HprofGcRootJavaFrame, record)
}

func (s *SmallRecordsWriteStorage) PutHprofGcRootStickyClass(record core.HprofGcRootStickyClass) {
	s.HprofGcRootStickyClass = append(s.HprofGcRootStickyClass, record)
}

func (s *SmallRecordsWriteStorage) PutHprofGcRootThreadObj(record core.HprofGcRootThreadObj) {
	s.HprofGcRootThreadObj = append(s.HprofGcRootThreadObj, record)
}

func (s *SmallRecordsWriteStorage) PutHprofGcClassDump(record core.HprofGcClassDump) {
	s.HprofGcClassDump[record.ClassObjectId] = record
}

func (s *SmallRecordsReadStorage) GetHprofUtf8(nameId core.Identifier) (core.HprofUtf8, error) {
	res, ok := s.HprofUtf8[nameId]
	if !ok {
		return core.HprofUtf8{}, fmt.Errorf("Cannot find HprofUtf8 record with nameId = %v", nameId)
	}
	return res, nil
}

func (s *SmallRecordsReadStorage) GetHprofLoadClassByClassObjectId(classObjectId core.Identifier) (core.HprofLoadClass, error) {
	for _, rec := range s.HprofLoadClass {
		if rec.ClassObjectId == classObjectId {
			return rec, nil
		}
	}
	return core.HprofLoadClass{}, fmt.Errorf("Cannot find HprofLoadClass record with classObjectId = %v", classObjectId)
}

func (s *SmallRecordsReadStorage) GetHprofLoadClassByClassSerialNumer(classSerialNumber uint32) (core.HprofLoadClass, error) {
	for _, rec := range s.HprofLoadClass {
		if rec.ClassSerialNumber == classSerialNumber {
			return rec, nil
		}
	}
	return core.HprofLoadClass{}, fmt.Errorf("Cannot find HprofLoadClass record with classSerialNumber = %v", classSerialNumber)
}

func (s *SmallRecordsReadStorage) ListHprofLoadClass() []core.HprofLoadClass {
	return s.HprofLoadClass
}

func (s *SmallRecordsReadStorage) GetHprofFrame(stackFrameId core.Identifier) (core.HprofFrame, error) {
	res, ok := s.HprofFrame[stackFrameId]
	if !ok {
		return core.HprofFrame{}, fmt.Errorf("Cannot find HprofFrame record with stackFrameId = %v", stackFrameId)
	}
	return res, nil
}

func (s *SmallRecordsReadStorage) GetHprofTrace(threadSerialNumber uint32) (core.HprofTrace, error) {
	res, ok := s.HprofTrace[threadSerialNumber]
	if !ok {
		return core.HprofTrace{}, fmt.Errorf("Cannot find HprofTrace record with threadSerialNumber = %v", threadSerialNumber)
	}
	return res, nil
}

func (s *SmallRecordsReadStorage) ListHprofGcRootJniGlobal() []core.HprofGcRootJniGlobal {
	return s.HprofGcRootJniGlobal
}

func (s *SmallRecordsReadStorage) ListHprofGcRootJniLocal() []core.HprofGcRootJniLocal {
	return s.HprofGcRootJniLocal
}

func (s *SmallRecordsReadStorage) ListHprofGcRootJavaFrame() []core.HprofGcRootJavaFrame {
	return s.HprofGcRootJavaFrame
}

func (s *SmallRecordsReadStorage) ListHprofGcRootStickyClass() []core.HprofGcRootStickyClass {
	return s.HprofGcRootStickyClass
}

func (s *SmallRecordsReadStorage) ListHprofGcRootThreadObj() []core.HprofGcRootThreadObj {
	return s.HprofGcRootThreadObj
}

func (s *SmallRecordsReadStorage) GetHprofGcClassDump(classObjectId core.Identifier) (core.HprofGcClassDump, error) {
	res, ok := s.HprofGcClassDump[classObjectId]
	if !ok {
		return core.HprofGcClassDump{}, fmt.Errorf("Cannot find HprofGcClassDump record with classObjectId = %v", classObjectId)
	}
	return res, nil
}
