// dump is the package responsible parsing the whole .hprof file and
// accessing parsed information saved on disk. For big records (instances,
// object arrays and primitive arrays) parser uses indexing techinque -
// it does not store parsed records, instead, it saves position of the
// record to index file that will be used later by accessor to find position
// of the record in the dump and parse it on demand. Tipically, big
// records can occupy >90% of heap dump size and if heap dump is big then
// all records simply do not fit into RAM. Moreover, heap is dumped to
// .hprof by traversing all objects in JVM and they are put to the file
// as they are presented in the heap. So, accessing arbitrary object by its
// identifier requires searching and it could be optimized using index.
// Other records are small and usually fit into RAM easily. For such records
// parser puts them into in-memory storage and after parsing the storage is
// serialized to the file. Accessor simply restores the storage from file when
// needed.
package dump

import (
	"bufio"
	"fmt"
	"io"

	"github.com/danielleontiev/neojhat/internal/core"
	"github.com/danielleontiev/neojhat/internal/storage"
)

// Parser traverses .hprof file and saves parsed information to storages.
type Parser struct {
	pos                      int
	heapDump                 io.Reader
	smallRecordsWriteStorage *storage.SmallRecordsWriteStorage
	bigRecordsWriteStorage   *storage.BigRecordsWriteStorage
	metaWriteStorage         *storage.MetaWriteStorage
}

func NewParser(
	heapDump io.Reader,
	smallRecordsWriteStorage *storage.SmallRecordsWriteStorage,
	bigRecordsWriteStorage *storage.BigRecordsWriteStorage,
	metaWriteStorage *storage.MetaWriteStorage,
) *Parser {
	return &Parser{
		pos:                      0,
		heapDump:                 heapDump,
		smallRecordsWriteStorage: smallRecordsWriteStorage,
		bigRecordsWriteStorage:   bigRecordsWriteStorage,
		metaWriteStorage:         metaWriteStorage,
	}
}

// GetPosition returns the current position while
// parsing. Mainly used to provide interactive
// progress bar since parsing large heap dumps
// can be time consuming.
func (creator *Parser) GetPosition() int {
	return creator.pos
}

// ParseHeapDump parses heap dump to storages.
// Can be used with arbitrary io.Reader.
func (parser *Parser) ParseHeapDump() error {
	bufferedHeapDump := bufio.NewReader(parser.heapDump)
	fileHeader, err := core.ParseFileHeader(bufferedHeapDump)
	if err != nil {
		return fmt.Errorf("error parsing .hprof header: %w", err)
	}
	if fileHeader.Header != core.ValidProfileVersion {
		return fmt.Errorf("unsupported profile version %v, only %v supported", fileHeader.Header, core.ValidProfileVersion)
	}
	parser.smallRecordsWriteStorage.PutIdSize(fileHeader.IdentifierSize)
	parser.smallRecordsWriteStorage.PutTimestamp(fileHeader.Timestamp)

	size := core.NewSizeInfo(fileHeader.IdentifierSize)
	recordParser := core.NewRecordParser(bufferedHeapDump, fileHeader.IdentifierSize)

	parser.pos = 31

	defer parser.bigRecordsWriteStorage.Close()

	for {
		header, err := recordParser.ParseRecordHeader()
		parser.pos += 9
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("error parsing record header: %w", err)
		}
		switch header.Tag {
		case core.HprofUtf8Tag:
			record, err := recordParser.ParseHprofUtf8(header.Remaining)
			if err != nil {
				return fmt.Errorf("error parsing HprofUtf8: %w", err)
			}
			parser.smallRecordsWriteStorage.PutHprofUtf8(record)
			parser.pos += int(header.Remaining)
		case core.HprofLoadClassTag:
			record, err := recordParser.ParseHprofLoadClass()
			if err != nil {
				return fmt.Errorf("error parsing HprofLoadClass: %w", err)
			}
			parser.smallRecordsWriteStorage.PutHprofLoadClass(record)
			parser.pos += int(header.Remaining)
		case core.HprofFrameTag:
			record, err := recordParser.ParseHprofFrame()
			if err != nil {
				return fmt.Errorf("error parsing HprofFrame: %w", err)
			}
			parser.smallRecordsWriteStorage.PutHprofFrame(record)
			parser.pos += int(header.Remaining)
		case core.HprofTraceTag:
			record, err := recordParser.ParseHprofTrace()
			if err != nil {
				return fmt.Errorf("error parsing HprofTrace: %w", err)
			}
			parser.smallRecordsWriteStorage.PutHprofTrace(record)
			parser.pos += int(header.Remaining)
		case core.HprofHeapDumpSegmentTag:
		loop:
			for {
				subRecordHeader, err := recordParser.ParseSubRecordHeader()
				if err != nil {
					return fmt.Errorf("error parsing sub-record type: %w", err)
				}
				parser.pos++
				switch subRecordHeader.SubRecordType {
				case core.HprofGcRootJniGlobalType:
					record, err := recordParser.ParseHprofGcRootJniGlobal()
					if err != nil {
						return fmt.Errorf("error parsing HprofGcRootJniGlobal: %w", err)
					}
					parser.smallRecordsWriteStorage.PutHprofGcRootJniGlobal(record)
					parser.pos += size.Of(record)
				case core.HprofGcRootJniLocalType:
					record, err := recordParser.ParseHprofGcRootJniLocal()
					if err != nil {
						return fmt.Errorf("error parsing HprofGcRootJniLocal: %w", err)
					}
					parser.smallRecordsWriteStorage.PutHprofGcRootJniLocal(record)
					parser.pos += size.Of(record)
				case core.HprofGcRootJavaFrameType:
					record, err := recordParser.ParseHprofGcRootJavaFrame()
					if err != nil {
						return fmt.Errorf("error parsing HprofGcRootJavaFrame: %w", err)
					}
					parser.smallRecordsWriteStorage.PutHprofGcRootJavaFrame(record)
					parser.pos += size.Of(record)
				case core.HprofGcRootStickyClassType:
					record, err := recordParser.ParseHprofGcRootStickyClass()
					if err != nil {
						return fmt.Errorf("error parsing HprofGcRootStickyClass: %w", err)
					}
					parser.smallRecordsWriteStorage.PutHprofGcRootStickyClass(record)
					parser.pos += size.Of(record)
				case core.HprofGcRootThreadObjType:
					record, err := recordParser.ParseHprofGcRootThreadObj()
					if err != nil {
						return fmt.Errorf("error parsing HprofGcRootThreadObj: %w", err)
					}
					parser.smallRecordsWriteStorage.PutHprofGcRootThreadObj(record)
					parser.pos += size.Of(record)
				case core.HprofGcClassDumpType:
					record, err := recordParser.ParseHprofGcClassDump()
					if err != nil {
						return fmt.Errorf("error parsing HprofGcClassDump: %w", err)
					}
					parser.smallRecordsWriteStorage.PutHprofGcClassDump(record)
					parser.pos += size.Of(record)
				case core.HprofGcInstanceDumpType:
					record, err := recordParser.ParseHprofGcClassDumpInstanceDumpHeader()
					if err != nil {
						return fmt.Errorf("error parsing HprofGcClassDumpInstanceDump: %w", err)
					}
					if err := parser.bigRecordsWriteStorage.HprofGcInstanceDumpPutOffset(record.ObjectId, parser.pos); err != nil {
						return fmt.Errorf("indexing error: HprofGcInstanceDumpPutOffset: %w", err)
					}
					fullSize, recordsSize := size.OfObject(record)
					parser.pos += fullSize
					parser.metaWriteStorage.AddInstance(record)
					if err := skip(recordsSize, bufferedHeapDump); err != nil {
						return fmt.Errorf("error discarding records of HprofGcClassDumpInstanceDump: %w", err)
					}
				case core.HprofGcObjArrayDumpType:
					record, err := recordParser.ParseHprofGcObjArrayDumpHeader()
					if err != nil {
						return fmt.Errorf("error parsing HprofGcObjArrayDump: %w", err)
					}
					if err := parser.bigRecordsWriteStorage.HprofGcObjArrayDumpPutOffset(record.ArrayObjectId, parser.pos); err != nil {
						return fmt.Errorf("indexing error: HprofGcObjArrayDumpPutOffset: %w", err)
					}
					fullSize, recordsSize := size.OfObject(record)
					parser.pos += fullSize
					parser.metaWriteStorage.AddInstance(record)
					if err := skip(recordsSize, bufferedHeapDump); err != nil {
						return fmt.Errorf("error discarding records of HprofGcObjArrayDump: %w", err)
					}
				case core.HprofGcPrimArrayDumpType:
					record, err := recordParser.ParseHprofGcPrimArrayDumpHeader()
					if err != nil {
						return fmt.Errorf("error parsing HprofGcPrimArrayDump: %w", err)
					}
					if err := parser.bigRecordsWriteStorage.HprofGcPrimArrayDumpPutOffset(record.ArrayObjectId, parser.pos); err != nil {
						return fmt.Errorf("indexing error: HprofGcPrimArrayDumpPutOffset: %w", err)
					}
					fullSize, recordsSize := size.OfObject(record)
					parser.pos += fullSize
					parser.metaWriteStorage.AddInstance(record)
					if err := skip(recordsSize, bufferedHeapDump); err != nil {
						return fmt.Errorf("error discarding records of HprofGcPrimArrayDump: %w", err)
					}
				case core.HprofHeapDumpEndSubRecord:
					if err := unreadByte(bufferedHeapDump); err != nil {
						return fmt.Errorf("error unreading byte at HprofHeapDumpEndSubRecord: %w", err)
					}
					parser.pos--
					break loop
				case core.HprofHeapDumpSegmentSubRecord:
					if err := unreadByte(bufferedHeapDump); err != nil {
						return fmt.Errorf("error unreading byte at HprofHeapDumpSegmentSubRecord: %w", err)
					}
					parser.pos--
					break loop
				}
			}
		case core.HprofHeapDumpEndTag:
			return nil
		default:
			return fmt.Errorf("unknown tag: %v", header.Tag)
		}
	}
}

// skip calls underlying bufio.Reader.Discard
func skip(n int, bufReader *bufio.Reader) error {
	_, err := bufReader.Discard(n)
	if err != nil {
		return fmt.Errorf("cannot skip bytes: %w", err)
	}
	return nil
}

// unreadByte calls underlying bufio.Reader.UnreadByte
func unreadByte(bufReader *bufio.Reader) error {
	if err := bufReader.UnreadByte(); err != nil {
		return fmt.Errorf("cannot unread byte: %w", err)
	}
	return nil
}
