package helpers

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
