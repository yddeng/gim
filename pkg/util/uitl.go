package util

func HasString(s string, m []string) bool {
	for _, v := range m {
		if v == s {
			return true
		}
	}
	return false
}
