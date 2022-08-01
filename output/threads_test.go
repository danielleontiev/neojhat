package output_test

import (
	_ "embed"

	"strings"
	"testing"

	"github.com/danielleontiev/neojhat/output"
	"github.com/danielleontiev/neojhat/threads"
)

var threads1 = threads.ThreadDump{
	StackTraces: []threads.StackTrace{
		{
			ThreadName:     "main",
			ThreadId:       1,
			ThreadDaemon:   false,
			ThreadPriority: 5,
			ThreadStatus:   threads.ThreadStateWaitingWithTimeout,
			NumberOfFrames: 2,
			Frames: []threads.StackFrame{
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
					LocalFrames: []threads.LocalFrame{
						{
							ObjectId:            1,
							ObjectTypeSignature: "[Ljava/lang/String;",
							Type:                threads.Frame,
						},
						{
							ObjectId:            2,
							ObjectTypeSignature: "java/lang/String",
							Type:                threads.Frame,
						},
					},
				},
			},
		},
	},
}

var (
	//go:embed test-data/threads1.txt
	threads1txt string
	//go:embed test-data/threads2.txt
	threads2txt string
	//go:embed test-data/threads1.html
	threads1html string
	//go:embed test-data/threads2.html
	threads2html string
)

func TestThreadPlain1(t *testing.T) {
	builder := &strings.Builder{}
	output.ThreadsPlain(threads1, true, builder)
	result := builder.String()
	if result != threads1txt {
		compareLineByLine(t, result, threads1txt)
	}
}

func TestThreadPlain2(t *testing.T) {
	builder := &strings.Builder{}
	output.ThreadsPlain(threads1, false, builder)
	result := builder.String()
	if result != threads2txt {
		compareLineByLine(t, result, threads2txt)
	}
}

func TestThreadHtml1(t *testing.T) {
	builder := &strings.Builder{}
	output.ThreadsHtml(threads1, true, builder)
	result := builder.String()
	if result != threads1html {
		compareLineByLine(t, result, threads1html)
	}
}

func TestThreadHtml2(t *testing.T) {
	builder := &strings.Builder{}
	output.ThreadsHtml(threads1, false, builder)
	result := builder.String()
	if result != threads2html {
		compareLineByLine(t, result, threads2html)
	}
}
