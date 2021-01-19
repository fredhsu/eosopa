package eosparse

import (
	"bufio"
	"strconv"
	"strings"
)

// SWVersion stores EOS version with integer based Major and Minor
type SWVersion struct {
	Major int    `json:"major"`
	Minor int    `json:"minor"`
	Patch string `json:"patch"`
	Meta  string `json:"meta"`
}

type HWVersion = string

// parseDeviceInfo assumes it is receiving a commented line with device information
// ex. ! device: DMZ-LF18 (DCS-7060SX2-48YC6, EOS-4.24.2.1F)
func parseDeviceInfo(d EOSDevice, scanner *bufio.Scanner) EOSDevice {
	line := strings.Fields(scanner.Text())
	d.SWVersion = parseSWVersion(line[4])
	d.HWVersion = parseHWVersion(line[3])
	return d
}

func parseSWVersion(line string) SWVersion {
	cleanLine := strings.TrimPrefix(line, "EOS-")
	cleanLine = strings.TrimSuffix(cleanLine, ")")
	fields := strings.Split(cleanLine, ".")
	var major, minor int
	if i, err := strconv.Atoi(fields[0]); err == nil {
		major = i
	}
	if i, err := strconv.Atoi(fields[1]); err == nil {
		minor = i
	}
	swv := SWVersion{Major: major, Minor: minor, Patch: fields[2]}
	if len(fields) > 3 {
		swv.Meta = fields[3]
	}
	return swv
}

func parseHWVersion(line string) string {
	cleanLine := strings.TrimPrefix(line, "(")
	cleanLine = strings.TrimSuffix(cleanLine, ",")
	return cleanLine
}
