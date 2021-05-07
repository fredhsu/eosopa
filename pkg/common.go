package eosparse

import "bufio"

// ParseNoCommand will negate the config of a feature
func ParseNoCommand(d EOSDevice, scanner *bufio.Scanner) EOSDevice {
	// TODO How to best implement?
	return d
}

// Contains indicates if a string s is present in a slice of strings xs
func Contains(xs []string, s string) bool {
	for _, x := range xs {
		if x == s {
			return true
		}
	}
	return false
}

// ParseShutdown indicates if a feature is shutdown or no shutdown
func ParseShutdown(line []string) bool {
	if line[0] == "no" {
		return false
	}
	return true
}

// StringSliceEq returns true if all the strings in s1 are equal to s2
func StringSliceEq(s1, s2 []string) bool {
	for i, s := range s1 {
		if s != s2[i] {
			return false
		}
	}
	return true
}
