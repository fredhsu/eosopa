package eosparse

import (
	"bufio"
	"log"
	"strings"
)

// ManagementType restricts the possible types of management interfaces
type ManagementType int

// Enum of different possible management type
const (
	TELNET ManagementType = iota
	SSH
	HTTPAPI
	GNMI
	GRIBI
	NETCONF
	OPENCONFIG
	RESTCONF
)

func (mt ManagementType) String() string {
	switch mt {
	case TELNET:
		return "telnet"
	case SSH:
		return "ssh"
	case HTTPAPI:
		return "httpapi"
	case GNMI:
		return "gnmi"
	case GRIBI:
		return "gribi"
	case NETCONF:
		return "netconf"
	case OPENCONFIG:
		return "openconfig"
	case RESTCONF:
		return "restconf"
	default:
		return "unknown"
	}
}

// Management is an umbrella type for json encoding -- UNUSED
type Management struct {
	Management map[string]interface{} `json:"management"`
}

// ManagementIntf gives common management functions
type ManagementIntf interface {
	Type() string
	Enabled() bool
}

// ManagementAPIHTTP stores data about eAPI config
type ManagementAPIHTTP struct {
	TypeID    string
	protocols []string
	Shutdown  bool `json:"shutdown"`
}

// Enabled indicates if the feature is enabled
func (m ManagementAPIHTTP) Enabled() bool {
	return !m.Shutdown
}

// Type gives a string value of the management type
func (m ManagementAPIHTTP) Type() string {
	return m.TypeID
}

// ManagementTelnet stores telnet settings
type ManagementTelnet struct {
	IdleTimeout  int    `json:"idleTimeout"`
	IP           string `json:"ip"`   //TODO: create struct
	IPv6         string `json:"ipv6"` //TODO: create struct
	SessionLimit int    `json:"sessionLimit"`
	VRF          string `json:"vrf"` //TODO: create struct
	TypeID       string `json:"typeId"`
	Shutdown     bool   `json:"shutdown"`
}

// Type gives a string value of the management type
func (m ManagementTelnet) Type() string {
	return m.TypeID
}

// Enabled indicates if the feature is enabled
func (m ManagementTelnet) Enabled() bool {
	return !m.Shutdown
}

// ManagementSSH stores ssh settings
type ManagementSSH struct {
	TypeID     string
	Shutdown   bool `json:"shutdown"`
	ServerPort int
}

// Enabled indicates if the feature is enabled
func (m ManagementSSH) Enabled() bool {
	return !m.Shutdown
}

// Type gives a string value of the management type
func (m ManagementSSH) Type() string {
	return m.TypeID
}

// NewManagement Creates an empty mangagement block with default values
func NewManagement() map[string]interface{} {
	// The following are defaults
	m := map[string]interface{}{
		"ssh":    ManagementSSH{TypeID: "ssh", Shutdown: false},
		"telnet": ManagementTelnet{TypeID: "telnet", Shutdown: true},
		"api":    ManagementAPIHTTP{TypeID: "http", Shutdown: true},
	}
	return m
}

// ParseManagement takes an EOSDevice (assumed management has been created) and scanner
// then returns an EOSDevice with the parsed management information added
func ParseManagement(d EOSDevice, scanner *bufio.Scanner) EOSDevice {
	line := strings.Fields(scanner.Text())
	switch mgmt := line[1]; mgmt {
	case "api":
		d.Management[mgmt] = ManagementAPIHTTP{TypeID: "api"}
	case "telnet":
		d.Management[mgmt] = parseTelnet(scanner)
	case "ssh":
		d.Management[mgmt] = parseSSH(scanner)
	default:
		log.Printf("%s is not a recognized management type", mgmt)
		// log.Fatal("Cannot continue parsing")
	}
	return d
}

func parseSSH(scanner *bufio.Scanner) ManagementSSH {
	m := ManagementSSH{TypeID: "ssh", Shutdown: false}
	line := strings.Fields(scanner.Text())
	if Contains(line, "shutdown") {
		m.Shutdown = ParseShutdown(line)
	}
	return m
}

func parseTelnet(scanner *bufio.Scanner) ManagementTelnet {
	m := ManagementTelnet{TypeID: "telnet", Shutdown: true}
	for scanner.Scan() {
		// for line := strings.Fields(scanner.Text()); line[0] != "!"; line = strings.Fields(scanner.Text()) {
		line := strings.Fields(scanner.Text())
		if line[0] == "!" {
			break
		}
		if Contains(line, "shutdown") {
			m.Shutdown = ParseShutdown(line)
			// log.Printf("parsing shutdown line : %s and result is %+v\n", line, m.Shutdown)
		}
		// idle-timeout
		// ip access-group
		// ipv6 access-group
		// session-limit
		// vrf
	}
	return m
}
