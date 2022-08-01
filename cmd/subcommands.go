package cmd

import (
	"fmt"
	"os"

	"github.com/danielleontiev/neojhat/dump"
	"github.com/danielleontiev/neojhat/objects"
	"github.com/danielleontiev/neojhat/storage"
	"github.com/danielleontiev/neojhat/summary"
	"github.com/danielleontiev/neojhat/threads"
)

const (
	maxMemory = 8 << 30
)

const (
	storageDirSuffix           = ".db/"
	instanceDumpIndexFileName  = "instance-dump.idx.bin"
	objArrayDumpIndexFileName  = "obj-array-dump.idx.bin"
	primArrayDumpIndexFileName = "prim-array-dump.idx.bin"
	smallRecordsFileName       = "small-records.bin"
	metaFileName               = "meta.bin"
)

func GetThreads(hprofFileName string, noColor, localVars bool) error {
	hprof, err := os.Open(hprofFileName)
	if err != nil {
		return fmt.Errorf("can't open file [%s]: %w", hprofFileName, err)
	}
	defer hprof.Close()

	smallRecordsDumpFile, err := os.Open(hprofFileName + storageDirSuffix + smallRecordsFileName)
	if err != nil {
		return err
	}
	defer smallRecordsDumpFile.Close()

	metaDumpFile, err := os.Open(hprofFileName + storageDirSuffix + metaFileName)
	if err != nil {
		return err
	}
	defer metaDumpFile.Close()

	bigReader, err := createBigReader(hprofFileName)
	if err != nil {
		return err
	}
	smallReader, err := createSmallReader(smallRecordsDumpFile)
	if err != nil {
		return err
	}
	metaReader, err := createMetaReader(metaDumpFile)
	if err != nil {
		return err
	}
	parsedAccessor := dump.NewParsedAccessor(hprof, bigReader, smallReader, metaReader)
	threadDump, err := threads.GetThreadDump(parsedAccessor)
	if err != nil {
		return fmt.Errorf("can't parse thread dump: %w", err)
	}
	if noColor {
		threads.PrettyPrint(threadDump, localVars)
		return nil
	}
	threads.PrettyPrintColor(threadDump, localVars)
	return nil
}

func GetSummary(hprofFileName string, noColor, allProps bool) error {
	hprof, err := os.Open(hprofFileName)
	if err != nil {
		return fmt.Errorf("can't open file [%s]: %w", hprofFileName, err)
	}
	defer hprof.Close()

	smallRecordsDumpFile, err := os.Open(hprofFileName + storageDirSuffix + smallRecordsFileName)
	if err != nil {
		return err
	}
	defer smallRecordsDumpFile.Close()

	metaDumpFile, err := os.Open(hprofFileName + storageDirSuffix + metaFileName)
	if err != nil {
		return err
	}
	defer metaDumpFile.Close()

	bigReader, err := createBigReader(hprofFileName)
	if err != nil {
		return err
	}
	smallReader, err := createSmallReader(smallRecordsDumpFile)
	if err != nil {
		return err
	}
	metaReader, err := createMetaReader(metaDumpFile)
	if err != nil {
		return err
	}
	parsedAccessor := dump.NewParsedAccessor(hprof, bigReader, smallReader, metaReader)
	s, err := summary.GetSummary(parsedAccessor, allProps)
	if err != nil {
		return fmt.Errorf("can't parse summary: %w", err)
	}
	if noColor {
		summary.PrettyPrint(s)
		return nil
	}
	summary.PrettyPrintColor(s)
	return nil
}

func GetObjects(hprofFileName string, noColor bool, sortBy objects.SortBy) error {
	hprof, err := os.Open(hprofFileName)
	if err != nil {
		return fmt.Errorf("can't open file [%s]: %w", hprofFileName, err)
	}
	defer hprof.Close()

	smallRecordsDumpFile, err := os.Open(hprofFileName + storageDirSuffix + smallRecordsFileName)
	if err != nil {
		return err
	}
	defer smallRecordsDumpFile.Close()

	metaDumpFile, err := os.Open(hprofFileName + storageDirSuffix + metaFileName)
	if err != nil {
		return err
	}
	defer metaDumpFile.Close()

	bigReader, err := createBigReader(hprofFileName)
	if err != nil {
		return err
	}
	smallReader, err := createSmallReader(smallRecordsDumpFile)
	if err != nil {
		return err
	}
	metaReader, err := createMetaReader(metaDumpFile)
	if err != nil {
		return err
	}
	parsedAccessor := dump.NewParsedAccessor(hprof, bigReader, smallReader, metaReader)
	obj, err := objects.GetObjects(parsedAccessor, sortBy)
	if err != nil {
		return fmt.Errorf("can't parse objects: %w", err)
	}
	if noColor {
		objects.PrettyPrint(obj)
		return nil
	}
	objects.PrettyPrintColor(obj)
	return nil
}

func ParseHprof(hprofFileName string, nonInteractive bool) error {
	hprof, err := os.Open(hprofFileName)
	if err != nil {
		return fmt.Errorf("can't open file [%s]: %w", hprofFileName, err)
	}
	defer hprof.Close()

	if err = os.Mkdir(hprofFileName+storageDirSuffix, os.ModePerm); err != nil {
		if os.IsExist(err) {
			return nil
		}
		return fmt.Errorf("can't create index: %w", err)
	}

	stat, err := hprof.Stat()
	if err != nil {
		return fmt.Errorf("can't get file stats: %w", err)
	}

	smallWriter := storage.NewSmallRecordsWriteStorage()
	instanceDumpIndexFile, err := os.OpenFile(hprofFileName+storageDirSuffix+instanceDumpIndexFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	objArrayDumpIndexFile, err := os.OpenFile(hprofFileName+storageDirSuffix+objArrayDumpIndexFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	primArrayDumpIndexFile, err := os.OpenFile(hprofFileName+storageDirSuffix+primArrayDumpIndexFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	bigWriter := storage.NewBigRecordsWriteStorage(instanceDumpIndexFile, objArrayDumpIndexFile, primArrayDumpIndexFile)
	metaWriter := storage.NewMetaWriteStorage()
	parser := dump.NewParser(hprof, smallWriter, bigWriter, metaWriter)
	cancel := interactive(progressBar(int(stat.Size()), parser.GetPosition, "Parsing"), nonInteractive)
	if err := parser.ParseHeapDump(); err != nil {
		return fmt.Errorf("can't create index: %w", err)
	}
	cancel()

	smallRecordsFile, err := os.OpenFile(hprofFileName+storageDirSuffix+smallRecordsFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer smallRecordsFile.Close()
	if err := smallWriter.SerializeTo(smallRecordsFile); err != nil {
		return fmt.Errorf("can't close small writer: %w", err)
	}
	metaFile, err := os.OpenFile(hprofFileName+storageDirSuffix+metaFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer metaFile.Close()
	if err := metaWriter.SerializeTo(metaFile); err != nil {
		return fmt.Errorf("can't close meta writer: %w", err)
	}
	return nil
}

func createBigReader(hprofFileName string) (*storage.BigRecordsReadStorage, error) {
	instanceDumpFile, err := os.Open(hprofFileName + storageDirSuffix + instanceDumpIndexFileName)
	if err != nil {
		return nil, err
	}
	instanceDumpFileStat, err := instanceDumpFile.Stat()
	if err != nil {
		return nil, err
	}
	objArrayDumpFile, err := os.Open(hprofFileName + storageDirSuffix + objArrayDumpIndexFileName)
	if err != nil {
		return nil, err
	}
	objArrayDumpFileStat, err := objArrayDumpFile.Stat()
	if err != nil {
		return nil, err
	}
	primArrayDumpFile, err := os.Open(hprofFileName + storageDirSuffix + primArrayDumpIndexFileName)
	if err != nil {
		return nil, err
	}
	primArrayDumpFileStat, err := primArrayDumpFile.Stat()
	if err != nil {
		return nil, err
	}
	bigReader, err := storage.NewBigRecordsReadStorage(
		instanceDumpFile, int(instanceDumpFileStat.Size()),
		objArrayDumpFile, int(objArrayDumpFileStat.Size()),
		primArrayDumpFile, int(primArrayDumpFileStat.Size()),
	)
	if err != nil {
		return bigReader, fmt.Errorf("can't create big reader: %w", err)
	}
	return bigReader, nil
}

func createSmallReader(smallRecordsReadStorageFile *os.File) (*storage.SmallRecordsReadStorage, error) {
	smallReader := storage.NewSmallRecordsReadStorage()
	if err := smallReader.RestoreFrom(smallRecordsReadStorageFile); err != nil {
		return nil, err
	}
	return smallReader, nil
}

func createMetaReader(metaReadStorageFile *os.File) (*storage.MetaReadStorage, error) {
	metaReader := storage.NewMetaReadStorage()
	if err := metaReader.RestoreFrom(metaReadStorageFile); err != nil {
		return nil, err
	}
	return metaReader, nil
}
