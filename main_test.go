package main

import (
	"bufio"
	"log"
	"os"
	"testing"
)

func TestParseHostname(t *testing.T) {
	d := EOSDevice{}
	line := []string{"hostname", "foo"}
	d = parseHostname(d, line)
	if d.Hostname != "foo" {
		t.Fatalf("Expected hostname to be %s, but got %s", line[1], d.Hostname)
	}
}

func TestParseTelnet(t *testing.T) {
	// Test for telnet no shutdown
	enabledFile, err := os.Open("./telnetenabled.config")
	if err != nil {
		log.Fatal(err)
	}
	defer enabledFile.Close()

	scanner := bufio.NewScanner(enabledFile)
	mt := parseTelnet(scanner)
	if mt.TypeID != "telnet" {
		t.Fatalf(`Wrong typeID parsed`)
	}

	// mt.Shutdown means it is shutdown
	if mt.Shutdown {
		t.Fatalf("Did not parse telnet no shutdown correctly %+v", mt)
	}

	// Test for shutdown
	shutFile, err := os.Open("./telnetshutdown.config")
	if err != nil {
		log.Fatal(err)
	}
	defer shutFile.Close()

	scanner = bufio.NewScanner(shutFile)
	mt = parseTelnet(scanner)

	// Fails if not shutdown
	if !mt.Shutdown {
		t.Fatalf("Did not parse telnet shutdown correctly")
	}

	// Test for default

	defaultFile, err := os.Open("./emptyconfig.config")
	if err != nil {
		log.Fatal(err)
	}
	defer defaultFile.Close()
	scanner = bufio.NewScanner(defaultFile)
	mt = parseTelnet(scanner)

	// Fails if not shutdown
	if !mt.Shutdown {
		t.Fatalf("Did not parse default value of shutdown (shut) correctly")
	}
}

func TestParseNameServers(t *testing.T) {
	file, err := os.Open("./nameserver.config")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	d := EOSDevice{}
	vrf := "default"
	addresses := []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"}
	scanner.Scan()
	d = parseNameServers(d, scanner)
	if d.IPNameServers.Vrf != vrf {
		t.Fatalf("Expected vrf to be %s, but got %s", d.IPNameServers.Vrf, vrf)
	}
	if !stringSliceEq(d.IPNameServers.Addresses, addresses) {
		t.Fatalf("Expected addresses to be %+v, but got %+v", d.IPNameServers.Addresses, addresses)
	}
}

func stringSliceEq(a, b []string) bool {
	for i, x := range a {
		if x != b[i] {
			return false
		}
	}
	return true
}
