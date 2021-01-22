package eosparse

import (
	"bufio"
	"log"
	"os"
	"testing"
)

func TestParseTelnet(t *testing.T) {
	// Test for telnet no shutdown
	enabledFile, err := os.Open("../telnetenabled.config")
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
	shutFile, err := os.Open("../telnetshutdown.config")
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

	defaultFile, err := os.Open("../emptyconfig.config")
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
