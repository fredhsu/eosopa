package eosparse

import (
	"bufio"
	"log"
	"os"
	"testing"
)

func TestLoggingHost(t *testing.T) {
	file, err := os.Open("../configs/logginghost.config")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	d := EOSDevice{}

	loggingHost := "10.90.226.100"
	scanner.Scan()
	d = parseLogging(d, scanner)
	if d.Logging.Host.Hostname != loggingHost {
		t.Fatalf("Expected logging host to be %s, but got %s", loggingHost, d.Logging.Host.Hostname)
	}
}
