/*
	summary extracts common information from heap dump
	such as system properties of the system that was
	running application.
*/
package summary

import (
	"fmt"
	"sort"
	"time"

	"github.com/danielleontiev/neojhat/core"
	"github.com/danielleontiev/neojhat/dump"
	"github.com/danielleontiev/neojhat/java"
)

const (
	javaLangSystemClassName          = "java/lang/System"
	managementFactoryHelperClassName = "sun/management/ManagementFactoryHelper"
)

type Properties = map[string]string
type Kv struct {
	Key string
	Val string
}
type EnvProperties struct {
	System       string
	Architecture string
	JavaHome     string
	JavaVersion  string
	JavaName     string
	JavaVendor   string
}
type HeapProperties struct {
	Classes   int
	GcRoots   int
	Instances int
	HeapSize  int
}
type SystemProperties struct {
	JvmUptime string
}
type Summary struct {
	Env        EnvProperties
	Heap       HeapProperties
	System     SystemProperties
	Properties []Kv
}

// GetSummary parses given .hprof file and extracts Summary from
// it. Summary consists of short heap info, important env info and
// all system properties from java.lang.System class
// (private static java.util.Properties props).
// "props" are set by JVM on startup and contains some info about host
// system and host JVM. Take a look at java.util.Properties javadoc for
// more info. java.util.Properties contains ConcurrentHashMap<String, String>
// which effectiveley the storage with the info. GetSummary takes an attempt
// to extract the info from hash map by inspecting underlying array with nodes
// of hash map.
func GetSummary(parsedAccessor *dump.ParsedAccessor, allProps bool) (Summary, error) {
	properties, err := getAllProps(parsedAccessor)
	if err != nil {
		return Summary{}, err
	}
	env := getEnv(properties)
	heap, err := getHeap(parsedAccessor)
	if err != nil {
		return Summary{}, err
	}
	system, err := getSystem(parsedAccessor)
	if err != nil {
		return Summary{}, err
	}
	if allProps {
		return Summary{
			Env:        env,
			Heap:       heap,
			System:     system,
			Properties: makeSortedList(properties),
		}, nil
	}
	return Summary{
		Env:        env,
		Heap:       heap,
		System:     system,
		Properties: nil,
	}, nil
}

// GetAllProps is similar to GetSummary but returns all properties from
// java.lang.System class
func getAllProps(parsedAccessor *dump.ParsedAccessor) (Properties, error) {
	heap := java.NewHeap(parsedAccessor)

	javaLangSystemClass, err := getClassByName(parsedAccessor, javaLangSystemClassName)
	if err != nil {
		return nil, err
	}
	var propsObjectId core.Identifier
	for _, field := range javaLangSystemClass.StaticFields {
		if field.Name == "props" {
			id, err := field.Value.ToObject()
			if err != nil {
				return nil, err
			}
			propsObjectId = id
			break
		}
	}
	props, err := heap.ParseNormalObject(propsObjectId)
	if err != nil {
		return nil, err
	}
	propsMapField, err := props.GetFieldValueByName("map")
	if err != nil {
		return nil, err
	}
	propsMapId, err := propsMapField.Value.ToObject()
	if err != nil {
		return nil, err
	}
	propsMap, err := heap.ParseNormalObject(propsMapId)
	if err != nil {
		return nil, err
	}
	tableField, err := propsMap.GetFieldValueByName("table")
	if err != nil {
		return nil, err
	}
	tableId, err := tableField.Value.ToObject()
	if err != nil {
		return nil, err
	}
	table, err := heap.ParseObjectArrayFull(tableId)
	if err != nil {
		return nil, err
	}
	resultProps := make(map[string]string)
	for _, nodeObjectId := range table.Elements {
		if nodeObjectId != 0 {
			node, err := heap.ParseNormalObject(nodeObjectId)
			if err != nil {
				return nil, err
			}
			key, err := node.GetFieldValueByName("key")
			if err != nil {
				return nil, err
			}
			val, err := node.GetFieldValueByName("val")
			if err != nil {
				return nil, err
			}
			keyString, err := heap.ParseJavaString(key.Value)
			if err != nil {
				return nil, err
			}
			valString, err := heap.ParseJavaString(val.Value)
			if err != nil {
				return nil, err
			}
			resultProps[keyString] = valString
		}
	}
	return Properties(resultProps), nil
}

func getEnv(props Properties) EnvProperties {
	return EnvProperties{
		System:       props["os.name"],
		Architecture: props["os.arch"],
		JavaHome:     props["java.home"],
		JavaVersion:  props["java.version"],
		JavaName: fmt.Sprintf("%s (%s, %s)",
			props["java.vm.name"],
			props["java.vm.version"],
			props["java.vm.info"]),
		JavaVendor: props["java.vm.vendor"],
	}
}

func getHeap(parsedAccessor *dump.ParsedAccessor) (HeapProperties, error) {
	classes := parsedAccessor.ListHprofLoadClass()
	gcRootsCount := len(parsedAccessor.HprofGcRootJavaFrame) + len(parsedAccessor.HprofGcRootJniGlobal) +
		len(parsedAccessor.HprofGcRootJniLocal) + len(parsedAccessor.ListHprofGcRootStickyClass()) +
		len(parsedAccessor.ListHprofGcRootThreadObj())
	classSet := make(map[core.Identifier]any)
	var void any
	for _, c := range classes {
		classSet[c.ClassObjectId] = void
	}
	meta := parsedAccessor.MetaStorage
	var totalSize int
	sizeInfo := core.NewSizeInfo(parsedAccessor.IdSize)
	for tpe, num := range meta.Counters.PrimArrayElementsCount {
		totalSize += sizeInfo.OfType(tpe) * num
	}
	for _, num := range meta.Counters.ObjArrayElementsCount {
		totalSize += sizeInfo.OfType(core.Object) * num
	}
	for classId, num := range meta.Counters.InstancesCount {
		class, err := parsedAccessor.GetHprofGcClassDump(classId)
		if err != nil {
			return HeapProperties{}, err
		}
		totalSize += int(class.InstanceSize) * num
	}
	var totalCount int
	for _, num := range meta.Counters.PrimArraysCount {
		totalCount += num
	}
	for _, num := range meta.Counters.ObjArraysCount {
		totalCount += num
	}
	for _, num := range meta.Counters.InstancesCount {
		totalCount += num
	}
	return HeapProperties{
		Classes:   len(classSet),
		GcRoots:   gcRootsCount,
		Instances: totalCount,
		HeapSize:  totalSize,
	}, nil
}

func getSystem(parsedAccessor *dump.ParsedAccessor) (SystemProperties, error) {
	jvmStartupTime, err := getJvmStartTime(parsedAccessor)
	if err != nil {
		return SystemProperties{}, err
	}
	uptime := parsedAccessor.Timestamp.Sub(jvmStartupTime)
	return SystemProperties{
		JvmUptime: uptime.String(),
	}, nil
}

func getJvmStartTime(parsedAccessor *dump.ParsedAccessor) (time.Time, error) {
	jmxFactoryHelper, err := getClassByName(parsedAccessor, managementFactoryHelperClassName)
	if err != nil {
		return time.Time{}, err
	}
	var runtimeMBean core.Identifier
	for _, field := range jmxFactoryHelper.StaticFields {
		if field.Name == "runtimeMBean" {
			id, err := field.Value.ToObject()
			if err != nil {
				return time.Time{}, err
			}
			runtimeMBean = id
			break
		}
	}
	heap := java.NewHeap(parsedAccessor)
	instance, err := heap.ParseNormalObject(runtimeMBean)
	if err != nil {
		return time.Time{}, err
	}
	vmStartupTime, err := instance.GetFieldValueByName("vmStartupTime")
	if err != nil {
		return time.Time{}, err
	}
	timeLong, err := vmStartupTime.Value.ToLong()
	if err != nil {
		return time.Time{}, err
	}
	return time.UnixMilli(int64(timeLong)), nil
}

func getClassByName(parsedAccessor *dump.ParsedAccessor, className string) (java.Class, error) {
	heap := java.NewHeap(parsedAccessor)
	loadClasses := parsedAccessor.ListHprofLoadClass()
	var targetLoadClass core.HprofLoadClass
	for _, loadClass := range loadClasses {
		utf8, err := parsedAccessor.GetHprofUtf8(loadClass.ClassNameId)
		if err != nil {
			return java.Class{}, err
		}
		if utf8.Characters == className {
			targetLoadClass = loadClass
			break
		}
	}
	class, err := heap.ParseClass(targetLoadClass.ClassObjectId)
	return class, err
}

func makeSortedList(m map[string]string) []Kv {
	var list []Kv
	for k, v := range m {
		list = append(list, Kv{k, v})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Key < list[j].Key })
	return list
}
