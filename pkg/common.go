package eosparse

func Contains(xs []string, s string) bool {
	for _, x := range xs {
		if x == s {
			return true
		}
	}
	return false
}

func ParseShutdown(line []string) bool {
	if line[0] == "no" {
		return false
	}
	return true
}
