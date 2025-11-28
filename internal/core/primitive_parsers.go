package core

import (
	"encoding/binary"
	"fmt"
	"io"
)

func NewPrimitiveParser(heapDump io.Reader, idSize uint32) *PrimitiveParser {
	return &PrimitiveParser{heapDump: heapDump, identifierSize: idSize}
}

// PrimitiveParser is a parser that allows to read
// low-level primitives like u1, u2, u4, etc. numbers
// and other types that presented in heap dump.
type PrimitiveParser struct {
	heapDump       io.Reader
	identifierSize uint32
}

// ParseUint8 reads 8-bit unsigned (u1) number.
func (parser *PrimitiveParser) ParseUint8() (uint8, error) {
	result := make([]byte, 1)
	_, err := io.ReadFull(parser.heapDump, result)
	if err != nil {
		return 0, fmt.Errorf("error in ParseUint8: %w", err)
	}
	return result[0], nil
}

// ParseUint16 reads 16-bit unsigned (u2) number.
func (parser *PrimitiveParser) ParseUint16() (uint16, error) {
	var result uint16
	err := binary.Read(parser.heapDump, binary.BigEndian, &result)
	if err != nil {
		return 0, fmt.Errorf("error in ParseUint16: %w", err)
	}
	return result, nil
}

// ParseUint32 reads 32-bit unsigned (u4) number.
func (parser *PrimitiveParser) ParseUint32() (uint32, error) {
	var result uint32
	err := binary.Read(parser.heapDump, binary.BigEndian, &result)
	if err != nil {
		return 0, fmt.Errorf("error in ParseUint32: %w", err)
	}
	return result, nil
}

// ParseInt32 reads 32-bit signed (i4) number.
func (parser *PrimitiveParser) ParseInt32() (int32, error) {
	var result int32
	err := binary.Read(parser.heapDump, binary.BigEndian, &result)
	if err != nil {
		return 0, fmt.Errorf("error in ParseInt32: %w", err)
	}
	return result, nil
}

// ParseIdentifier reads JVM identifier that can be 32 or 64 bit
// value. It depends on the machine on which JVM runs.
func (parser *PrimitiveParser) ParseIdentifier() (Identifier, error) {
	var identifier uint64
	switch parser.identifierSize {
	case 8:
		err := binary.Read(parser.heapDump, binary.BigEndian, &identifier)
		if err != nil {
			return 0, fmt.Errorf("error in ParseIdentifier(8): %w", err)
		}
	case 4:
		var shortIdentifier uint32
		err := binary.Read(parser.heapDump, binary.BigEndian, &shortIdentifier)
		if err != nil {
			return 0, fmt.Errorf("error in ParseIdentifier(4): %w", err)
		}
		identifier = uint64(shortIdentifier)
	default:
		return 0, fmt.Errorf("unsupported identifier size, only [4, 8] are supported")
	}
	return Identifier(identifier), nil
}

// ParseLineNumber reads LineNumber which is 32-bit signed under the hood.
func (parser *PrimitiveParser) ParseLineNumber() (LineNumber, error) {
	result, err := parser.ParseInt32()
	if err != nil {
		return 0, fmt.Errorf("error in ParseLineNumber: %w", err)
	}
	return LineNumber(result), nil
}

// ParseJavaType reads JavaType which is 8-bit unsigned under the hood.
func (parser *PrimitiveParser) ParseJavaType() (JavaType, error) {
	result, err := parser.ParseUint8()
	if err != nil {
		return 0, fmt.Errorf("error in ParseJavaType: %w", err)
	}
	return JavaType(result), nil
}

// ParseJavaValue reads JavaValue of given JavaType and stores it
// in struct along with the type for future use.
func (parser *PrimitiveParser) ParseJavaValue(ty JavaType) (JavaValue, error) {
	switch ty {
	case Object:
		result, err := parser.ParseIdentifier()
		if err != nil {
			return JavaValue{}, fmt.Errorf("error in ParseJavaValue(Object): %w", err)
		}
		return JavaValue{Type: ty, Value: result}, nil
	case Boolean:
		var result bool
		err := binary.Read(parser.heapDump, binary.BigEndian, &result)
		if err != nil {
			return JavaValue{}, fmt.Errorf("error in ParseJavaValue(Boolean): %w", err)
		}
		return JavaValue{Type: ty, Value: result}, nil
	case Char:
		buf := make([]byte, 2)
		_, err := io.ReadFull(parser.heapDump, buf)
		if err != nil {
			return JavaValue{}, fmt.Errorf("error in ParseJavaValue(Char): %w", err)
		}
		return JavaValue{Type: ty, Value: string(buf)}, nil
	case Float:
		var result float32
		err := binary.Read(parser.heapDump, binary.BigEndian, &result)
		if err != nil {
			return JavaValue{}, fmt.Errorf("error in ParseJavaValue(Float): %w", err)
		}
		return JavaValue{Type: ty, Value: result}, nil
	case Double:
		var result float64
		err := binary.Read(parser.heapDump, binary.BigEndian, &result)
		if err != nil {
			return JavaValue{}, fmt.Errorf("error in ParseJavaValue(Double): %w", err)
		}
		return JavaValue{Type: ty, Value: result}, nil
	case Byte:
		var result int8
		err := binary.Read(parser.heapDump, binary.BigEndian, &result)
		if err != nil {
			return JavaValue{}, fmt.Errorf("error in ParseJavaValue(Byte): %w", err)
		}
		return JavaValue{Type: ty, Value: result}, nil
	case Short:
		var result int16
		err := binary.Read(parser.heapDump, binary.BigEndian, &result)
		if err != nil {
			return JavaValue{}, fmt.Errorf("error in ParseJavaValue(Short): %w", err)
		}
		return JavaValue{Type: ty, Value: result}, nil
	case Int:
		result, err := parser.ParseInt32()
		if err != nil {
			return JavaValue{}, fmt.Errorf("error in ParseJavaValue(Int): %w", err)
		}
		return JavaValue{Type: ty, Value: result}, nil
	case Long:
		var result int64
		err := binary.Read(parser.heapDump, binary.BigEndian, &result)
		if err != nil {
			return JavaValue{}, fmt.Errorf("error in ParseJavaValue(Long): %w", err)
		}
		return JavaValue{Type: ty, Value: result}, nil
	}
	return JavaValue{}, fmt.Errorf("unexpected Java type: %v", byte(ty))
}

// ParseByteSeq reads n bytes from underlying reader
func (parser *PrimitiveParser) ParseByteSeq(n int) ([]byte, error) {
	buf := make([]byte, n)
	_, err := io.ReadFull(parser.heapDump, buf)
	if err != nil {
		return nil, fmt.Errorf("error in ParseByteSeq: %w", err)
	}
	return buf, nil
}
