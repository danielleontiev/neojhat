package storage

import (
	"fmt"
	"io"
	"strings"

	"github.com/danielleontiev/neojhat/core"
)

type BigRecordsWriteStorage struct {
	instanceDumpPersistent  *IndexRecordsWriteStorage
	objArrayDumpPersistent  *IndexRecordsWriteStorage
	primArrayDumpPersistent *IndexRecordsWriteStorage
}

func NewBigRecordsWriteStorage(
	instanceDumpPersistent io.WriteCloser,
	objArrayDumpPersistent io.WriteCloser,
	primArrayDumpPersistent io.WriteCloser,
) *BigRecordsWriteStorage {
	instanceStorage := NewIndexRecordsWriteStorage(instanceDumpPersistent, DefaultBatchSize)
	objArrayStorage := NewIndexRecordsWriteStorage(objArrayDumpPersistent, DefaultBatchSize)
	primArrayStorage := NewIndexRecordsWriteStorage(primArrayDumpPersistent, DefaultBatchSize)
	return &BigRecordsWriteStorage{
		instanceDumpPersistent:  instanceStorage,
		objArrayDumpPersistent:  objArrayStorage,
		primArrayDumpPersistent: primArrayStorage,
	}
}

type BigRecordsReadStorage struct {
	instanceDumpPersistent  *IndexRecordsReadStorage
	objArrayDumpPersistent  *IndexRecordsReadStorage
	primArrayDumpPersistent *IndexRecordsReadStorage
}

func NewBigRecordsReadStorage(
	instanceDumpPersistent IndexRecordsReaderAtCloser,
	instanceDumpPersistentSize int,
	objArrayDumpPersistent IndexRecordsReaderAtCloser,
	objArrayDumpPersistentSize int,
	primArrayDumpPersistent IndexRecordsReaderAtCloser,
	primArrayDumpPersistentSize int,
) (*BigRecordsReadStorage, error) {
	instanceStorage, err := NewIndexRecordsReadStorage(instanceDumpPersistent, instanceDumpPersistentSize)
	if err != nil {
		return nil, fmt.Errorf("Cannot create bytestorage: %w", err)
	}
	objArrayStorage, err := NewIndexRecordsReadStorage(objArrayDumpPersistent, objArrayDumpPersistentSize)
	if err != nil {
		return nil, fmt.Errorf("Cannot create bytestorage: %w", err)
	}
	primArrayStorage, err := NewIndexRecordsReadStorage(primArrayDumpPersistent, primArrayDumpPersistentSize)
	if err != nil {
		return nil, fmt.Errorf("Cannot create bytestorage: %w", err)
	}
	return &BigRecordsReadStorage{
		instanceDumpPersistent:  instanceStorage,
		objArrayDumpPersistent:  objArrayStorage,
		primArrayDumpPersistent: primArrayStorage,
	}, nil
}

func (r *BigRecordsReadStorage) HprofGcInstanceDumpGetOffset(objectId core.Identifier) (int, error) {
	offset, err := r.instanceDumpPersistent.Get(uint64(objectId))
	return int(offset), err
}

func (r *BigRecordsReadStorage) HprofGcObjArrayDumpGetOffset(arrayObjectId core.Identifier) (int, error) {
	offset, err := r.objArrayDumpPersistent.Get(uint64(arrayObjectId))
	return int(offset), err
}

func (r *BigRecordsReadStorage) HprofGcPrimArrayDumpGetOffset(arrayObjectId core.Identifier) (int, error) {
	offset, err := r.primArrayDumpPersistent.Get(uint64(arrayObjectId))
	return int(offset), err
}

func (r *BigRecordsReadStorage) Close() error {
	err1 := r.instanceDumpPersistent.Close()
	err2 := r.objArrayDumpPersistent.Close()
	err3 := r.primArrayDumpPersistent.Close()
	err := combineErrors("Cannot close BigRecordsReadStorage", err1, err2, err3)
	return err
}

func (w *BigRecordsWriteStorage) HprofGcInstanceDumpPutOffset(objectId core.Identifier, offset int) error {
	err := w.instanceDumpPersistent.Put(uint64(objectId), uint64(offset))
	return err
}

func (w *BigRecordsWriteStorage) HprofGcObjArrayDumpPutOffset(arrayObjectId core.Identifier, offset int) error {
	err := w.objArrayDumpPersistent.Put(uint64(arrayObjectId), uint64(offset))
	return err
}

func (w *BigRecordsWriteStorage) HprofGcPrimArrayDumpPutOffset(arrayObjectId core.Identifier, offset int) error {
	err := w.primArrayDumpPersistent.Put(uint64(arrayObjectId), uint64(offset))
	return err
}

func (w *BigRecordsWriteStorage) Close() error {
	err1 := w.instanceDumpPersistent.Close()
	err2 := w.objArrayDumpPersistent.Close()
	err3 := w.primArrayDumpPersistent.Close()
	err := combineErrors("Cannot close BigRecordsWriteStorage", err1, err2, err3)
	return err
}

func combineErrors(label string, errors ...error) error {
	var messages []string
	for _, err := range errors {
		if err != nil {
			messages = append(messages, err.Error())
		}
	}
	if len(messages) != 0 {
		return fmt.Errorf("%s: %s", label, strings.Join(messages, ", "))
	}
	return nil
}
