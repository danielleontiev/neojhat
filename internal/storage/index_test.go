package storage

import (
	"testing"
)

func TestIndexRecordsStorage(t *testing.T) {
	writeVolume := NewRamWriteVolume()
	writer := NewIndexRecordsWriteStorage(writeVolume, 100)

	for i := 0; i < 1000; i++ {
		err := writer.Put(uint64(i), uint64(i))
		if err != nil {
			t.Errorf("error putting bytes: %v", err)
		}
	}

	if err := writer.Close(); err != nil {
		t.Errorf("cannot close byte storage after writing: %v", err)
	}

	reader, err := NewIndexRecordsReadStorage(NewRamReadVolume(writeVolume.Bytes()), writeVolume.Len())
	if err != nil {
		t.Errorf("cannot create byte reader: %v", err)
	}
	defer reader.Close()

	var keyvals []uint64
	// should find records with keys 0 - 999
	for i := 0; i < 1000; i++ {
		keyval := uint64(i)
		keyvals = append(keyvals, keyval)
		res, err := reader.Get(keyval)
		if err != nil {
			t.Errorf("error looking for key %v: %v", i, err)
		}
		if res != uint64(i) {
			t.Errorf("Get(%v) = %v, want %v", i, res, i)
		}
	}

	// should fail
	_, err = reader.Get(1000)
	if err == nil {
		t.Errorf("not found error expected")
	}
}

func Test_GoBackward(t *testing.T) {
	writeVolume := NewRamWriteVolume()
	writer := NewIndexRecordsWriteStorage(writeVolume, 100)

	for i := 0; i < 1000; i++ {
		err := writer.Put(uint64(i), uint64(i))
		if err != nil {
			t.Errorf("error putting bytes: %v", err)
		}
	}

	if err := writer.Put(uint64(1000), uint64(0)); err != nil {
		t.Errorf("putting 1000 = %v, no error expected", err)
	}

	if err := writer.Put(uint64(1000), uint64(0)); err != nil {
		t.Errorf("putting 1000 again = %v, no error expected", err)
	}

	if err := writer.Put(uint64(999), uint64(0)); err == nil {
		t.Errorf("putting 999 have not failed but error was expected")
	}
}
