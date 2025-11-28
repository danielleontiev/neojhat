package core

import (
	"bytes"
	"testing"
	"time"
)

func TestParseFileHeader(t *testing.T) {
	sampleTime, err := time.Parse(time.RFC3339, "2021-09-17T03:59:04.147+04:00")
	if err != nil {
		t.Errorf("Cannot parse time string, error: [%v]", err)
	}
	header := []byte{
		0x4a, 0x41, 0x56, 0x41, 0x20, 0x50, 0x52, 0x4f, 0x46, 0x49, 0x4c, 0x45, 0x20, 0x31, 0x2e, 0x30, 0x2e, 0x32, 0x00, // JAVA PROFILE 1.0.2
		0x00, 0x00, 0x00, 0x08, // identifier size
		0x00, 0x00, 0x01, 0x7b, // high word
		0xf1, 0x0c, 0xa9, 0xd3, // low word
	}
	tests := []struct {
		name    string
		input   []byte
		want    FileHeader
		wantErr bool
	}{
		{
			name:  "success parse file header",
			input: header,
			want: FileHeader{
				Header:         "JAVA PROFILE 1.0.2",
				IdentifierSize: 8,
				Timestamp:      sampleTime,
			},
		},
		{
			name:    "error parse file header",
			input:   empty,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFileHeader(bytes.NewReader(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFileHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !compareHeaders(got, tt.want) {
				t.Errorf("ParseFileHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func compareHeaders(this, that FileHeader) bool {
	str := this.Header == that.Header
	ids := this.IdentifierSize == that.IdentifierSize
	times := this.Timestamp.Equal(that.Timestamp)
	return str && ids && times
}
