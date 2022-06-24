/*
	heap provides functions for extracting meaningfull data from heap dump.
	For example, it is able to parse class hierarchies from HPROF_GC_CLASS_DUMP
	records into the linked-list like data structure or parse values of heap
	objects like instances of instances fields.

	Assembling the such data requires a lot of references traverse by the identifiers
	of heap objects. Package dump.ParsedAccessor is used for that purpose. For example,
	to get the class from class dump it's important to follow the links and collect
	superclass, superclass of superclass, etc.
*/
package java

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/danielleontiev/neojhat/core"
	"github.com/danielleontiev/neojhat/dump"
)

// Heap contains methods for reading data
// from the heap dump.
type Heap struct {
	parsedAccessor *dump.ParsedAccessor
}

func NewHeap(parsedAccessor *dump.ParsedAccessor) *Heap {
	return &Heap{parsedAccessor: parsedAccessor}
}

// ParseNormalObject is used to parse NormalObject by given objectId.
// It also triggers ObjectReader.ParseClass because NormalObject should
// contain information about class it represents.
func (h *Heap) ParseNormalObject(objectId core.Identifier) (NormalObject, error) {
	instance, err := h.parsedAccessor.GetHprofGcInstanceDump(objectId)
	if err != nil {
		return NormalObject{}, fmt.Errorf("error reading instance with id %v: %w", objectId, err)
	}
	instanceBytes, err := h.parsedAccessor.GetBytesFromCurrent(int(instance.NumberOfBytesThatFollow))
	if err != nil {
		return NormalObject{}, fmt.Errorf("error reading instance(id = %v) payload: %w", objectId, err)
	}
	class, err := h.ParseClass(instance.ClassObjectId)
	if err != nil {
		return NormalObject{}, fmt.Errorf("error reading class for instance with id %v: %w", objectId, err)
	}
	return NormalObject{
		identifierSize: h.parsedAccessor.IdentifierSize,
		Class:          class,
		Bytes:          instanceBytes,
	}, nil
}

// ParseJavaString takes the instance values known to be a reference to
// java.lang.String object and converts it to convenient Go string.
func (h *Heap) ParseJavaString(str core.JavaValue) (string, error) {
	if str.Type == core.Object {
		id, ok := str.Value.(core.Identifier)
		if ok {
			javaString, err := h.ParseNormalObject(id)
			if err != nil {
				return "", fmt.Errorf("error parsing string object: %w", err)
			}
			if javaString.Class.Name == "java/lang/String" {
				value, err := javaString.GetFieldValueByName("value")
				if err != nil {
					return "", fmt.Errorf("error getting `value` array from java.lang.String: %w", err)
				}
				arrayObj, ok := value.Value.Value.(core.Identifier)
				if ok {
					byteArr, err := h.parseByteArrayFull(arrayObj)
					if err != nil {
						return "", fmt.Errorf("error reading `value` array for string %v: %w", str, err)
					}
					return string(byteArr), nil
				}
				return "", fmt.Errorf("malformed java string, `value` is not a reference to array (%T)", value.Value.Value)
			}
			return "", fmt.Errorf("cannot run StringValue on non java.lang.String (%v)", javaString.Class.Name)
		}
		return "", fmt.Errorf("malformed value (%v), identifier is not core.Identifier value (%T)", str.Value, str.Value)
	}
	return "", fmt.Errorf("passed value (%v) is not reference to object", str)
}

func (h *Heap) parseByteArrayFull(arrayObjectId core.Identifier) ([]byte, error) {
	header, err := h.parsedAccessor.GetHprofGcPrimArray(arrayObjectId)
	if err != nil {
		return nil, fmt.Errorf("error parsing byte array with id = %v: %w", arrayObjectId, err)
	}
	if header.ElementType != core.Byte {
		return nil, fmt.Errorf("array with id = %v is not byte array (%T)", arrayObjectId, header.ElementType)
	}
	payload, err := h.parsedAccessor.GetBytesFromCurrent(int(header.NumberOfElements))
	if err != nil {
		return nil, fmt.Errorf("error parsing payload of array with id = %v: %w", arrayObjectId, err)
	}
	return payload, nil
}

// ParseClass parses inforamsion from HPROF_DC_CLASS_DUMP record of a class and all
// its superclasses.
func (h *Heap) ParseClass(classId core.Identifier) (Class, error) {
	class, err := h.parsedAccessor.GetHprofGcClassDump(classId)
	if err != nil {
		return Class{}, fmt.Errorf("error reading class with id %v", classId)
	}
	loadClass, err := h.parsedAccessor.GetHprofLoadClassByClassObjectId(class.ClassObjectId)
	if err != nil {
		return Class{}, fmt.Errorf("error reading HprofLoadClass by id %v: %w", class.ClassObjectId, err)
	}
	className, err := h.parsedAccessor.GetHprofUtf8(loadClass.ClassNameId)
	if err != nil {
		return Class{}, fmt.Errorf("error reading class name: %w", err)
	}
	var staticFields []StaticField
	for _, sf := range class.StaticFieldRecords {
		staticFieldName, err := h.parsedAccessor.GetHprofUtf8(sf.StaticFieldName)
		if err != nil {
			return Class{}, fmt.Errorf("error reading field name with id %v: %w", sf.StaticFieldName, err)
		}
		staticField := StaticField{
			Name:  staticFieldName.Characters,
			Type:  sf.Ty,
			Value: sf.Value,
		}
		staticFields = append(staticFields, staticField)
	}
	var instanceFields []InstanceField
	for _, ins := range class.InstanceFieldRecords {
		instanceFieldName, err := h.parsedAccessor.GetHprofUtf8(ins.InstanceFieldName)
		if err != nil {
			return Class{}, fmt.Errorf("error reading field name with id %v: %w", ins.InstanceFieldName, err)
		}
		instanceField := InstanceField{
			Name: instanceFieldName.Characters,
			Type: ins.Ty,
		}
		instanceFields = append(instanceFields, instanceField)
	}

	var superClass *Class
	// java.lang.Object has superclass id = 0
	if class.SuperclassObjectId == 0 {
		superClass = nil
	} else {
		// TODO not stack safe, fix
		sc, err := h.ParseClass(class.SuperclassObjectId)
		if err != nil {
			return Class{}, fmt.Errorf("error parsing superclass: %w", err)
		}
		superClass = &sc
	}
	return Class{
		Name:           className.Characters,
		Superclass:     superClass,
		StaticFields:   staticFields,
		InstanceFields: instanceFields,
	}, nil
}

// GetFieldValueByName extracts instance field value by name of the field as
// well as name of the class that actually defines requested field.
func (o *NormalObject) GetFieldValueByName(name string) (FieldValue, error) {
	reader := bytes.NewReader(o.Bytes)
	primitiveParser := core.NewPrimitiveParser(reader, o.identifierSize)
	var ty core.JavaType
	var nameFound bool
	var offset int
	var origin string
	class := &o.Class
	for class != nil {
		origin = class.Name
		found, typ, length := findField(name, class.InstanceFields, o.identifierSize)
		offset += length
		if found {
			nameFound = true
			ty = typ
			break
		}
		class = class.Superclass
		if class == nil {
			nameFound = false
			break
		}
	}
	if !nameFound {
		return FieldValue{}, fmt.Errorf("field {%v} not found", name)
	}
	_, err := reader.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return FieldValue{}, fmt.Errorf("cannot parse value of field %v: %w", name, err)
	}
	value, err := primitiveParser.ParseJavaValue(ty)
	if err != nil {
		return FieldValue{}, fmt.Errorf("cannot parse value of field %v: %w", name, err)
	}
	return FieldValue{
		Value:  value,
		Origin: origin,
	}, nil
}

func findField(name string, fields []InstanceField, idSize uint32) (found bool, ty core.JavaType, offset int) {
	var size = core.NewSizeInfo(idSize)
	for _, field := range fields {
		if field.Name == name {
			found = true
			ty = field.Type
			return
		}
		offset += size.OfType(field.Type)
	}
	found = false
	return
}

// ParseObjectArrayFull reads the whole object array and
// returns slice with ObjectArray structs
func (h *Heap) ParseObjectArrayFull(arrayObjectId core.Identifier) (ObjectArray, error) {
	header, err := h.parsedAccessor.GetHprofGcObjArray(arrayObjectId)
	if err != nil {
		return ObjectArray{}, fmt.Errorf("error parsing object array header with id %v: %w", arrayObjectId, err)
	}
	var elements []core.Identifier
	for i := header.NumberOfElements; i > 0; i-- {
		e, err := h.parsedAccessor.GetBytesFromCurrent(int(h.parsedAccessor.IdentifierSize))
		if err != nil {
			return ObjectArray{}, fmt.Errorf("error reading element of object array with id %v: %w", arrayObjectId, err)
		}
		id := binary.BigEndian.Uint64(e)
		elements = append(elements, core.Identifier(id))
	}
	return ObjectArray{
		ArrayClassId: header.ArrayClassId,
		Elements:     elements,
	}, nil
}
