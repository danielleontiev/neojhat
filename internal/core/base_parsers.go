package core

import (
	"fmt"
	"io"
)

func NewRecordParser(heapDump io.Reader, idSize uint32) *RecordParser {
	return &RecordParser{
		primitiveParser: NewPrimitiveParser(heapDump, idSize),
	}
}

// RecordParser is used to parse HprofRecords from raw heap dump.
type RecordParser struct {
	primitiveParser *PrimitiveParser
}

// ParseRecordHeader reads the header of Hprof records which contain
// information about record's TAG and remaining bytes.
func (parser *RecordParser) ParseRecordHeader() (RecordHeader, error) {
	tag, err := parser.primitiveParser.ParseUint8()
	if err != nil {
		return RecordHeader{}, fmt.Errorf("error in ParseRecordHeader: %w", err)
	}
	// time offset is always zero for some reason
	// https://github.com/openjdk/jdk/blob/c1e39faaa99ee62ff626ffec9f978ed0f8ffaca1/src/hotspot/share/services/heapDumper.cpp#L688
	_, err = parser.primitiveParser.ParseUint32()
	if err != nil {
		return RecordHeader{}, fmt.Errorf("error in ParseRecordHeader: %w", err)
	}
	remaining, err := parser.primitiveParser.ParseUint32()
	if err != nil {
		return RecordHeader{}, fmt.Errorf("error in ParseRecordHeader: %w", err)
	}
	return RecordHeader{
		Tag:       Tag(tag),
		Remaining: remaining,
	}, nil
}

// ParseHprofUtf8 reads HPROF_UTF8 record
func (parser *RecordParser) ParseHprofUtf8(remaining uint32) (HprofUtf8, error) {
	identifier, err := parser.primitiveParser.ParseIdentifier()
	if err != nil {
		return HprofUtf8{}, fmt.Errorf("error in ParseHprofUtf8: %w", err)
	}
	buf := make([]byte, remaining-parser.primitiveParser.identifierSize)
	_, err = io.ReadFull(parser.primitiveParser.heapDump, buf)
	if err != nil {
		return HprofUtf8{}, fmt.Errorf("error in ParseHprofUtf8: %w", err)
	}
	return HprofUtf8{
		Identifier: identifier,
		Characters: string(buf),
	}, nil
}

// ParseHprofLoadClass reads HPROF_LOAD_CLASS record.
func (parser *RecordParser) ParseHprofLoadClass() (HprofLoadClass, error) {
	classSerialNumber, err := parser.primitiveParser.ParseUint32()
	if err != nil {
		return HprofLoadClass{}, fmt.Errorf("error in ParseHprofLoadClass: %w", err)
	}

	classObjectId, err := parser.primitiveParser.ParseIdentifier()
	if err != nil {
		return HprofLoadClass{}, fmt.Errorf("error in ParseHprofLoadClass: %w", err)
	}

	stackTraceSerialNumber, err := parser.primitiveParser.ParseUint32()
	if err != nil {
		return HprofLoadClass{}, fmt.Errorf("error in ParseHprofLoadClass: %w", err)
	}

	classNameId, err := parser.primitiveParser.ParseIdentifier()
	if err != nil {
		return HprofLoadClass{}, fmt.Errorf("error in ParseHprofLoadClass: %w", err)
	}
	return HprofLoadClass{
		ClassSerialNumber:      classSerialNumber,
		ClassObjectId:          classObjectId,
		StackTraceSerialNumber: stackTraceSerialNumber,
		ClassNameId:            classNameId,
	}, nil
}

// ParseHprofFrame reads HPROF_FRAME record.
func (parser *RecordParser) ParseHprofFrame() (HprofFrame, error) {
	stackFrameId, err := parser.primitiveParser.ParseIdentifier()
	if err != nil {
		return HprofFrame{}, fmt.Errorf("error in ParseHprofFrame: %w", err)
	}

	methodNameId, err := parser.primitiveParser.ParseIdentifier()
	if err != nil {
		return HprofFrame{}, fmt.Errorf("error in ParseHprofFrame: %w", err)
	}

	methodSignatureId, err := parser.primitiveParser.ParseIdentifier()
	if err != nil {
		return HprofFrame{}, fmt.Errorf("error in ParseHprofFrame: %w", err)
	}

	sourceFileNameId, err := parser.primitiveParser.ParseIdentifier()
	if err != nil {
		return HprofFrame{}, fmt.Errorf("error in ParseHprofFrame: %w", err)
	}

	classSerialNumber, err := parser.primitiveParser.ParseUint32()
	if err != nil {
		return HprofFrame{}, fmt.Errorf("error in ParseHprofFrame: %w", err)
	}

	lineNumber, err := parser.primitiveParser.ParseLineNumber()
	if err != nil {
		return HprofFrame{}, fmt.Errorf("error in ParseHprofFrame: %w", err)
	}
	return HprofFrame{
		StackFrameId:      stackFrameId,
		MethodNameId:      methodNameId,
		MethodSignatureId: methodSignatureId,
		SourceFileNameId:  sourceFileNameId,
		ClassSerialNumber: classSerialNumber,
		LineNumber:        lineNumber,
	}, nil
}

// ParseHprofTrace reads HPROF_TRACE record.
func (parser *RecordParser) ParseHprofTrace() (HprofTrace, error) {
	stackTraceSerialNumber, err := parser.primitiveParser.ParseUint32()
	if err != nil {
		return HprofTrace{}, fmt.Errorf("error in ParseHprofTrace: %w", err)
	}

	threadSerialNumber, err := parser.primitiveParser.ParseUint32()
	if err != nil {
		return HprofTrace{}, fmt.Errorf("error in ParseHprofTrace: %w", err)
	}

	numberOfFrames, err := parser.primitiveParser.ParseUint32()
	if err != nil {
		return HprofTrace{}, fmt.Errorf("error in ParseHprofTrace: %w", err)
	}

	var stackFrameIds []Identifier
	for i := numberOfFrames; i > 0; i-- {
		stackFrameId, err := parser.primitiveParser.ParseIdentifier()
		if err != nil {
			return HprofTrace{}, fmt.Errorf("error in ParseHprofTrace: %w", err)
		}
		stackFrameIds = append(stackFrameIds, stackFrameId)
	}

	return HprofTrace{
		StackTraceSerialNumber: stackTraceSerialNumber,
		ThreadSerialNumber:     threadSerialNumber,
		NumberOfFrames:         numberOfFrames,
		StackFrameIds:          stackFrameIds,
	}, nil
}

// ParseSubRecordHeader reads the type of the sub-record inside HPROF_HEAP_DUMP_SEGMENT.
func (parser *RecordParser) ParseSubRecordHeader() (SubRecordHeader, error) {
	subRecordType, err := parser.primitiveParser.ParseUint8()
	if err != nil {
		return SubRecordHeader{}, fmt.Errorf("error in ParseSubRecordHeader: %w", err)
	}
	return SubRecordHeader{SubRecordType: SubRecordType(subRecordType)}, nil
}

// ParseHprofGcRootThreadObj reads HPROF_GC_ROOT_THREAD_OBJ sub-record.
func (parser *RecordParser) ParseHprofGcRootThreadObj() (HprofGcRootThreadObj, error) {
	threadObjectId, err := parser.primitiveParser.ParseIdentifier()
	if err != nil {
		return HprofGcRootThreadObj{}, fmt.Errorf("error in ParseHprofGcRootThreadObj: %w", err)
	}

	threadSequenceNumber, err := parser.primitiveParser.ParseUint32()
	if err != nil {
		return HprofGcRootThreadObj{}, fmt.Errorf("error in ParseHprofGcRootThreadObj: %w", err)
	}

	stackTraceSequenceNumber, err := parser.primitiveParser.ParseUint32()
	if err != nil {
		return HprofGcRootThreadObj{}, fmt.Errorf("error in ParseHprofGcRootThreadObj: %w", err)
	}
	return HprofGcRootThreadObj{
		ThreadObjectId:           threadObjectId,
		ThreadSequenceNumber:     threadSequenceNumber,
		StackTraceSequenceNumber: stackTraceSequenceNumber,
	}, nil
}

// ParseHprofGcRootJniGlobal reads HPROF_GC_ROOT_JNI_GLOBAL sub-record.
func (parser *RecordParser) ParseHprofGcRootJniGlobal() (HprofGcRootJniGlobal, error) {
	objectId, err := parser.primitiveParser.ParseIdentifier()
	if err != nil {
		return HprofGcRootJniGlobal{}, fmt.Errorf("error in ParseHprofGcRootJniGlobal: %w", err)
	}

	jniGlobalRefId, err := parser.primitiveParser.ParseIdentifier()
	if err != nil {
		return HprofGcRootJniGlobal{}, fmt.Errorf("error in ParseHprofGcRootJniGlobal: %w", err)
	}
	return HprofGcRootJniGlobal{
		ObjectId:       objectId,
		JniGlobalRefId: jniGlobalRefId,
	}, nil
}

// ParseHprofGcRootJniLocal reads HPROF_GC_ROOT_JNI_LOCAL sub-record.
func (parser *RecordParser) ParseHprofGcRootJniLocal() (HprofGcRootJniLocal, error) {
	objectId, err := parser.primitiveParser.ParseIdentifier()
	if err != nil {
		return HprofGcRootJniLocal{}, fmt.Errorf("error in ParseHprofGcRootJniLocal: %w", err)
	}

	threadSerialNumber, err := parser.primitiveParser.ParseUint32()
	if err != nil {
		return HprofGcRootJniLocal{}, fmt.Errorf("error in ParseHprofGcRootJniLocal: %w", err)
	}

	frameNumberInStackTrace, err := parser.primitiveParser.ParseUint32()
	if err != nil {
		return HprofGcRootJniLocal{}, fmt.Errorf("error in ParseHprofGcRootJniLocal: %w", err)
	}
	return HprofGcRootJniLocal{
		ObjectId:                objectId,
		ThreadSerialNumber:      threadSerialNumber,
		FrameNumberInStackTrace: frameNumberInStackTrace,
	}, nil
}

// ParseHprofGcRootJavaFrame reads HPROF_GC_ROOT_JAVA_FRAME sub-record.
func (parser *RecordParser) ParseHprofGcRootJavaFrame() (HprofGcRootJavaFrame, error) {
	objectId, err := parser.primitiveParser.ParseIdentifier()
	if err != nil {
		return HprofGcRootJavaFrame{}, fmt.Errorf("error in ParseHprofGcRootJavaFrame: %w", err)
	}

	threadSerialNumber, err := parser.primitiveParser.ParseUint32()
	if err != nil {
		return HprofGcRootJavaFrame{}, fmt.Errorf("error in ParseHprofGcRootJavaFrame: %w", err)
	}

	frameNumberInStackTrace, err := parser.primitiveParser.ParseUint32()
	if err != nil {
		return HprofGcRootJavaFrame{}, fmt.Errorf("error in ParseHprofGcRootJavaFrame: %w", err)
	}
	return HprofGcRootJavaFrame{
		ObjectId:                objectId,
		ThreadSerialNumber:      threadSerialNumber,
		FrameNumberInStackTrace: frameNumberInStackTrace,
	}, nil
}

// ParseHprofGcRootStickyClass reads HPROF_GC_ROOT_STICKY_CLASS sub-record.
func (parser *RecordParser) ParseHprofGcRootStickyClass() (HprofGcRootStickyClass, error) {
	objectId, err := parser.primitiveParser.ParseIdentifier()
	if err != nil {
		return HprofGcRootStickyClass{}, fmt.Errorf("error in ParseHprofGcRootStickyClass: %w", err)
	}
	return HprofGcRootStickyClass{ObjectId: objectId}, nil
}

// ParseHprofGcClassDump reads HPROF_GC_CLASS_DUMP sub-record.
func (parser *RecordParser) ParseHprofGcClassDump() (HprofGcClassDump, error) {
	classObjectId, err := parser.primitiveParser.ParseIdentifier()
	if err != nil {
		return HprofGcClassDump{}, fmt.Errorf("error in ParseHprofGcClassDump: %w", err)
	}

	stackTraceSerialNumber, err := parser.primitiveParser.ParseUint32()
	if err != nil {
		return HprofGcClassDump{}, fmt.Errorf("error in ParseHprofGcClassDump: %w", err)
	}

	superclassObjectId, err := parser.primitiveParser.ParseIdentifier()
	if err != nil {
		return HprofGcClassDump{}, fmt.Errorf("error in ParseHprofGcClassDump: %w", err)
	}

	classloaderObjectId, err := parser.primitiveParser.ParseIdentifier()
	if err != nil {
		return HprofGcClassDump{}, fmt.Errorf("error in ParseHprofGcClassDump: %w", err)
	}

	signersObjectId, err := parser.primitiveParser.ParseIdentifier()
	if err != nil {
		return HprofGcClassDump{}, fmt.Errorf("error in ParseHprofGcClassDump: %w", err)
	}

	protectionDomainObjectId, err := parser.primitiveParser.ParseIdentifier()
	if err != nil {
		return HprofGcClassDump{}, fmt.Errorf("error in ParseHprofGcClassDump: %w", err)
	}

	// reserved 1
	_, err = parser.primitiveParser.ParseIdentifier()
	if err != nil {
		return HprofGcClassDump{}, fmt.Errorf("error in ParseHprofGcClassDump: %w", err)
	}

	// reserved 2
	_, err = parser.primitiveParser.ParseIdentifier()
	if err != nil {
		return HprofGcClassDump{}, fmt.Errorf("error in ParseHprofGcClassDump: %w", err)
	}

	instanceSize, err := parser.primitiveParser.ParseInt32()
	if err != nil {
		return HprofGcClassDump{}, fmt.Errorf("error in ParseHprofGcClassDump: %w", err)
	}

	sizeOfConstantPool, err := parser.primitiveParser.ParseUint16()
	if err != nil {
		return HprofGcClassDump{}, fmt.Errorf("error in ParseHprofGcClassDump: %w", err)
	}

	var constantPoolRecords []HprofGcClassDumpConstantPoolRecord
	for i := sizeOfConstantPool; i > 0; i-- {
		record, err := parser.parseHprofGcClassDumpConstantPoolRecord()
		if err != nil {
			return HprofGcClassDump{}, fmt.Errorf("error in ParseHprofGcClassDump: %w", err)
		}
		constantPoolRecords = append(constantPoolRecords, record)
	}

	numberOfStaticFields, err := parser.primitiveParser.ParseUint16()
	if err != nil {
		return HprofGcClassDump{}, fmt.Errorf("error in ParseHprofGcClassDump: %w", err)
	}

	var staticFieldRecords []HprofGcClassDumpStaticFieldsRecord
	for i := numberOfStaticFields; i > 0; i-- {
		record, err := parser.parseHprofGcClassDumpStaticFieldsRecord()
		if err != nil {
			return HprofGcClassDump{}, fmt.Errorf("error in ParseHprofGcClassDump: %w", err)
		}
		staticFieldRecords = append(staticFieldRecords, record)
	}

	numberOfInstanceFields, err := parser.primitiveParser.ParseUint16()
	if err != nil {
		return HprofGcClassDump{}, fmt.Errorf("error in ParseHprofGcClassDump: %w", err)
	}

	var instanceFieldRecords []HprofGcClassDumpInstanceFieldsRecord
	for i := numberOfInstanceFields; i > 0; i-- {
		record, err := parser.parseHprofGcClassDumpInstanceFieldsRecord()
		if err != nil {
			return HprofGcClassDump{}, fmt.Errorf("error in ParseHprofGcClassDump: %w", err)
		}
		instanceFieldRecords = append(instanceFieldRecords, record)
	}

	return HprofGcClassDump{
		ClassObjectId:            classObjectId,
		StackTraceSerialNumber:   stackTraceSerialNumber,
		SuperclassObjectId:       superclassObjectId,
		ClassloaderObjectId:      classloaderObjectId,
		SignersObjectId:          signersObjectId,
		ProtectionDomainObjectId: protectionDomainObjectId,
		InstanceSize:             instanceSize,
		SizeOfConstantPool:       sizeOfConstantPool,
		ConstantPoolRecords:      constantPoolRecords,
		NumberOfStaticFields:     numberOfStaticFields,
		StaticFieldRecords:       staticFieldRecords,
		NumberOfInstanceFields:   numberOfInstanceFields,
		InstanceFieldRecords:     instanceFieldRecords,
	}, nil
}

// parseHprofGcClassDumpConstantPoolRecord reads one constant pool record. Used
// by ParseHprofGcClassDump
func (parser *RecordParser) parseHprofGcClassDumpConstantPoolRecord() (HprofGcClassDumpConstantPoolRecord, error) {
	constantPoolIndex, err := parser.primitiveParser.ParseUint16()
	if err != nil {
		return HprofGcClassDumpConstantPoolRecord{}, fmt.Errorf("error in ParseHprofGcClassDumpConstantPoolRecord: %w", err)
	}
	javaType, err := parser.primitiveParser.ParseJavaType()
	if err != nil {
		return HprofGcClassDumpConstantPoolRecord{}, fmt.Errorf("error in ParseHprofGcClassDumpConstantPoolRecord: %w", err)
	}
	javaValue, err := parser.primitiveParser.ParseJavaValue(javaType)
	if err != nil {
		return HprofGcClassDumpConstantPoolRecord{}, fmt.Errorf("error in ParseHprofGcClassDumpConstantPoolRecord: %w", err)
	}
	return HprofGcClassDumpConstantPoolRecord{
		ConstantPoolIndex: constantPoolIndex,
		Ty:                javaType,
		Value:             javaValue,
	}, nil
}

// parseHprofGcClassDumpStaticFieldsRecord reads one static field record. Used
// by ParseHprofGcClassDump
func (parser *RecordParser) parseHprofGcClassDumpStaticFieldsRecord() (HprofGcClassDumpStaticFieldsRecord, error) {
	staticFieldName, err := parser.primitiveParser.ParseIdentifier()
	if err != nil {
		return HprofGcClassDumpStaticFieldsRecord{}, fmt.Errorf("error in ParseHprofGcClassDumpStaticFieldsRecord: %w", err)
	}
	javaType, err := parser.primitiveParser.ParseJavaType()
	if err != nil {
		return HprofGcClassDumpStaticFieldsRecord{}, fmt.Errorf("error in ParseHprofGcClassDumpStaticFieldsRecord: %w", err)
	}
	javaValue, err := parser.primitiveParser.ParseJavaValue(javaType)
	if err != nil {
		return HprofGcClassDumpStaticFieldsRecord{}, fmt.Errorf("error in ParseHprofGcClassDumpStaticFieldsRecord: %w", err)
	}
	return HprofGcClassDumpStaticFieldsRecord{
		StaticFieldName: staticFieldName,
		Ty:              javaType,
		Value:           javaValue,
	}, nil
}

// parseHprofGcClassDumpInstanceFieldsRecord reads one instance field record. Used
// by ParseHprofGcClassDump
func (parser *RecordParser) parseHprofGcClassDumpInstanceFieldsRecord() (HprofGcClassDumpInstanceFieldsRecord, error) {
	instanceFieldName, err := parser.primitiveParser.ParseIdentifier()
	if err != nil {
		return HprofGcClassDumpInstanceFieldsRecord{}, fmt.Errorf("error in ParseHprofGcClassDumpInstanceFieldsRecord: %w", err)
	}

	ty, err := parser.primitiveParser.ParseJavaType()
	if err != nil {
		return HprofGcClassDumpInstanceFieldsRecord{}, fmt.Errorf("error in ParseHprofGcClassDumpInstanceFieldsRecord: %w", err)
	}
	return HprofGcClassDumpInstanceFieldsRecord{
		InstanceFieldName: instanceFieldName,
		Ty:                ty,
	}, nil
}

// ParseHprofGcClassDumpInstanceDumpHeader reads "header" information from HPROF_GC_INSTANCE_DUMP.
// It does not read array of instance field values, so ParseHprofGcClassDumpInstanceDumpRecord does.
func (parser *RecordParser) ParseHprofGcClassDumpInstanceDumpHeader() (HprofGcClassDumpInstanceDumpHeader, error) {
	objectId, err := parser.primitiveParser.ParseIdentifier()
	if err != nil {
		return HprofGcClassDumpInstanceDumpHeader{}, fmt.Errorf("error in ParseHprofGcClassDumpInstanceDumpHeader: %w", err)
	}

	stackTraceSerialNumber, err := parser.primitiveParser.ParseUint32()
	if err != nil {
		return HprofGcClassDumpInstanceDumpHeader{}, fmt.Errorf("error in ParseHprofGcClassDumpInstanceDumpHeader: %w", err)
	}

	classObjectId, err := parser.primitiveParser.ParseIdentifier()
	if err != nil {
		return HprofGcClassDumpInstanceDumpHeader{}, fmt.Errorf("error in ParseHprofGcClassDumpInstanceDumpHeader: %w", err)
	}

	numberOfBytesThatFollow, err := parser.primitiveParser.ParseUint32()
	if err != nil {
		return HprofGcClassDumpInstanceDumpHeader{}, fmt.Errorf("error in ParseHprofGcClassDumpInstanceDumpHeader: %w", err)
	}
	return HprofGcClassDumpInstanceDumpHeader{
		ObjectId:                objectId,
		StackTraceSerialNumber:  stackTraceSerialNumber,
		ClassObjectId:           classObjectId,
		NumberOfBytesThatFollow: numberOfBytesThatFollow,
	}, nil
}

// ParseHprofGcObjArrayDumpHeader reads "header" information from HPROF_GC_OBJ_ARRAY_DUMP.
// It does not read array of elements, so ParseHprofGcObjArrayDumpRecord does.
func (parser *RecordParser) ParseHprofGcObjArrayDumpHeader() (HprofGcObjArrayDumpHeader, error) {
	arrayObjectId, err := parser.primitiveParser.ParseIdentifier()
	if err != nil {
		return HprofGcObjArrayDumpHeader{}, fmt.Errorf("error in ParseHprofGcObjArrayDumpHeader: %w", err)
	}

	stackTraceSerialNumber, err := parser.primitiveParser.ParseUint32()
	if err != nil {
		return HprofGcObjArrayDumpHeader{}, fmt.Errorf("error in ParseHprofGcObjArrayDumpHeader: %w", err)
	}

	numberOfElements, err := parser.primitiveParser.ParseUint32()
	if err != nil {
		return HprofGcObjArrayDumpHeader{}, fmt.Errorf("error in ParseHprofGcObjArrayDumpHeader: %w", err)
	}

	arrayClassId, err := parser.primitiveParser.ParseIdentifier()
	if err != nil {
		return HprofGcObjArrayDumpHeader{}, fmt.Errorf("error in ParseHprofGcObjArrayDumpHeader: %w", err)
	}
	return HprofGcObjArrayDumpHeader{
		ArrayObjectId:          arrayObjectId,
		StackTraceSerialNumber: stackTraceSerialNumber,
		NumberOfElements:       numberOfElements,
		ArrayClassId:           arrayClassId,
	}, nil
}

// ParseHprofGcPrimArrayDumpHeader reads "header" information from HPROF_GC_PRIM_ARRAY_DUMP.
// It does not read array of elements, so ParseHprofGcPrimArrayDumpRecord does.
func (parser *RecordParser) ParseHprofGcPrimArrayDumpHeader() (HprofGcPrimArrayDumpHeader, error) {
	arrayObjectId, err := parser.primitiveParser.ParseIdentifier()
	if err != nil {
		return HprofGcPrimArrayDumpHeader{}, fmt.Errorf("error in ParseHprofGcPrimArrayDumpHeader: %w", err)
	}

	stackTraceSerialNumber, err := parser.primitiveParser.ParseUint32()
	if err != nil {
		return HprofGcPrimArrayDumpHeader{}, fmt.Errorf("error in ParseHprofGcPrimArrayDumpHeader: %w", err)
	}

	numberOfElements, err := parser.primitiveParser.ParseUint32()
	if err != nil {
		return HprofGcPrimArrayDumpHeader{}, fmt.Errorf("error in ParseHprofGcPrimArrayDumpHeader: %w", err)
	}

	elementType, err := parser.primitiveParser.ParseJavaType()
	if err != nil {
		return HprofGcPrimArrayDumpHeader{}, fmt.Errorf("error in ParseHprofGcPrimArrayDumpHeader: %w", err)
	}
	return HprofGcPrimArrayDumpHeader{
		ArrayObjectId:          arrayObjectId,
		StackTraceSerialNumber: stackTraceSerialNumber,
		NumberOfElements:       numberOfElements,
		ElementType:            elementType,
	}, nil
}

// ReadBytes read n bytes from underlying io.Reader
func (parser *RecordParser) ReadBytes(n int) ([]byte, error) {
	res, err := parser.primitiveParser.ParseByteSeq(n)
	if err != nil {
		return nil, fmt.Errorf("error in ReadBytes: %w", err)
	}
	return res, nil
}
