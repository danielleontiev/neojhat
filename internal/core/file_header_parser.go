package core

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

const (
	ValidProfileVersion = "JAVA PROFILE 1.0.2"
)

// ParseFileHeader reads the header of .hprof file.
func ParseFileHeader(heapDump io.Reader) (FileHeader, error) {
	str := make([]byte, len(ValidProfileVersion)+1)
	_, err := io.ReadFull(heapDump, str)
	if err != nil {
		return FileHeader{}, fmt.Errorf("error in ParseFileHeader: %w", err)
	}
	header := string(str[:len(str)-1])
	if header != ValidProfileVersion {
		return FileHeader{}, fmt.Errorf("unsupported profile %v, only %v is supported", header, ValidProfileVersion)
	}
	identifiersSize, err := parserUint32(heapDump)
	if err != nil {
		return FileHeader{}, fmt.Errorf("error in ParseFileHeader: %w", err)
	}
	high, err := parserUint32(heapDump)
	if err != nil {
		return FileHeader{}, fmt.Errorf("error in ParseFileHeader: %w", err)
	}
	low, err := parserUint32(heapDump)
	if err != nil {
		return FileHeader{}, fmt.Errorf("error in ParseFileHeader: %w", err)
	}
	milli := int64(high)
	milli <<= 32
	milli += int64(low)
	timestamp := time.Unix(0, 0).Add(time.Duration(milli) * time.Millisecond)
	return FileHeader{
		Header:         header,
		IdentifierSize: identifiersSize,
		Timestamp:      timestamp,
	}, nil
}

// this function duplicates the functionality
// of PrimitiveParser.ParseUint32()
// because we need to parse identifier size
// before we can use PrimitiveParser
func parserUint32(reader io.Reader) (uint32, error) {
	var result uint32
	err := binary.Read(reader, binary.BigEndian, &result)
	if err != nil {
		return 0, fmt.Errorf("error in parserUint32: %w", err)
	}
	return result, nil
}
