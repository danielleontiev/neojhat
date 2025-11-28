package output

import (
	_ "embed"

	"fmt"
	"html/template"
	"io"
	"sort"
	"strings"

	"github.com/danielleontiev/neojhat/internal/core"
	"github.com/danielleontiev/neojhat/internal/format"
	"github.com/danielleontiev/neojhat/internal/threads"
)

// ThreadsPlain prints given thread dump
// in beautiful manner
func ThreadsPlain(threadDump threads.ThreadDump, localVars bool, destination io.Writer) {
	traces := getSortedStackTraces(threadDump)
	for _, stackTrace := range traces {
		fmt.Fprintln(destination, createPrettyThread(stackTrace))
		for _, frame := range stackTrace.Frames {
			fmt.Fprintf(destination, "    %s\n", createPrettyFrame(frame))
			if localVars {
				for _, local := range frame.LocalFrames {
					fmt.Fprintf(destination, "        local %s\n", createPrettyStackVariable(local))
				}
			}
		}
		fmt.Fprintln(destination)
	}
}

func createPrettyThread(stackTrace threads.StackTrace) string {
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

func createPrettyFrame(frame threads.StackFrame) string {
	prettyClassName := format.ClassName(frame.ClassName)
	args, ret := format.Signature(frame.MethodSignature)
	prettyLocation := createLocation(frame.FileName, frame.LineNumber)
	return ret + " " + format.ClassName(prettyClassName) + "." + frame.MethodName + "(" + args + ")" + " " + prettyLocation
}

func createPrettyStackVariable(localFrame threads.LocalFrame) string {
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
	if fileName == threads.UnknownString {
		return ""
	}
	if lineNumber == core.Unknown.String() {
		return fileName
	}
	return fmt.Sprintf("%s:%s", fileName, lineNumber)
}

func getSortedStackTraces(threadDump threads.ThreadDump) []threads.StackTrace {
	traces := threadDump.StackTraces
	sort.Slice(traces, func(i, j int) bool {
		return traces[i].ThreadId < traces[j].ThreadId
	})
	return traces
}

// ThreadsPlainColor prints given thread dump
// in beautiful manner with ANSI colors
func ThreadsPlainColor(threadDump threads.ThreadDump, localVars bool) {
	traces := getSortedStackTraces(threadDump)
	for _, stackTrace := range traces {
		fmt.Println(Bold(createPrettyThread(stackTrace)))
		for _, frame := range stackTrace.Frames {
			fmt.Printf("	%s\n", createPrettyColorfulFrame(frame))
			if localVars {
				for _, local := range frame.LocalFrames {
					localString := createPrettyStackVariable(local)
					fmt.Printf("		local %s\n", Blue(localString))
				}
			}
		}
		fmt.Println()
	}
}

func createPrettyColorfulFrame(frame threads.StackFrame) string {
	prettyClassName := Yellow(format.ClassName(frame.ClassName))
	args, ret := format.Signature(frame.MethodSignature)
	prettyLocation := Cyan(createLocation(frame.FileName, frame.LineNumber))
	return ret + " " + prettyClassName + Yellow(".") + Red(frame.MethodName) + "(" + args + ")" + " " + prettyLocation
}

var (
	//go:embed templates/threads.html
	threadsHtml string
)

type printTrace struct {
	ThreadName     string
	ThreadId       int
	ThreadDaemon   bool
	ThreadPriority int
	ThreadStatus   string
	Frames         []printFrame
}

type printFrame struct {
	ClassName   string
	MethodName  string
	Args        string
	Ret         string
	Location    string
	LocalFrames []string
}

// ThreadsHtml prints the output of summary command in nice
// beautifully-formatted HTML
func ThreadsHtml(threadDump threads.ThreadDump, localVars bool, destination io.Writer) error {
	coreTemplate, err := template.New("core").Parse(coreHtml)
	if err != nil {
		return err
	}
	threadsTemplate, err := coreTemplate.Parse(threadsHtml)
	if err != nil {
		return err
	}
	var traces []printTrace
	for _, t := range getSortedStackTraces(threadDump) {
		var frames []printFrame
		for _, f := range t.Frames {
			var stackVariables []string
			for _, s := range f.LocalFrames {
				stackVariables = append(stackVariables, createPrettyStackVariable(s))
			}
			args, ret := format.Signature(f.MethodSignature)
			loc := createLocation(f.FileName, f.LineNumber)
			frames = append(frames, printFrame{
				ClassName:   format.ClassName(f.ClassName),
				MethodName:  f.MethodName,
				Args:        args,
				Ret:         ret,
				Location:    loc,
				LocalFrames: stackVariables,
			})
		}
		traces = append(traces, printTrace{
			ThreadName:     t.ThreadName,
			ThreadId:       t.ThreadId,
			ThreadPriority: t.ThreadPriority,
			ThreadStatus:   t.ThreadStatus.String(),
			ThreadDaemon:   t.ThreadDaemon,
			Frames:         frames,
		})
	}
	return threadsTemplate.Execute(destination, data{
		Title:   "Thread Dump",
		Favicon: faviconBase64,
		Payload: struct {
			Traces    []printTrace
			LocalVars bool
		}{
			Traces:    traces,
			LocalVars: localVars,
		},
	})
}
