package core

import (
	"testing"
)

func TestLineNumber_String(t *testing.T) {
	tests := []struct {
		name string
		l    LineNumber
		want string
	}{
		{
			name: "unknown",
			l:    LineNumber(-1),
			want: "Unknown",
		},
		{
			name: "compiled",
			l:    LineNumber(-2),
			want: "CompiledMethod",
		},
		{
			name: "native",
			l:    LineNumber(-3),
			want: "NativeMethod",
		},
		{
			name: "unrecognized",
			l:    LineNumber(-12),
			want: "Error",
		},
		{
			name: "positive",
			l:    12,
			want: "12",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.String(); got != tt.want {
				t.Errorf("LineNumber.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJavaType_String(t *testing.T) {
	tests := []struct {
		name string
		j    JavaType
		want string
	}{
		{
			name: "Object",
			j:    Object,
			want: "object",
		},
		{
			name: "Boolean",
			j:    Boolean,
			want: "boolean",
		},
		{
			name: "Char",
			j:    Char,
			want: "char",
		},
		{
			name: "Float",
			j:    Float,
			want: "float",
		},
		{
			name: "Double",
			j:    Double,
			want: "double",
		},
		{
			name: "Byte",
			j:    Byte,
			want: "byte",
		},
		{
			name: "Short",
			j:    Short,
			want: "short",
		},
		{
			name: "Int",
			j:    Int,
			want: "int",
		},
		{
			name: "Long",
			j:    Long,
			want: "long",
		},
		{
			name: "Unknown",
			j:    42,
			want: "unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.j.String(); got != tt.want {
				t.Errorf("JavaType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJavaValue_ToBool(t *testing.T) {
	tests := []struct {
		name    string
		jv      JavaValue
		want    bool
		wantErr bool
	}{
		{
			name: "success",
			jv:   JavaValue{Type: Boolean, Value: true},
			want: true,
		},
		{
			name:    "wrong type",
			jv:      JavaValue{Type: Int, Value: true},
			wantErr: true,
		},
		{
			name:    "wrong value",
			jv:      JavaValue{Type: Boolean, Value: 42},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.jv.ToBool()
			if (err != nil) != tt.wantErr {
				t.Errorf("JavaValue.ToBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("JavaValue.ToBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJavaValue_ToInt(t *testing.T) {
	tests := []struct {
		name    string
		jv      JavaValue
		want    int
		wantErr bool
	}{
		{
			name: "success",
			jv:   JavaValue{Type: Int, Value: int32(42)},
			want: 42,
		},
		{
			name:    "wrong type",
			jv:      JavaValue{Type: Long, Value: 42},
			wantErr: true,
		},
		{
			name:    "wrong value",
			jv:      JavaValue{Type: Int, Value: true},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.jv.ToInt()
			if (err != nil) != tt.wantErr {
				t.Errorf("JavaValue.ToInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("JavaValue.ToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJavaValue_ToLong(t *testing.T) {
	tests := []struct {
		name    string
		jv      JavaValue
		want    int
		wantErr bool
	}{
		{
			name: "success",
			jv:   JavaValue{Type: Long, Value: int64(42)},
			want: 42,
		},
		{
			name:    "wrong type",
			jv:      JavaValue{Type: Int, Value: 42},
			wantErr: true,
		},
		{
			name:    "wrong value",
			jv:      JavaValue{Type: Long, Value: false},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.jv.ToLong()
			if (err != nil) != tt.wantErr {
				t.Errorf("JavaValue.ToLong() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("JavaValue.ToLong() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJavaValue_ToObject(t *testing.T) {
	tests := []struct {
		name    string
		jv      JavaValue
		want    Identifier
		wantErr bool
	}{
		{
			name: "success",
			jv:   JavaValue{Type: Object, Value: Identifier(42)},
			want: 42,
		},
		{
			name:    "wrong type",
			jv:      JavaValue{Type: Int, Value: 42},
			wantErr: true,
		},
		{
			name:    "wrong value",
			jv:      JavaValue{Type: Object, Value: false},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.jv.ToObject()
			if (err != nil) != tt.wantErr {
				t.Errorf("JavaValue.ToObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("JavaValue.ToObject() = %v, want %v", got, tt.want)
			}
		})
	}
}
