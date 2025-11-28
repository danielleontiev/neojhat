// threads extracts thread dump of JVM at the moment heap dump
// was taken and outputs it.
//
// Collecting the thread dump from .hprof file is the following
// process:
//
//  1. List all HPROF_GC_ROOT_THREAD_OBJ
//  2. Find corresponding HPROF_LOAD_CLASS
//  3. Find all HPROF_GC_ROOT_JAVA_FRAME
//  4. For each HPROF_GC_ROOT_THREAD_OBJ collect all HPROF_TRACE of
//     the thread and for each HPROF_TRACE collect all the HPROF_FRAME
//  5. For each HPROF_FRAME find corresponding HPROF_GC_ROOT_JAVA_FRAME
//     records and match them by frame number in stack trace.
//  6. Resolve all names presented using HPROF_UTF8
//  7. For each thread read the instance and class values from the heap
//     (meaning HPROF_GC_CLASS_DUMP, HPROF_GC_INSTANCE_DUMP) and extract
//     useful information such as thread name, thread id, thread priority,
//     and so on.
//  8. Provide the informative output for the above
package threads

import (
	"github.com/danielleontiev/neojhat/internal/core"
	"github.com/danielleontiev/neojhat/internal/dump"
	"github.com/danielleontiev/neojhat/internal/java"
)

var (
	UnknownString = "<unknown string>"
)

// GetThreadDump implements the whole collecting process described above.
func GetThreadDump(parsedAccessor *dump.ParsedAccessor) (ThreadDump, error) {
	heap := java.NewHeap(parsedAccessor)

	type threadSerialNumber int
	type positionInStack int
	var localFrames = map[threadSerialNumber]map[positionInStack][]LocalFrame{}
	initNestedMap := func(outerIndex threadSerialNumber) {
		_, ok := localFrames[outerIndex]
		if !ok {
			localFrames[outerIndex] = make(map[positionInStack][]LocalFrame)
		}
	}

	readObjectName := func(id core.Identifier) (string, error) {
		object, err := parsedAccessor.GetHprofGcInstanceDump(id)
		if err != nil {
			// we fall here when local var is not simple object
			objectArr, err := parsedAccessor.GetHprofGcObjArray(id)
			if err != nil {
				// try prim array
				primArr, err := parsedAccessor.GetHprofGcPrimArray(id)
				if err != nil {
					// then it's Class<?>
					loadClass, err := parsedAccessor.GetHprofLoadClassByClassObjectId(id)
					if err != nil {
						return "", err
					}
					className, err := parsedAccessor.GetHprofUtf8(loadClass.ClassNameId)
					if err != nil {
						return "", err
					}
					return "class " + className.Characters, nil
				}
				return "[" + getTypeSignature(primArr.ElementType), nil
			}
			object.ClassObjectId = objectArr.ArrayClassId
		}
		loadClass, err := parsedAccessor.GetHprofLoadClassByClassObjectId(object.ClassObjectId)
		if err != nil {
			return "", err
		}
		className, err := parsedAccessor.GetHprofUtf8(loadClass.ClassNameId)
		if err != nil {
			return "", err
		}
		return className.Characters, nil
	}

	jniLocals := parsedAccessor.ListHprofGcRootJniLocal()
	for _, jniLocal := range jniLocals {
		objectName, err := readObjectName(jniLocal.ObjectId)
		if err != nil {
			return ThreadDump{}, err
		}
		frame := LocalFrame{ObjectId: int(jniLocal.ObjectId), ObjectTypeSignature: objectName, Type: JniLocal}
		tn := threadSerialNumber(jniLocal.ThreadSerialNumber)
		pos := positionInStack(jniLocal.FrameNumberInStackTrace)
		initNestedMap(tn)
		localFrames[tn][pos] = append(localFrames[tn][pos], frame)
	}

	javaFrames := parsedAccessor.ListHprofGcRootJavaFrame()
	for _, javaFrame := range javaFrames {
		objectName, err := readObjectName(javaFrame.ObjectId)
		if err != nil {
			return ThreadDump{}, err
		}
		frame := LocalFrame{ObjectId: int(javaFrame.ObjectId), ObjectTypeSignature: objectName, Type: Frame}
		tn := threadSerialNumber(javaFrame.ThreadSerialNumber)
		pos := positionInStack(javaFrame.FrameNumberInStackTrace)
		initNestedMap(tn)
		localFrames[tn][pos] = append(localFrames[tn][pos], frame)
	}

	threadObjects := parsedAccessor.ListHprofGcRootThreadObj()
	var stackTraces []StackTrace
	for _, threadObj := range threadObjects {
		threadInstance, err := heap.ParseNormalObject(threadObj.ThreadObjectId)
		if err != nil {
			return ThreadDump{}, err
		}
		threadName, err := threadInstance.GetFieldValueByName("name")
		if err != nil {
			return ThreadDump{}, err
		}
		threadNameString, err := heap.ParseJavaString(threadName.Value)
		if err != nil {
			return ThreadDump{}, err
		}
		daemon, err := threadInstance.GetFieldValueByName("daemon")
		if err != nil {
			return ThreadDump{}, err
		}
		daemonBool, err := daemon.Value.ToBool()
		if err != nil {
			return ThreadDump{}, err
		}
		priority, err := threadInstance.GetFieldValueByName("priority")
		if err != nil {
			return ThreadDump{}, err
		}
		priorityInt, err := priority.Value.ToInt()
		if err != nil {
			return ThreadDump{}, err
		}
		threadId, err := threadInstance.GetFieldValueByName("tid")
		if err != nil {
			return ThreadDump{}, err
		}
		threadIdLong, err := threadId.Value.ToLong()
		if err != nil {
			return ThreadDump{}, err
		}
		threadStatus, err := threadInstance.GetFieldValueByName("threadStatus")
		if err != nil {
			return ThreadDump{}, err
		}
		threadStatusInt, err := threadStatus.Value.ToInt()
		if err != nil {
			return ThreadDump{}, err
		}
		stackTrace, err := parsedAccessor.GetHprofTrace(threadObj.ThreadSequenceNumber)
		if err != nil {
			return ThreadDump{}, err
		}
		var stackFrames []StackFrame
		for position, frameId := range stackTrace.StackFrameIds {
			stackFrame, err := parsedAccessor.GetHprofFrame(frameId)
			if err != nil {
				return ThreadDump{}, err
			}
			methodName, err := parsedAccessor.GetHprofUtf8(stackFrame.MethodNameId)
			if err != nil {
				methodName = core.HprofUtf8{Characters: UnknownString}
			}
			methodSignature, err := parsedAccessor.GetHprofUtf8(stackFrame.MethodSignatureId)
			if err != nil {
				methodSignature = core.HprofUtf8{Characters: UnknownString}
			}
			fileName, err := parsedAccessor.GetHprofUtf8(stackFrame.SourceFileNameId)
			if err != nil {
				fileName = core.HprofUtf8{Characters: UnknownString}
			}
			class, err := parsedAccessor.GetHprofLoadClassByClassSerialNumer(stackFrame.ClassSerialNumber)
			if err != nil {
				return ThreadDump{}, err
			}
			className, err := parsedAccessor.GetHprofUtf8(class.ClassNameId)
			if err != nil {
				className = core.HprofUtf8{Characters: UnknownString}
			}
			threadSequenceNumber := threadSerialNumber(threadObj.ThreadSequenceNumber)
			positionInStackFrame := positionInStack(position)
			stackLocalFrames := localFrames[threadSequenceNumber][positionInStackFrame]
			frame := StackFrame{
				MethodName:      methodName.Characters,
				MethodSignature: methodSignature.Characters,
				FileName:        fileName.Characters,
				ClassName:       className.Characters,
				LineNumber:      stackFrame.LineNumber.String(),
				LocalFrames:     stackLocalFrames,
			}
			stackFrames = append(stackFrames, frame)
		}
		trace := StackTrace{
			ThreadName:     threadNameString,
			ThreadId:       threadIdLong,
			ThreadDaemon:   daemonBool,
			ThreadPriority: priorityInt,
			ThreadStatus:   ThreadStatus(threadStatusInt),
			NumberOfFrames: stackTrace.NumberOfFrames,
			Frames:         stackFrames,
		}
		stackTraces = append(stackTraces, trace)
	}
	return ThreadDump{
		StackTraces: stackTraces,
	}, nil
}

var signaturesMap = map[core.JavaType]string{
	core.Object:  "L",
	core.Boolean: "Z",
	core.Char:    "C",
	core.Float:   "F",
	core.Double:  "D",
	core.Byte:    "B",
	core.Short:   "S",
	core.Int:     "I",
	core.Long:    "J",
}

func getTypeSignature(j core.JavaType) string {
	str, ok := signaturesMap[j]
	if ok {
		return str
	}
	return "unknown"
}
