package java

import (
	"github.com/danielleontiev/neojhat/internal/core"
)

// StaticField is representation
// of class'es static field name, type and
// value, obtained from HPROF_GC_CLASS_DUMP
type StaticField struct {
	Name  string
	Type  core.JavaType
	Value core.JavaValue
}

// InstanceField is representation
// of class'es instance field name and type,
// obtained from HPROF_GC_CLASS_DUMP
type InstanceField struct {
	Name string
	Type core.JavaType
}

// Class is linked-list of information extracted
// from HPROF_GC_CLASS_DUMP as well as the Class
// data of all its superclasses.
type Class struct {
	Name           string
	Superclass     *Class
	StaticFields   []StaticField
	InstanceFields []InstanceField
}

// NormalObject is representation of HPROF_GC_INSTANCE_DUMP.
// It has corresponding Class. The values of fields are raw
// bytes. NormalObject.GetFieldValueByName should be used
// to read particular field value.
type NormalObject struct {
	identifierSize uint32
	Class          Class
	Bytes          []byte
}

// PrimitiveArray represents HPROF_GC_PRIM_ARRAY_DUMP
// but only it's "header" part. To obtain the body, i.e.
// elements of the array, raw bytes should be read from
// underlying index.Reader.
type PrimitiveArray struct {
	Type   core.JavaType
	Length uint32
}

// FieldValue is a wrapper returns from
// NormalObject.GetFieldValueByName. Origin
// field has the name of the class (or any superclass)
// that actually defines requested field.
type FieldValue struct {
	Value  core.JavaValue
	Origin string
}

// ObjectArray holds type and elements
// of object array
type ObjectArray struct {
	ArrayClassId core.Identifier
	Elements     []core.Identifier
}
