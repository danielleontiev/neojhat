/*
	core is a package for parsing JVM .hprof files.

	Models are based on the following JDK source:
	https://github.com/openjdk/jdk/blob/c1e39faaa99ee62ff626ffec9f978ed0f8ffaca1/src/hotspot/share/services/heapDumper.cpp

	It contains basic utility primitives that should be used to parse
	the given .hprof dump in the form of io.Reader. Functions in
	this packages are not aware of the position in the input, they are
	simple primitives to parse structures that .hprof could contain. The whole
	structure of .hprof file is written in code comment at the link above.
*/
package core

import (
	"time"
)

type FileHeader struct {
	Header         string
	IdentifierSize uint32
	Timestamp      time.Time
}

type RecordHeader struct {
	Tag       Tag
	Remaining uint32
}

type HprofUtf8 struct {
	Identifier Identifier
	Characters string
}

type HprofLoadClass struct {
	ClassSerialNumber      uint32
	ClassObjectId          Identifier
	StackTraceSerialNumber uint32
	ClassNameId            Identifier
}

type HprofFrame struct {
	StackFrameId      Identifier
	MethodNameId      Identifier
	MethodSignatureId Identifier
	SourceFileNameId  Identifier
	ClassSerialNumber uint32
	LineNumber        LineNumber
}

type HprofTrace struct {
	StackTraceSerialNumber uint32
	ThreadSerialNumber     uint32
	NumberOfFrames         uint32
	StackFrameIds          []Identifier
}

type SubRecordHeader struct {
	SubRecordType SubRecordType
}

type HprofGcRootThreadObj struct {
	ThreadObjectId           Identifier
	ThreadSequenceNumber     uint32
	StackTraceSequenceNumber uint32
}

type HprofGcRootJniGlobal struct {
	ObjectId       Identifier
	JniGlobalRefId Identifier
}

type HprofGcRootJniLocal struct {
	ObjectId                Identifier
	ThreadSerialNumber      uint32
	FrameNumberInStackTrace uint32
}

type HprofGcRootJavaFrame struct {
	ObjectId                Identifier
	ThreadSerialNumber      uint32
	FrameNumberInStackTrace uint32
}

type HprofGcRootStickyClass struct {
	ObjectId Identifier
}

type HprofGcClassDump struct {
	ClassObjectId            Identifier
	StackTraceSerialNumber   uint32
	SuperclassObjectId       Identifier
	ClassloaderObjectId      Identifier
	SignersObjectId          Identifier
	ProtectionDomainObjectId Identifier
	InstanceSize             int32
	SizeOfConstantPool       uint16
	ConstantPoolRecords      []HprofGcClassDumpConstantPoolRecord
	NumberOfStaticFields     uint16
	StaticFieldRecords       []HprofGcClassDumpStaticFieldsRecord
	NumberOfInstanceFields   uint16
	InstanceFieldRecords     []HprofGcClassDumpInstanceFieldsRecord
}

type HprofGcClassDumpConstantPoolRecord struct {
	ConstantPoolIndex uint16
	Ty                JavaType
	Value             JavaValue
}

type HprofGcClassDumpStaticFieldsRecord struct {
	StaticFieldName Identifier
	Ty              JavaType
	Value           JavaValue
}

type HprofGcClassDumpInstanceFieldsRecord struct {
	InstanceFieldName Identifier
	Ty                JavaType
}

type HprofGcClassDumpInstanceDumpHeader struct {
	ObjectId                Identifier
	StackTraceSerialNumber  uint32
	ClassObjectId           Identifier
	NumberOfBytesThatFollow uint32
}

type HprofGcObjArrayDumpHeader struct {
	ArrayObjectId          Identifier
	StackTraceSerialNumber uint32
	NumberOfElements       uint32
	ArrayClassId           Identifier
}

type HprofGcPrimArrayDumpHeader struct {
	ArrayObjectId          Identifier
	StackTraceSerialNumber uint32
	NumberOfElements       uint32
	ElementType            JavaType
}
