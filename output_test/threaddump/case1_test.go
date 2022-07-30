package threaddump_test

import (
	_ "embed"

	"github.com/danielleontiev/neojhat/threaddump"
)

var td1 = threaddump.ThreadDump{
	StackTraces: []threaddump.StackTrace{
		{
			ThreadName:     "main",
			ThreadId:       1,
			ThreadDaemon:   false,
			ThreadPriority: 5,
			ThreadStatus:   threaddump.ThreadStateWaitingWithTimeout,
			NumberOfFrames: 2,
			Frames: []threaddump.StackFrame{
				{
					MethodName:      "sleep",
					MethodSignature: "(J)V",
					FileName:        "Thread.java",
					ClassName:       "java/lang/Thread",
					LineNumber:      "NativeMethod",
					LocalFrames:     nil,
				},
				{
					MethodName:      "main",
					MethodSignature: "([Ljava/lang/String;)V",
					FileName:        "Main.java",
					ClassName:       "Main",
					LineNumber:      "6",
					LocalFrames: []threaddump.LocalFrame{
						{
							ObjectId:            1,
							ObjectTypeSignature: "[Ljava/lang/String;",
							Type:                threaddump.Frame,
						},
						{
							ObjectId:            2,
							ObjectTypeSignature: "java/lang/String",
							Type:                threaddump.Frame,
						},
					},
				},
			},
		},
	},
}

//go:embed case1.txt
var out1 string

//go:embed case1-no-local.txt
var out1noLocal string
