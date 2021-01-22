package eosparse

import (
	"bufio"
	"io"
	"strings"
	"testing"
)

func TestParseDeviceInfo(t *testing.T) {
	var r io.Reader = strings.NewReader("! device: DMZ-LF18 (DCS-7060SX2-48YC6, EOS-4.24.2.1F)\n")
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	d := EOSDevice{}
	d = ParseDeviceInfo(d, scanner)
	swver := SWVersion{
		Major: 4,
		Minor: 24,
		Patch: "2",
		Meta:  "1F",
	}
	hwver := "DCS-7060SX2-48YC6"
	if d.SWVersion != swver {
		t.Fatalf("Expected %+v :: got %+v", swver, d.SWVersion)
	}
	if d.HWVersion != hwver {
		t.Fatalf("Expected %s :: got %s", hwver, d.HWVersion)
	}

}

func TestParseHostname(t *testing.T) {
	d := EOSDevice{}
	line := []string{"hostname", "foo"}
	d = ParseHostname(d, line)
	if d.Hostname != "foo" {
		t.Fatalf("Expected hostname to be %s, but got %s", line[1], d.Hostname)
	}
}
