package threads

import (
	"fmt"
	"sort"
	"strings"

	"github.com/danielleontiev/neojhat/core"
	"github.com/danielleontiev/neojhat/format"
	"github.com/danielleontiev/neojhat/printing"
)

// PrettyPrint prints given thread dump
// in beautiful manner
func PrettyPrint(threadDump ThreadDump, localVars bool) {
	traces := threadDump.StackTraces
	sort.Slice(traces, func(i, j int) bool {
		return traces[i].ThreadId < traces[j].ThreadId
	})
	for _, stackTrace := range traces {
		fmt.Println(createPrettyThread(stackTrace))
		for _, frame := range stackTrace.Frames {
			fmt.Printf("	%s\n", createPrettyFrame(frame))
			if localVars {
				for _, local := range frame.LocalFrames {
					fmt.Printf("		local %s\n", createPrettyStackVariable(local))
				}
			}
		}
		fmt.Println()
	}
}

func createPrettyThread(stackTrace StackTrace) string {
	threadDesc := fmt.Sprintf(
		"\"%v\", ID=%v, prio=%v, status=%v",
		stackTrace.ThreadName,
		stackTrace.ThreadId,
		stackTrace.ThreadPriority,
		stackTrace.ThreadStatus,
	)
	if stackTrace.ThreadDaemon {
		threadDesc += " (daemon)"
	}
	return threadDesc
}

func createPrettyFrame(frame StackFrame) string {
	prettyClassName := format.ClassName(frame.ClassName)
	args, ret := format.Signature(frame.MethodSignature)
	prettyLocation := createLocation(frame.FileName, frame.LineNumber)
	return ret + " " + format.ClassName(prettyClassName) + "." + frame.MethodName + "(" + args + ")" + " " + prettyLocation
}

func createPrettyStackVariable(localFrame LocalFrame) string {
	signature := localFrame.ObjectTypeSignature
	if !strings.HasPrefix(signature, "[") { // it's object
		signature = "L" + signature + ";"
	}
	if strings.HasPrefix(signature, "class") {
		return signature
	}
	arg, _ := format.Signature(signature)
	return arg
}

func createLocation(fileName, lineNumber string) string {
	if fileName == unknownString {
		return ""
	}
	if lineNumber == core.Unknown.String() {
		return fileName
	}
	return fmt.Sprintf("%s:%s", fileName, lineNumber)
}

// PrettyPrintColor prints given thread dump
// in beautiful manner with ANSI colors
func PrettyPrintColor(threadDump ThreadDump, localVars bool) {
	traces := threadDump.StackTraces
	sort.Slice(traces, func(i, j int) bool {
		return traces[i].ThreadId < traces[j].ThreadId
	})
	for _, stackTrace := range traces {
		fmt.Println(printing.Bold(createPrettyThread(stackTrace)))
		for _, frame := range stackTrace.Frames {
			fmt.Printf("	%s\n", createPrettyColorfulFrame(frame))
			if localVars {
				for _, local := range frame.LocalFrames {
					localString := createPrettyStackVariable(local)
					fmt.Printf("		local %s\n", printing.Blue(localString))
				}
			}
		}
		fmt.Println()
	}
}

func createPrettyColorfulFrame(frame StackFrame) string {
	prettyClassName := printing.Yellow(format.ClassName(frame.ClassName))
	args, ret := format.Signature(frame.MethodSignature)
	prettyLocation := printing.Cyan(createLocation(frame.FileName, frame.LineNumber))
	return ret + " " + prettyClassName + printing.Yellow(".") + printing.Red(frame.MethodName) + "(" + args + ")" + " " + prettyLocation
}
