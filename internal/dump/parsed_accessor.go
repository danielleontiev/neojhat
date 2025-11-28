package dump

import (
	"fmt"
	"io"

	"github.com/danielleontiev/neojhat/internal/core"
	"github.com/danielleontiev/neojhat/internal/storage"
)

// ParsedAccessor uses storages to retrieve parsed information. It reads
// offsets of records from index and parses the objects from the position
// obtained from index (for big objects) and restores in-memory storage for
// small objects and provides access to it.
type ParsedAccessor struct {
	heapDump              io.ReadSeeker
	recordParser          *core.RecordParser
	IdentifierSize        uint32
	bigRecordsReadStorage *storage.BigRecordsReadStorage
	*storage.SmallRecordsReadStorage
	*storage.MetaReadStorage
}

func NewParsedAccessor(
	headDump io.ReadSeeker,
	bigRecordsReadStorage *storage.BigRecordsReadStorage,
	smallRecordsReadStorage *storage.SmallRecordsReadStorage,
	metaReadStorage *storage.MetaReadStorage,
) *ParsedAccessor {
	recordParser := core.NewRecordParser(headDump, smallRecordsReadStorage.IdSize)
	return &ParsedAccessor{
		heapDump:                headDump,
		recordParser:            recordParser,
		IdentifierSize:          smallRecordsReadStorage.IdSize,
		bigRecordsReadStorage:   bigRecordsReadStorage,
		SmallRecordsReadStorage: smallRecordsReadStorage,
		MetaReadStorage:         metaReadStorage,
	}
}

func (a *ParsedAccessor) GetHprofGcInstanceDump(objectId core.Identifier) (core.HprofGcClassDumpInstanceDumpHeader, error) {
	offset, err := a.bigRecordsReadStorage.HprofGcInstanceDumpGetOffset(objectId)
	if err != nil {
		return core.HprofGcClassDumpInstanceDumpHeader{}, fmt.Errorf("error getting offset of HprofGcClassDumpInstanceDumpHeader with objectId %v: %w", objectId, err)
	}
	if err := a.seek(offset); err != nil {
		return core.HprofGcClassDumpInstanceDumpHeader{}, err
	}
	res, err := a.recordParser.ParseHprofGcClassDumpInstanceDumpHeader()
	if err != nil {
		return core.HprofGcClassDumpInstanceDumpHeader{}, fmt.Errorf("error reading HprofGcClassDumpInstanceDumpHeader at offset %v: %w", offset, err)
	}
	return res, nil
}

func (a *ParsedAccessor) GetHprofGcObjArray(arrayObjectId core.Identifier) (core.HprofGcObjArrayDumpHeader, error) {
	offset, err := a.bigRecordsReadStorage.HprofGcObjArrayDumpGetOffset(arrayObjectId)
	if err != nil {
		return core.HprofGcObjArrayDumpHeader{}, fmt.Errorf("error getting offset of HprofGcObjArrayDumpHeader with arrayObjectId %v: %w", arrayObjectId, err)
	}
	if err := a.seek(offset); err != nil {
		return core.HprofGcObjArrayDumpHeader{}, err
	}
	res, err := a.recordParser.ParseHprofGcObjArrayDumpHeader()
	if err != nil {
		return core.HprofGcObjArrayDumpHeader{}, fmt.Errorf("error reading HprofGcObjArrayDumpHeader at offset %v: %w", offset, err)
	}
	return res, nil
}

func (a *ParsedAccessor) GetHprofGcPrimArray(arrayObjectId core.Identifier) (core.HprofGcPrimArrayDumpHeader, error) {
	offset, err := a.bigRecordsReadStorage.HprofGcPrimArrayDumpGetOffset(arrayObjectId)
	if err != nil {
		return core.HprofGcPrimArrayDumpHeader{}, fmt.Errorf("error getting offset of HprofGcPrimArrayDumpHeader with arrayObjectId %v: %w", arrayObjectId, err)
	}
	if err := a.seek(offset); err != nil {
		return core.HprofGcPrimArrayDumpHeader{}, err
	}
	res, err := a.recordParser.ParseHprofGcPrimArrayDumpHeader()
	if err != nil {
		return core.HprofGcPrimArrayDumpHeader{}, fmt.Errorf("error reading HprofGcPrimArrayDumpHeader at offset %v: %w", offset, err)
	}
	return res, nil
}

func (a *ParsedAccessor) GetBytesFromCurrent(n int) ([]byte, error) {
	res, err := a.recordParser.ReadBytes(n)
	if err != nil {
		return nil, fmt.Errorf("error in ReadBytesFromCurrent: %w", err)
	}
	return res, nil
}

func (a *ParsedAccessor) seek(offset int) error {
	_, err := a.heapDump.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return fmt.Errorf("cannot goto offset %v: %v", offset, err)
	}
	return nil
}
