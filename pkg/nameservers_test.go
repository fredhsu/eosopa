package eosparse

import (
	"bufio"
	"log"
	"os"
	"testing"
)

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
