package objects

import "fmt"

type SortBy int

const (
	Count SortBy = iota
	Size
)

func (s *SortBy) String() string {
	switch *s {
	case Count:
		return "count"
	case Size:
		return "size"
	}
	return "unknown"
}

func (s *SortBy) Set(value string) error {
	switch value {
	case "count":
		*s = Count
		return nil
	case "size":
		*s = Size
		return nil
	case "":
		*s = Count
		return nil
	}
	return fmt.Errorf("Use \"count\" or \"size\" instead")
}

type ObjectItem struct {
	Name           string
	TotalSize      int
	InstancesCount int
}

type printItem struct {
	Name           string
	TotalSize      string
	InstancesCount string
}

type Objects struct {
	Items      []ObjectItem
	TotalSize  int
	TotalCount int
	SortBy     SortBy
}
