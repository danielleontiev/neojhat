package threads

type ThreadDump struct {
	StackTraces []StackTrace
}

type StackTrace struct {
	ThreadName     string
	ThreadId       int
	ThreadDaemon   bool
	ThreadPriority int
	ThreadStatus   ThreadStatus
	NumberOfFrames uint32
	Frames         []StackFrame
}

type StackFrame struct {
	MethodName      string
	MethodSignature string
	FileName        string
	ClassName       string
	LineNumber      string
	LocalFrames     []LocalFrame
}

type LocalFrame struct {
	ObjectId            int
	ObjectTypeSignature string
	Type                FrameType
}

type FrameType int

const (
	JniLocal = iota
	Frame
)

const (
	ThreadStateAlive                 = 0x0001
	ThreadStateTerminated            = 0x0002
	ThreadStateRunnable              = 0x0004
	ThreadStateBlockedOnMonitorEnter = 0x0400
	ThreadStateWaitingIndefinitely   = 0x0010
	ThreadStateWaitingWithTimeout    = 0x0020
)

type ThreadStatus int

func (ts ThreadStatus) String() string {
	if ts&ThreadStateRunnable != 0 {
		return "RUNNABLE"
	} else if ts&ThreadStateBlockedOnMonitorEnter != 0 {
		return "BLOCKED"
	} else if ts&ThreadStateWaitingIndefinitely != 0 {
		return "WAITING"
	} else if ts&ThreadStateWaitingWithTimeout != 0 {
		return "TIMED_WAITING"
	} else if ts&ThreadStateTerminated != 0 {
		return "TERMINATED"
	} else if ts&ThreadStateAlive != 0 {
		return "NEW"
	} else {
		return "RUNNABLE"
	}
}
