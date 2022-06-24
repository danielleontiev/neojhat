package format

import (
	"testing"
)

func TestClassName(t *testing.T) {
	tests := []struct {
		className string
		want      string
	}{
		{
			className: "Main",
			want:      "Main",
		},
		{
			className: "java/lang/String",
			want:      "java.lang.String",
		},
	}
	for _, tt := range tests {
		t.Run(tt.className, func(t *testing.T) {
			if got := ClassName(tt.className); got != tt.want {
				t.Errorf("ClassName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSignature(t *testing.T) {
	tests := []struct {
		signature string
		wantArg   string
		wantRet   string
	}{
		{
			signature: "()V",
			wantArg:   "",
			wantRet:   "void",
		},
		{
			signature: "(BCDFIJSZ)V",
			wantArg:   "byte, char, double, float, int, long, short, boolean",
			wantRet:   "void",
		},
		{
			signature: "(Ljava/lang/String;)Ljava/lang/Object;",
			wantArg:   "java.lang.String",
			wantRet:   "java.lang.Object",
		},
		{
			signature: "([B[[C[[[D)V",
			wantArg:   "byte[], char[][], double[][][]",
			wantRet:   "void",
		},
		{
			signature: "[B[[C[[[D",
			wantArg:   "byte[], char[][], double[][][]",
			wantRet:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.signature, func(t *testing.T) {
			arg, ret := Signature(tt.signature)
			if arg != tt.wantArg {
				t.Errorf("Signature() got = %v, want %v", arg, tt.wantArg)
			}
			if ret != tt.wantRet {
				t.Errorf("Signature() got1 = %v, want %v", ret, tt.wantRet)
			}
		})
	}
}

func TestSize(t *testing.T) {
	tests := []struct {
		name  string
		bytes int
		want  string
	}{
		{
			name:  "b",
			bytes: 42,
			want:  "42B",
		},
		{
			name:  "k",
			bytes: 2456,
			want:  "2K",
		},
		{
			name:  "m",
			bytes: 1234987,
			want:  "1M",
		},
		{
			name:  "g",
			bytes: 4365876354,
			want:  "4G",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Size(tt.bytes); got != tt.want {
				t.Errorf("Size() = %v, want %v", got, tt.want)
			}
		})
	}
}
