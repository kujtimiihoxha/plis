package helpers

import (
	"bytes"
	"regexp"
)

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func FromStringArrayToString(list []string) string {
	s := ""
	for _, v := range list {
		s += v + string('\n')
	}
	return s
}

func ToCamelCaseOrUnderscore(src string) string {
	camelingRegex := regexp.MustCompile("[0-9A-Za-z_]+")
	byteSrc := []byte(src)
	chunks := camelingRegex.FindAll(byteSrc, -1)
	for idx, val := range chunks {
		if idx > 0 {
			chunks[idx] = bytes.Title(val)
		}
	}
	return string(bytes.Join(chunks, nil))
}
