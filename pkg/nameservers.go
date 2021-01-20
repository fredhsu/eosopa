package eosparse

import (
	"bufio"
	"strings"
)

// NameServers is an ip name-server
type NameServers struct {
	Vrf       string
	Addresses []string
}

// parses ip name-server command (max: 3)
// ex: ip name-server vrf default 172.22.22.40
func ParseNameServers(d EOSDevice, scanner *bufio.Scanner) EOSDevice {
	ns := NameServers{}
	line := strings.Fields(scanner.Text())
	addrs := []string{}
	ns.Vrf = line[3]
	for _, addr := range line[4:] {
		addrs = append(addrs, addr)
	}
	ns.Addresses = addrs
	d.IPNameServers = ns
	return d
}

// parses dns domain command
// ex: dns domain sjc.aristanetworks.com
func ParseDNSDomain(d EOSDevice, scanner *bufio.Scanner) EOSDevice {
	return d
}
