package threads

import (
	"testing"

	"github.com/danielleontiev/neojhat/core"
)

func TestCreatePrettyThread(t *testing.T) {
	stackTrace := StackTrace{
		ThreadName:     "thread",
		ThreadId:       1,
		ThreadDaemon:   true,
		ThreadPriority: 2,
		ThreadStatus:   1,
	}
	want := "\"thread\", ID=1, prio=2, status=NEW (daemon)"
	if got := createPrettyThread(stackTrace); got != want {
		t.Errorf("CreatePrettyThread() = %v, want %v", got, want)
	}
}

func TestCreatePrettyFrame(t *testing.T) {
	frame := StackFrame{
		MethodName:      "main",
		MethodSignature: "([Ljava/lang/String;)V",
		FileName:        "Main.java",
		ClassName:       "foo/bar/Main",
		LineNumber:      "42",
	}
	want := "void foo.bar.Main.main(java.lang.String[]) Main.java:42"
	if got := createPrettyFrame(frame); got != want {
		t.Errorf("CreatePrettyFrame() = %v, want %v", got, want)
	}
}

func Test_createLocation(t *testing.T) {
	tests := []struct {
		name       string
		fileName   string
		lineNumber string
		want       string
	}{
		{
			name:     "unknown file",
			fileName: unknownString,
			want:     "",
		},
		{
			name:       "unknown line",
			fileName:   "Main.java",
			lineNumber: core.Unknown.String(),
			want:       "Main.java",
		},
		{
			name:       "all known",
			fileName:   "Main.java",
			lineNumber: "42",
			want:       "Main.java:42",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createLocation(tt.fileName, tt.lineNumber); got != tt.want {
				t.Errorf("createLocation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createPrettyStackVariable(t *testing.T) {
	tests := []struct {
		name       string
		localFrame LocalFrame
		want       string
	}{
		{
			name:       "object",
			localFrame: LocalFrame{ObjectTypeSignature: "java/lang/String"},
			want:       "java.lang.String",
		},
		{
			name:       "array",
			localFrame: LocalFrame{ObjectTypeSignature: "[B"},
			want:       "byte[]",
		},
		{
			name:       "class",
			localFrame: LocalFrame{ObjectTypeSignature: "class Main"},
			want:       "class Main",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createPrettyStackVariable(tt.localFrame); got != tt.want {
				t.Errorf("createPrettyStackVariable() = %v, want %v", got, tt.want)
			}
		})
	}
}
