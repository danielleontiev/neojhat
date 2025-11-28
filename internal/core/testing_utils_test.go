package core

import (
	"bufio"
	"bytes"
)

var (
	empty = []byte{}
	one1  = []byte{0x01}
	one2  = []byte{0x00, 0x01}
	one4  = []byte{0x00, 0x00, 0x00, 0x01}
	one8  = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}
	zero4 = []byte{0x00, 0x00, 0x00, 0x00}
)

func concat(head []byte, tail ...[]byte) []byte {
	res := head[:]
	for _, b := range tail {
		res = append(res, b...)
	}
	return res
}

type CreateOpts struct {
	idSize uint32
}

func createRecordParser(b []byte, opts CreateOpts) RecordParser {
	primitiveParser := createPrimitiveParser(b, opts)
	return RecordParser{
		primitiveParser: &primitiveParser,
	}
}

func createPrimitiveParser(b []byte, opts CreateOpts) PrimitiveParser {
	return PrimitiveParser{
		heapDump:       bufio.NewReader(bytes.NewReader(b)),
		identifierSize: opts.idSize,
	}
}
