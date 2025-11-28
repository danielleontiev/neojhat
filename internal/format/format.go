package format

import (
	"fmt"
	"strings"
)

func ClassName(className string) string {
	return strings.ReplaceAll(className, "/", ".")
}

func Size(bytes int) string {
	k := 1_024
	m := k * k
	g := k * k * k
	if bytes < k {
		return fmt.Sprintf("%vB", bytes)
	}
	if bytes < m {
		return fmt.Sprintf("%vK", bytes/k)
	}
	if bytes < g {
		return fmt.Sprintf("%vM", bytes/m)
	}
	return fmt.Sprintf("%vG", bytes/g)
}

func Signature(signature string) (string, string) {
	var inClassName bool
	var dimCount int
	var classNameBuffer string
	var result []string

	var argLen int
	activeLen := &argLen
	for _, s := range signature {
		switch s {
		case '(':
			continue
		case ')':
			activeLen = new(int)
			continue
		case '[':
			dimCount++
			continue
		case 'L':
			if !inClassName {
				inClassName = true
				continue
			}
		case ';':
			inClassName = false
			result = append(result, classNameBuffer+strings.Repeat("[]", dimCount))
			*activeLen++
			classNameBuffer = ""
			dimCount = 0
			continue
		}
		if inClassName {
			if s == '/' {
				classNameBuffer += "."
			} else {
				classNameBuffer += string(s)
			}

		} else {
			var token string
			switch s {
			case 'B':
				token = "byte"
			case 'C':
				token = "char"
			case 'D':
				token = "double"
			case 'F':
				token = "float"
			case 'I':
				token = "int"
			case 'J':
				token = "long"
			case 'S':
				token = "short"
			case 'Z':
				token = "boolean"
			case 'V':
				token = "void"
			}
			result = append(result, token+strings.Repeat("[]", dimCount))
			*activeLen++
			dimCount = 0
		}
	}
	args := strings.Join(result[:argLen], ", ")
	if len(result) == argLen {
		return args, ""
	}
	return args, result[len(result)-1]
}
