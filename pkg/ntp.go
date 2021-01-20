package eosparse

import (
	"bufio"
)

// parses ntp server command
// ex.
// ntp server 172.22.22.50
// ntp server 198.55.111.50
// ntp server 216.229.0.17
func ParseNTPServer(d EOSDevice, scanner *bufio.Scanner) EOSDevice {
	return d
}
