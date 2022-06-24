package java

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/danielleontiev/neojhat/core"
	"github.com/danielleontiev/neojhat/dump"
	"github.com/danielleontiev/neojhat/storage"
)

// // Main.java
// public class Main {
//     public static void main(String[] args) {
//         MyClass myClass = new MyClass(42);
//         try {
//             Thread.sleep(1000000);
//         } catch (Exception e) {};
//         System.out.println(myClass);
//     }
// }
//
// // MyClass.java
// public class MyClass extends SuperClass {
//     private int my;
//
//     public MyClass(int my) {
//         super(my + 1);
//         this.my = my;
//     }
// }
//
// // SuperClass.java
// public class SuperClass extends SuperSuperClass {
//     private int superClass;
//
//     public SuperClass(int superClass) {
//         super(superClass + 1);
//         this.superClass = superClass;
//     }
// }
//
// // SuperSuperClass.java
// public class SuperSuperClass {
//     private int superSuper;
//
//     public SuperSuperClass(int superSuper) {
//         this.superSuper = superSuper;
//     }
// }

var (
	objectReaderTestFileHeader = []byte{
		0x4a, 0x41, 0x56, 0x41, 0x20, 0x50, 0x52, 0x4f, 0x46, 0x49, 0x4c, 0x45, 0x20, 0x31, 0x2e, 0x30, 0x2e, 0x32, 0x00, // header, 0-terminated
		0x00, 0x00, 0x00, 0x08, // identifier size
		0x00, 0x00, 0x01, 0x7b, // timestamp, low word
		0x7f, 0x28, 0xa8, 0x27, // timestamp, high word
	}
)

var (
	myClassString = []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, // id
		0x4d, 0x79, 0x43, 0x6c, 0x61, 0x73, 0x73, // MyClass
	}
	superClassString = []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, // id
		0x53, 0x75, 0x70, 0x65, 0x72, 0x43, 0x6c, 0x61, 0x73, 0x73, // SuperClass
	}
	superSuperClassString = []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, // id
		0x53, 0x75, 0x70, 0x65, 0x72, 0x53, 0x75, 0x70, 0x65, 0x72, 0x43, 0x6c, 0x61, 0x73, 0x73, // SuperSuperClass
	}
	myFieldString = []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, // id
		0x6d, 0x79, // my
	}
	superClassFieldString = []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, // id
		0x73, 0x75, 0x70, 0x65, 0x72, 0x43, 0x6c, 0x61, 0x73, 0x73, // superClass
	}
	superSuperFieldString = []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x06, // id
		0x73, 0x75, 0x70, 0x65, 0x72, 0x53, 0x75, 0x70, 0x65, 0x72, // superSuper
	}
	javaLangObject = []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x07, // id
		0x6a, 0x61, 0x76, 0x61, 0x2f, 0x6c, 0x61, 0x6e, 0x67, 0x2f, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, // java/lang/Object
	}
	resolvedReferencesString = []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08, //id
		0x3c, 0x72, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x64, 0x5f, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x3e, // <resolved_references>
	}
)

var (
	loadMyClass = []byte{
		0x00, 0x00, 0x00, 0x01, // class serial number
		0x00, 0x00, 0x00, 0x07, 0xff, 0x67, 0xaf, 0x78, // class object id
		0x00, 0x00, 0x00, 0x01, // stack trace serial number
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, // class name id
	}
	loadSuperClass = []byte{
		0x00, 0x00, 0x00, 0x02, // class serial number
		0x00, 0x00, 0x00, 0x07, 0xff, 0x67, 0xae, 0xd0, // class object id
		0x00, 0x00, 0x00, 0x01, // stack trace serial number
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, // class name id
	}
	loadSuperSuperClass = []byte{
		0x00, 0x00, 0x00, 0x03, // class serial number
		0x00, 0x00, 0x00, 0x07, 0xff, 0x67, 0xae, 0x28, // class object id
		0x00, 0x00, 0x00, 0x01, // stack trace serial number
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, // class name id
	}
	loadJavaLangObject = []byte{
		0x00, 0x00, 0x00, 0x04, // class serial number
		0x00, 0x00, 0x00, 0x07, 0xff, 0x60, 0x07, 0x80, // class object id
		0x00, 0x00, 0x00, 0x01, // stack trace serial number
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x07, // class name id
	}
)

var (
	myClassInstance = []byte{
		// MyClass instance
		0x00, 0x00, 0x00, 0x07, 0xff, 0x67, 0xb0, 0x18, // object id (34349756440)
		0x00, 0x00, 0x00, 0x01, // stack trace serial number
		0x00, 0x00, 0x00, 0x07, 0xff, 0x67, 0xaf, 0x78, // class object id (34349756280)
		0x00, 0x00, 0x00, 0x0c, // number of bytes that follow (12)
		0x00, 0x00, 0x00, 0x2a, 0x00, 0x00,
		0x00, 0x2b, 0x00, 0x00, 0x00, 0x2c,
	}
	myClassClass = []byte{
		// MyClass class dump
		0x00, 0x00, 0x00, 0x07, 0xff, 0x67, 0xaf, 0x78, // class object id
		0x00, 0x00, 0x00, 0x01, // stack trace serial number
		0x00, 0x00, 0x00, 0x07, 0xff, 0x67, 0xae, 0xd0, // super class object id
		0x00, 0x00, 0x00, 0x07, 0xff, 0x64, 0x95, 0x48, // class loader object id
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // signers object id
		0x00, 0x00, 0x00, 0x07, 0xff, 0x67, 0xa1, 0x30, // protection domain object id
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // reserved 1
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // reserved 2
		0x00, 0x00, 0x00, 0x0c, // instance size
		0x00, 0x00, // const. pool size
		0x00, 0x00, // number of static fields
		0x00, 0x01, // number of instance fields
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, // field name id
		0x0a, // type
	}
	superClassClass = []byte{
		// SuperClass class dump
		0x00, 0x00, 0x00, 0x07, 0xff, 0x67, 0xae, 0xd0, // class object id
		0x00, 0x00, 0x00, 0x01, // stack trace serial number
		0x00, 0x00, 0x00, 0x07, 0xff, 0x67, 0xae, 0x28, // super class object id
		0x00, 0x00, 0x00, 0x07, 0xff, 0x64, 0x95, 0x48, // class loader object id
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // signers object id
		0x00, 0x00, 0x00, 0x07, 0xff, 0x67, 0xa1, 0x30, // protection domain object id
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // reserved 1
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // reserved 2
		0x00, 0x00, 0x00, 0x08, // instance size
		0x00, 0x00, // const. pool size
		0x00, 0x00, // number of static fields
		0x00, 0x01, // number of instance fields
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, // field name id
		0x0a, // type
	}
	superSuperClass = []byte{
		// SuperSuperClass class dump
		0x00, 0x00, 0x00, 0x07, 0xff, 0x67, 0xae, 0x28, // class object id
		0x00, 0x00, 0x00, 0x01, // stack trace serial number
		0x00, 0x00, 0x00, 0x07, 0xff, 0x60, 0x07, 0x80, // super class object id
		0x00, 0x00, 0x00, 0x07, 0xff, 0x64, 0x95, 0x48, // class loader object id
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // signers object id
		0x00, 0x00, 0x00, 0x07, 0xff, 0x67, 0xa1, 0x30, // protection domain object id
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // reserved 1
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // reserved 2
		0x00, 0x00, 0x00, 0x04, // instance size
		0x00, 0x00, // const. pool size
		0x00, 0x00, // number of static fields
		0x00, 0x01, // number of instance fields
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x06, // field name id
		0x0a, // type
	}
	javaLangObjectClass = []byte{
		// java.lang.Object class dump
		0x00, 0x00, 0x00, 0x07, 0xff, 0x60, 0x07, 0x80, // class object id
		0x00, 0x00, 0x00, 0x01, // stack trace serial number
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // super class object id
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // class loader object id
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // signers object id
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // protection domain object id
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // reserved 1
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // reserved 2
		0x00, 0x00, 0x00, 0x00, // instance size
		0x00, 0x00, // const. pool size
		0x00, 0x01, // number of static fields
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08, // static field name
		0x02,                                           // type (object)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, // value
		0x00, 0x00, // number of instance filed
	}
)

var (
	objectReaderTestInput = concat(
		objectReaderTestFileHeader,
		createRecordHeader(core.HprofUtf8Tag, uint32(len(myClassString))), // MyClass string
		myClassString,
		createRecordHeader(core.HprofUtf8Tag, uint32(len(superClassString))), // SuperClass string
		superClassString,
		createRecordHeader(core.HprofUtf8Tag, uint32(len(superSuperClassString))), // SuperSuperClass string
		superSuperClassString,
		createRecordHeader(core.HprofUtf8Tag, uint32(len(myFieldString))), // my string
		myFieldString,
		createRecordHeader(core.HprofUtf8Tag, uint32(len(superClassFieldString))), // superClass string
		superClassFieldString,
		createRecordHeader(core.HprofUtf8Tag, uint32(len(superSuperFieldString))), // superSuper string
		superSuperFieldString,
		createRecordHeader(core.HprofUtf8Tag, uint32(len(javaLangObject))), // java/lang/Object
		javaLangObject,
		createRecordHeader(core.HprofUtf8Tag, uint32(len(resolvedReferencesString))), // resolved_references
		resolvedReferencesString,
		createRecordHeader(core.HprofLoadClassTag, 24), // load MyClass
		loadMyClass,
		createRecordHeader(core.HprofLoadClassTag, 24), // load SuperClass
		loadSuperClass,
		createRecordHeader(core.HprofLoadClassTag, 24), // load SuperSuperClass
		loadSuperSuperClass,
		createRecordHeader(core.HprofLoadClassTag, 24), // load java.lang.Object
		loadJavaLangObject,
		createRecordHeader(core.HprofHeapDumpSegmentTag, 0), // remaining for segments is obsolete
		createSubRecordHeader(core.HprofGcClassDumpType),    // MyClass class dump
		myClassClass,
		createSubRecordHeader(core.HprofGcClassDumpType), // SuperClass class dump
		superClassClass,
		createSubRecordHeader(core.HprofGcClassDumpType), // SuperSuperClass class dump
		superSuperClass,
		createSubRecordHeader(core.HprofGcClassDumpType), // java.lang.Object class dump
		javaLangObjectClass,
		createSubRecordHeader(core.HprofGcInstanceDumpType), // MyClass instance
		myClassInstance,
		createRecordHeader(core.HprofHeapDumpEndTag, 0),
	)
)

func TestObjectReader_ParseNormalObject(t *testing.T) {
	objectReader := createObjectReader(objectReaderTestInput, t)
	objId := 34349756440
	got, err := objectReader.ParseNormalObject(core.Identifier(objId))
	want := NormalObject{
		identifierSize: 8,
		Class: Class{
			Name: "MyClass",
			Superclass: &Class{
				Name: "SuperClass",
				Superclass: &Class{
					Name: "SuperSuperClass",
					Superclass: &Class{
						Name:       "java/lang/Object",
						Superclass: nil,
						StaticFields: []StaticField{
							{
								Name:  "<resolved_references>",
								Type:  core.Object,
								Value: core.JavaValue{Type: core.Object, Value: core.Identifier(1)},
							},
						},
					},
					InstanceFields: []InstanceField{
						{Name: "superSuper", Type: core.Int},
					},
				},
				InstanceFields: []InstanceField{
					{Name: "superClass", Type: core.Int},
				},
			},
			InstanceFields: []InstanceField{
				{Name: "my", Type: core.Int},
			},
		},
		Bytes: []byte{
			0x00, 0x00, 0x00, 0x2a, 0x00, 0x00,
			0x00, 0x2b, 0x00, 0x00, 0x00, 0x2c,
		},
	}
	if err != nil {
		t.Errorf("ObjectReader.ParseNormalObject() error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Instances = %v, want %v", got, want)
	}
}

func TestObjectReader_ParseClass(t *testing.T) {
	classId := 34349254528
	want := Class{
		Name:       "java/lang/Object",
		Superclass: nil,
		StaticFields: []StaticField{
			{
				Name:  "<resolved_references>",
				Type:  core.Object,
				Value: core.JavaValue{Type: core.Object, Value: core.Identifier(1)},
			},
		},
	}
	objectReader := createObjectReader(objectReaderTestInput, t)
	got, err := objectReader.ParseClass(core.Identifier(classId))
	if err != nil {
		t.Errorf("ParseClass error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParseClass = %v, want %v", got, want)
	}
}

func TestNormalObject_GetFieldValueByName(t *testing.T) {
	objectReader := createObjectReader(objectReaderTestInput, t)
	objId := 34349756440
	got, err := objectReader.ParseNormalObject(core.Identifier(objId))
	if err != nil {
		t.Errorf("ParseInstance error = %v", err)
	}

	my, err := got.GetFieldValueByName("my")
	wantMy := FieldValue{Value: core.JavaValue{Type: core.Int, Value: int32(42)}, Origin: "MyClass"}
	if err != nil {
		t.Errorf("GetInstanceValueByName error = %v", err)
	}

	if my != wantMy {
		t.Errorf("GetInstanceValueByName = %v, want %v", my, wantMy)
	}

	superClass, err := got.GetFieldValueByName("superClass")
	wantSuperClass := FieldValue{Value: core.JavaValue{Type: core.Int, Value: int32(43)}, Origin: "SuperClass"}
	if err != nil {
		t.Errorf("GetInstanceValueByName error = %v", err)
	}
	if superClass != wantSuperClass {
		t.Errorf("GetInstanceValueByName = %v, want %v", superClass, wantSuperClass)
	}

	superSuper, err := got.GetFieldValueByName("superSuper")
	wantSuperSuper := FieldValue{Value: core.JavaValue{Type: core.Int, Value: int32(44)}, Origin: "SuperSuperClass"}
	if err != nil {
		t.Errorf("GetInstanceValueByName error = %v", err)
	}
	if superSuper != wantSuperSuper {
		t.Errorf("GetInstanceValueByName = %v, want %v", superSuper, wantSuperSuper)
	}
}

var (
	stringInstance = []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, // obj id
		0x00, 0x00, 0x00, 0x01, // stack strace
		0x00, 0x00, 0x00, 0x07, 0xff, 0xa0, 0x0a, 0x90, // class id
		0x00, 0x00, 0x00, 0x0d, // number of bytes that follow
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
	}
	stringClass = []byte{
		0x00, 0x00, 0x00, 0x07, 0xff, 0xa0, 0x0a, 0x90, // class id
		0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x07, 0xff, 0x60, 0x07, 0x80,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x0d,
		0x00, 0x00,
		0x00, 0x07, // static fields
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x08,
		0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x08,
		0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x02,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x02,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x04,
		0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x0b,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x02,
		0x00, 0x00, 0x00, 0x07, 0xff, 0xac, 0x17, 0xe0,
		0x00, 0x03, // instance fields
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x0A,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x08,
		0x00, 0x00, 0x7f, 0xad, 0x5c, 0xe8, 0x28, 0x48,
		0x02,
	}
)

var (
	dummyString = []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, // id
		0x6a, 0x61, 0x76, 0x61, // java
	}
	valueString = []byte{
		0x00, 0x00, 0x7f, 0xad, 0x5c, 0xe8, 0x28, 0x48, // id
		0x76, 0x61, 0x6c, 0x75, 0x65, // value field name
	}
	javaLangString = []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, // id
		0x6a, 0x61, 0x76, 0x61, 0x2f, 0x6c, 0x61, 0x6e, 0x67, 0x2f, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, // java/lang/String
	}
)

var (
	loadJavaLangString = []byte{
		0x00, 0x00, 0x00, 0x01, // class serial number
		0x00, 0x00, 0x00, 0x07, 0xff, 0xa0, 0x0a, 0x90, // class object id
		0x00, 0x00, 0x00, 0x01, // stack trace serial number
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, // class name id
	}
)

var (
	stringArray = []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, // array object id
		0x00, 0x00, 0x00, 0x01, // stack trace serial number
		0x00, 0x00, 0x00, 0x04, // number of elements
		0x08,                   // type
		0x6a, 0x61, 0x76, 0x61, // java
	}
)

var (
	stringSample = concat(
		objectReaderTestFileHeader,
		createRecordHeader(core.HprofUtf8Tag, uint32(len(dummyString))), // java
		dummyString,
		createRecordHeader(core.HprofUtf8Tag, uint32(len(valueString))), // value
		valueString,
		createRecordHeader(core.HprofUtf8Tag, uint32(len(javaLangString))), // java/lang/String
		javaLangString,
		createRecordHeader(core.HprofUtf8Tag, uint32(len(javaLangObject))), // java/lang/Object
		javaLangObject,
		createRecordHeader(core.HprofUtf8Tag, uint32(len(resolvedReferencesString))), // resolved_references
		resolvedReferencesString,
		createRecordHeader(core.HprofLoadClassTag, 24),
		loadJavaLangString,
		createRecordHeader(core.HprofLoadClassTag, 24),
		loadJavaLangObject,
		createRecordHeader(core.HprofHeapDumpSegmentTag, 0),
		createSubRecordHeader(core.HprofGcInstanceDumpType),
		stringInstance,
		createSubRecordHeader(core.HprofGcClassDumpType),
		stringClass,
		createSubRecordHeader(core.HprofGcClassDumpType),
		javaLangObjectClass,
		createSubRecordHeader(core.HprofGcPrimArrayDumpType),
		stringArray,
		createRecordHeader(core.HprofHeapDumpEndTag, 0),
	)
)

func TestObjectReader_ParseJavaString(t *testing.T) {
	objectReader := createObjectReader(stringSample, t)
	str := core.JavaValue{Type: core.Object, Value: core.Identifier(1)}
	want := "java"
	got, err := objectReader.ParseJavaString(str)
	if err != nil {
		t.Errorf("ObjectReader.ParseJavaString() error = %v", err)
	}
	if got != want {
		t.Errorf("ObjectReader.ParseJavaString() = %v, want %v", got, want)
	}
}

var (
	sampleObjectArray = concat(
		objectReaderTestFileHeader,
		createRecordHeader(core.HprofHeapDumpSegmentTag, 0),
		createSubRecordHeader(core.HprofGcObjArrayDumpType),
		[]byte{
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, // array object id
			0x00, 0x00, 0x00, 0x01, // stack trace serial number
			0x00, 0x00, 0x00, 0x01, // number of elements
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, // array class id
			// elements
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		},
		createRecordHeader(core.HprofHeapDumpEndTag, 0),
	)
)

func TestObjectReader_ParseObjectArrayFull(t *testing.T) {
	objectReader := createObjectReader(sampleObjectArray, t)
	objectArrayId := core.Identifier(1)
	want := ObjectArray{
		ArrayClassId: core.Identifier(1),
		Elements: []core.Identifier{
			core.Identifier(1),
		},
	}
	got, err := objectReader.ParseObjectArrayFull(objectArrayId)
	if err != nil {
		t.Errorf("ObjectReader.ParseObjectArrayFull() error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ObjectReader.ParseObjectArrayFull() = %v, want %v", got, want)
	}
}

func concat(head []byte, tail ...[]byte) []byte {
	res := head[:]
	for _, b := range tail {
		res = append(res, b...)
	}
	return res
}

func createObjectReader(in []byte, t *testing.T) *Heap {
	return NewHeap(createReader(in, t))
}

func createReader(in []byte, t *testing.T) *dump.ParsedAccessor {
	heapDump := bytes.NewReader(in)
	smallWriter := storage.NewSmallRecordsWriteStorage()
	instanceDumpWriteVolume := storage.NewRamWriteVolume()
	objArrayDumpWriteVolume := storage.NewRamWriteVolume()
	primArrayDumpWriteVolume := storage.NewRamWriteVolume()
	bigWriter := storage.NewBigRecordsWriteStorage(
		instanceDumpWriteVolume, objArrayDumpWriteVolume, primArrayDumpWriteVolume)
	metaWriter := storage.NewMetaWriteStorage()
	parser := dump.NewParser(heapDump, smallWriter, bigWriter, metaWriter)

	if err := parser.ParseHeapDump(); err != nil {
		t.Errorf("error indexing sample input: %v", err)
	}
	bigReader, err := storage.NewBigRecordsReadStorage(
		storage.NewRamReadVolume(instanceDumpWriteVolume.Bytes()), instanceDumpWriteVolume.Len(),
		storage.NewRamReadVolume(objArrayDumpWriteVolume.Bytes()), objArrayDumpWriteVolume.Len(),
		storage.NewRamReadVolume(primArrayDumpWriteVolume.Bytes()), primArrayDumpWriteVolume.Len(),
	)
	if err != nil {
		t.Errorf("error creating bigreader: %v", err)
	}
	smallReader := storage.NewSmallRecordsReadStorage()
	metaReader := storage.NewMetaReadStorage()
	smallBuf := bytes.NewBuffer(nil)
	metaBuf := bytes.NewBuffer(nil)
	if err := smallWriter.SerializeTo(smallBuf); err != nil {
		t.Errorf("error serializing smallwriter: %v", err)
	}
	if err := metaWriter.SerializeTo(metaBuf); err != nil {
		t.Errorf("error serializing meta writer: %v", err)
	}
	if err := smallReader.RestoreFrom(smallBuf); err != nil {
		t.Errorf("error creating smallreader: %v", err)
	}
	if err := metaReader.RestoreFrom(metaBuf); err != nil {
		t.Errorf("error creating meta reader: %v", err)
	}
	reader := dump.NewParsedAccessor(heapDump, bigReader, smallReader, metaReader)
	return reader
}

// record header length = 9
func createRecordHeader(tag core.Tag, remaining uint32) []byte {
	start := []byte{
		byte(tag),              // record tag
		0x00, 0x00, 0x00, 0x00, // time, always zero
	}
	end := make([]byte, 4) // remaining bytes in record
	binary.BigEndian.PutUint32(end, remaining)
	return concat(start, end)
}

func createSubRecordHeader(ty core.SubRecordType) []byte {
	return []byte{byte(ty)}
}
