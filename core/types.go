package core

import "fmt"

// Identifier represents Java identifier. uint64 is
// used to store it because actual size can be 4 or 8 bytes.
type Identifier uint64

// LineNumber is the wrapper for int32 with verbose values
// for some constants (Unknown, CompiledMethod, NativeMethod).
// It's stringified integer if > 0.
type LineNumber int32

const (
	Unknown        LineNumber = -1
	CompiledMethod LineNumber = -2
	NativeMethod   LineNumber = -3
)

var lineNumberMap = map[LineNumber]string{
	Unknown:        "Unknown",
	CompiledMethod: "CompiledMethod",
	NativeMethod:   "NativeMethod",
}

func (l LineNumber) String() string {
	if l > 0 {
		return fmt.Sprintf("%d", l)
	}
	str, ok := lineNumberMap[l]
	if ok {
		return str
	}
	return "Error"
}

// JavaType if representation of Java types
// with string values for readability.
type JavaType uint8

const (
	Object  JavaType = 2
	Boolean JavaType = 4
	Char    JavaType = 5
	Float   JavaType = 6
	Double  JavaType = 7
	Byte    JavaType = 8
	Short   JavaType = 9
	Int     JavaType = 10
	Long    JavaType = 11
)

var javaTypeMap = map[JavaType]string{
	Object:  "object",
	Boolean: "boolean",
	Char:    "char",
	Float:   "float",
	Double:  "double",
	Byte:    "byte",
	Short:   "short",
	Int:     "int",
	Long:    "long",
}

func (j JavaType) String() string {
	str, ok := javaTypeMap[j]
	if ok {
		return str
	}
	return "unknown"
}

// JavaValue wraps type of the
// constant and the value converted
// to corresponding Go type.
type JavaValue struct {
	Type  JavaType
	Value any
}

func (jv JavaValue) ToBool() (bool, error) {
	if jv.Type == Boolean {
		b, ok := jv.Value.(bool)
		if ok {
			return b, nil
		}
		return false, fmt.Errorf("malformed boolen value %v (%T)", jv.Value, jv.Value)
	}
	return false, fmt.Errorf("given value %v is not bool", jv)
}

func (jv JavaValue) ToInt() (int, error) {
	if jv.Type == Int {
		i, ok := jv.Value.(int32)
		if ok {
			return int(i), nil
		}
		return 0, fmt.Errorf("malformed int value %v (%T)", jv.Value, jv.Value)
	}
	return 0, fmt.Errorf("given value %v is not int", jv)
}

func (jv JavaValue) ToLong() (int, error) {
	if jv.Type == Long {
		i, ok := jv.Value.(int64)
		if ok {
			return int(i), nil
		}
		return 0, fmt.Errorf("malformed long value %v (%T)", jv.Value, jv.Value)
	}
	return 0, fmt.Errorf("given value %v is not long", jv)
}

func (jv JavaValue) ToObject() (Identifier, error) {
	if jv.Type == Object {
		i, ok := jv.Value.(Identifier)
		if ok {
			return i, nil
		}
		return 0, fmt.Errorf("malformed identifier value %v (%T)", jv.Value, jv.Value)
	}
	return 0, fmt.Errorf("given value %v is not identifier", jv)
}
