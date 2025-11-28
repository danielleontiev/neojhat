package core

import (
	"reflect"
	"testing"
)

func TestPrimitiveParser_ParseUint32(t *testing.T) {
	tests := []struct {
		name    string
		parser  PrimitiveParser
		want    uint32
		wantErr bool
	}{
		{
			name:    "success",
			parser:  createPrimitiveParser(one4, CreateOpts{}),
			want:    1,
			wantErr: false,
		},
		{
			name:    "error",
			parser:  createPrimitiveParser(empty, CreateOpts{}),
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.parser.ParseUint32()
			if (err != nil) != tt.wantErr {
				t.Errorf("PrimitiveParser.ParseUint32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PrimitiveParser.ParseUint32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrimitiveParser_ParseUint8(t *testing.T) {
	tests := []struct {
		name    string
		parser  PrimitiveParser
		want    uint8
		wantErr bool
	}{
		{
			name:    "success",
			parser:  createPrimitiveParser(one1, CreateOpts{}),
			want:    1,
			wantErr: false,
		},
		{
			name:    "error",
			parser:  createPrimitiveParser(empty, CreateOpts{}),
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.parser.ParseUint8()
			if (err != nil) != tt.wantErr {
				t.Errorf("PrimitiveParser.ParseUint8() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PrimitiveParser.ParseUint8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrimitiveParser_ParseIdentifier(t *testing.T) {
	tests := []struct {
		name    string
		parser  PrimitiveParser
		want    Identifier
		wantErr bool
	}{
		{
			name:    "success parse id with length 4 (unit32)",
			parser:  createPrimitiveParser(one4, CreateOpts{idSize: 4}),
			want:    1,
			wantErr: false,
		},
		{
			name:    "success parse id with length 8 (unit64)",
			parser:  createPrimitiveParser(concat(zero4, one4), CreateOpts{idSize: 8}),
			want:    1,
			wantErr: false,
		},
		{
			name:    "error on not enough bytes with length 4",
			parser:  createPrimitiveParser(empty, CreateOpts{idSize: 4}),
			want:    0,
			wantErr: true,
		},
		{
			name:    "error on not enough bytes with length 8",
			parser:  createPrimitiveParser(empty, CreateOpts{idSize: 8}),
			want:    0,
			wantErr: true,
		},
		{
			name:    "error on length != 4 or 8",
			parser:  createPrimitiveParser(empty, CreateOpts{idSize: 7}),
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.parser.ParseIdentifier()
			if (err != nil) != tt.wantErr {
				t.Errorf("PrimitiveParser.ParseIdentifier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PrimitiveParser.ParseIdentifier() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrimitiveParser_ParseInt32(t *testing.T) {
	tests := []struct {
		name    string
		parser  PrimitiveParser
		want    int32
		wantErr bool
	}{
		{
			name:    "success positive",
			parser:  createPrimitiveParser(one4, CreateOpts{}),
			want:    1,
			wantErr: false,
		},
		{
			name:    "success negative",
			parser:  createPrimitiveParser([]byte{0xff, 0xff, 0xff, 0xff}, CreateOpts{}),
			want:    -1,
			wantErr: false,
		},
		{
			name:    "error parsing",
			parser:  createPrimitiveParser(empty, CreateOpts{}),
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.parser.ParseInt32()
			if (err != nil) != tt.wantErr {
				t.Errorf("PrimitiveParser.ParseInt32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PrimitiveParser.ParseInt32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrimitiveParser_ParseUint16(t *testing.T) {
	tests := []struct {
		name    string
		parser  PrimitiveParser
		want    uint16
		wantErr bool
	}{
		{
			name:   "success",
			parser: createPrimitiveParser(one2, CreateOpts{}),
			want:   1,
		},
		{
			name:    "error",
			parser:  createPrimitiveParser(empty, CreateOpts{}),
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.parser.ParseUint16()
			if (err != nil) != tt.wantErr {
				t.Errorf("PrimitiveParser.ParseUint16() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PrimitiveParser.ParseUint16() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrimitiveParser_lineNumber(t *testing.T) {
	tests := []struct {
		name    string
		parser  PrimitiveParser
		want    LineNumber
		wantErr bool
	}{
		{
			name:   "success",
			parser: createPrimitiveParser(one4, CreateOpts{}),
			want:   1,
		},
		{
			name:    "error",
			parser:  createPrimitiveParser(empty, CreateOpts{}),
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.parser.ParseLineNumber()
			if (err != nil) != tt.wantErr {
				t.Errorf("PrimitiveParser.lineNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PrimitiveParser.lineNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrimitiveParser_ParseJavaType(t *testing.T) {
	tests := []struct {
		name    string
		parser  PrimitiveParser
		want    JavaType
		wantErr bool
	}{
		{
			name:    "success",
			parser:  createPrimitiveParser(one1, CreateOpts{}),
			want:    1,
			wantErr: false,
		},
		{
			name:    "error",
			parser:  createPrimitiveParser(empty, CreateOpts{}),
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.parser.ParseJavaType()
			if (err != nil) != tt.wantErr {
				t.Errorf("PrimitiveParser.ParseJavaType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PrimitiveParser.ParseJavaType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrimitiveParser_ParseJavaValue(t *testing.T) {
	tests := []struct {
		name    string
		parser  PrimitiveParser
		ty      JavaType
		want    any
		wantErr bool
	}{
		{
			name:   "object",
			parser: createPrimitiveParser(one8, CreateOpts{idSize: 8}),
			ty:     Object,
			want:   Identifier(1),
		},
		{
			name:   "boolean",
			parser: createPrimitiveParser(one1, CreateOpts{}),
			ty:     Boolean,
			want:   true,
		},
		{
			name:   "char",
			parser: createPrimitiveParser([]byte{0x00, 0x61}, CreateOpts{}),
			ty:     Char,
			want:   string([]byte{0x00, 0x61}),
		},
		{
			name:   "float",
			parser: createPrimitiveParser(one4, CreateOpts{}),
			ty:     Float,
			want:   float32(1e-45),
		},
		{
			name:   "double",
			parser: createPrimitiveParser(one8, CreateOpts{}),
			ty:     Double,
			want:   5e-324,
		},
		{
			name:   "byte",
			parser: createPrimitiveParser(one1, CreateOpts{}),
			ty:     Byte,
			want:   int8(1),
		},
		{
			name:   "short",
			parser: createPrimitiveParser(one2, CreateOpts{}),
			ty:     Short,
			want:   int16(1),
		},
		{
			name:   "int",
			parser: createPrimitiveParser(one4, CreateOpts{}),
			ty:     Int,
			want:   int32(1),
		},
		{
			name:   "long",
			parser: createPrimitiveParser(one8, CreateOpts{}),
			ty:     Long,
			want:   int64(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.parser.ParseJavaValue(tt.ty)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrimitiveParser.ParseJavaValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Value != tt.want {
				t.Errorf("PrimitiveParser.ParseJavaValue() = %v (%T), want %v (%T)", got.Value, got.Value, tt.want, tt.want)
			}
		})
	}
}

func TestPrimitiveParser_ParseByteSeq(t *testing.T) {
	tests := []struct {
		name    string
		parser  PrimitiveParser
		n       int
		want    []byte
		wantErr bool
	}{
		{
			name:   "success",
			parser: createPrimitiveParser([]byte{0x00, 0x01}, CreateOpts{}),
			n:      2,
			want:   []byte{0x00, 0x01},
		},
		{
			name:    "error",
			parser:  createPrimitiveParser([]byte{0x00, 0x01}, CreateOpts{}),
			n:       4,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.parser.ParseByteSeq(tt.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrimitiveParser.ParseByteSeq() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PrimitiveParser.ParseByteSeq() = %v, want %v", got, tt.want)
			}
		})
	}
}
