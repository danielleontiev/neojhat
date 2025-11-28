// storage package has implementation of storages for small objects and
// index storage for big objects.
//
// Storage for small objects is simple structure with hash maps and lists
// which are filled with parsed records during parsing of heap dump. It's
// in-memory storage and to avoid parsing small objects every time, it has
// serialization machinery that allows to put all contents to the file and
// restore it. Data is serialized to <heap dump name>.db/small-records.bin
//
// Storage for big objects does not try to hold all in RAM. Instead it stores
// index to the index file that later is used to access that records. Index
// files are <heap dump>.db/instance-dump.idx.bin,
// <heap dump>.db/obj-array-dump.idx.bin and <heap dump>.db/prim-array-dump.idx.bin
//
// For effective way of accessing index records binary search is used.
// Index records are key:value pairs of 8+8=16 bytes which is used to store
// offsets of some objects in .hprof file. All objects in .hprof files have
// unique identifiers and correspoiding offsets. F.e. "instance dump" object
// with object id = 1 and offset = 1 could be written to index file and
// the offset could be obtained later and the whole class dump could be parsed.
// Typically, using index file is more effective than linear search in .hprof file.
// Hovewer, it could not be true for machines with HDD disk because this mechanism
// assumes the data is in increasing sorted order by key because Get() is
// binary search. The package could be easily misused - data should be
// written in sorted order or should be sorted after write before read. For most
// sections of .hprof such as instance dumps and array dumps it's true by default -
// data is stored in the dump already sorted. But anyway, the order is tracked
// during parsing to detect anomalies in records order.
package storage

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type BatchSize int

const (
	DefaultBatchSize BatchSize = 1_000_000 // Means the size of 1 batch by default if 16 Mb
)

type IndexRecordsWriteCloser interface {
	io.WriteCloser
}

type IndexRecordsReaderAtCloser interface {
	io.ReaderAt
	io.Closer
}

// IndexRecordsWriteStorage is the struct that tracks underlying file
// and currently processed batch.
type IndexRecordsWriteStorage struct {
	persistentStorage IndexRecordsWriteCloser
	curBatch          *Batch
	batchCount        int
	batchSize         BatchSize
	lastSeenKey       uint64
}

// NewIndexRecordsWriteStorage creates new index file for writing index there. The
// second argument should be used to specify batch size.
func NewIndexRecordsWriteStorage(persistentStorage io.WriteCloser, batchSize BatchSize) *IndexRecordsWriteStorage {
	return &IndexRecordsWriteStorage{
		persistentStorage: persistentStorage,
		batchSize:         batchSize,
	}
}

// IndexRecordsReadStorage uses binary search to read
// stored index file
type IndexRecordsReadStorage struct {
	recordsNumber     int
	persistentStorage IndexRecordsReaderAtCloser
}

// NewIndexRecordsReadStorage opens index file for reading. It assumes
// that the index file is sorted, it should be
// checked before using ByteReader.
func NewIndexRecordsReadStorage(persistentStorage IndexRecordsReaderAtCloser, size int) (*IndexRecordsReadStorage, error) {
	if size%16 != 0 {
		return nil, fmt.Errorf("index storage is corrupted, size = %v, size %% 16 != 0", size)
	}

	return &IndexRecordsReadStorage{
		persistentStorage: persistentStorage,
		recordsNumber:     size / 16,
	}, nil
}

// Batch should be used to group index records for
// writing them more effectiveley.
type Batch struct {
	data []byte
}

// put adds data to the batch
func (b *Batch) put(key uint64, val uint64) {
	keyBuf := make([]byte, 8)
	valBuf := make([]byte, 8)
	binary.BigEndian.PutUint64(keyBuf, key)
	binary.BigEndian.PutUint64(valBuf, val)
	b.data = append(b.data, keyBuf...)
	b.data = append(b.data, valBuf...)
}

// Write dumps the given batch to the index file
func (w *IndexRecordsWriteStorage) Write(batch Batch) error {
	_, err := w.persistentStorage.Write(batch.data)
	if err != nil {
		return fmt.Errorf("error writing batch: %w", err)
	}
	return nil
}

// Put takes one index record and puts it to the current
// active batch. If batch size equals to the limit it
// triggers writing to the disk.
func (w *IndexRecordsWriteStorage) Put(key uint64, val uint64) error {
	if key < w.lastSeenKey {
		return fmt.Errorf("error, keys have gone backward. Key = %v, last seen key = %v", key, w.lastSeenKey)
	}
	w.lastSeenKey = key
	if w.batchCount == int(w.batchSize) {
		w.curBatch.put(key, val)
		err := w.Write(*w.curBatch) // TODO shift to another goroutine to not block the current?
		if err != nil {
			return fmt.Errorf("cannot put current batch and continue: %w", err)
		}
		w.curBatch = nil
		w.batchCount = 0
	} else {
		if w.curBatch == nil {
			w.curBatch = new(Batch)
		}
		w.curBatch.put(key, val)
		w.batchCount++
	}
	return nil
}

// Close should be invoked after parsing
// is over to trigger writing of the active
// batch.
func (w *IndexRecordsWriteStorage) Close() error {
	if w.curBatch != nil {
		err := w.Write(*w.curBatch)
		if err != nil {
			return fmt.Errorf("fail to write last batch: %w", err)
		}
	}
	return w.persistentStorage.Close()
}

// Get is used to lookup the offset of the given key
// in the file. It uses binary search algorithm.
func (r *IndexRecordsReadStorage) Get(key uint64) (uint64, error) {
	left := 0
	right := r.recordsNumber - 1
	record := make([]byte, 16)
	for left <= right {
		mid := left + (right-left)/2
		_, err := r.persistentStorage.ReadAt(record, int64(16*mid))
		if err != nil {
			return 0, fmt.Errorf("cannot read record at offset %v: %w", mid, err)
		}
		keyValue := binary.BigEndian.Uint64(record[:8])
		if keyValue == key {
			return binary.BigEndian.Uint64(record[8:]), nil
		}
		if keyValue < key {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return 0, fmt.Errorf("key %v not found", key)
}

// Close closes the underlying file.
func (db *IndexRecordsReadStorage) Close() error {
	return db.persistentStorage.Close()
}

/* in-memory implementation */

type RamWriteVolume struct {
	*bytes.Buffer
}

type RamReadVolume struct {
	*bytes.Reader
}

func NewRamWriteVolume() *RamWriteVolume {
	buffer := bytes.NewBuffer(nil)
	return &RamWriteVolume{Buffer: buffer}
}

func (w *RamWriteVolume) Close() error {
	return nil
}

func NewRamReadVolume(data []byte) *RamReadVolume {
	reader := bytes.NewReader(data)
	return &RamReadVolume{Reader: reader}
}

func (r *RamReadVolume) Close() error {
	return nil
}
