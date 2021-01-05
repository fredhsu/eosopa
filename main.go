package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

// Hostname is used to hold the hostname for JSON serialization
type Hostname struct {
	Hostname string
}

// EOSDevices captures a list of switches
type EOSDevices struct {
	Switches []EOSDevice
}

// EOSDevice is an EOS endpoint
type EOSDevice struct {
	Hostname   string
	Management Management
}

// Management is an umbrella type for json encoding
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

// TODO : scan through multiple commented lines to get to non-commented line
func parseComments() {

}

func parseManagement(scanner *bufio.Scanner, line []string) ManagementIntf {
	switch mgmt := line[1]; mgmt {
	case "api":
		return ManagementAPIHTTP{TypeID: "api"}
	case "telnet":
		return parseTelnet(scanner)
	case "ssh":
		return parseSSH(scanner)
	default:
		log.Printf("%s is not a recognized management type", mgmt)
		return nil
	}
}

func parseSSH(scanner *bufio.Scanner) ManagementSSH {
	m := ManagementSSH{TypeID: "ssh", Shutdown: false}
	line := strings.Fields(scanner.Text())
	if contains(line, "shutdown") {
		m.Shutdown = parseShutdown(line)
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
		if contains(line, "shutdown") {
			m.Shutdown = parseShutdown(line)
			log.Printf("parsing shutdown line : %s and result is %+v\n", line, m.Shutdown)
		}
		// idle-timeout
		// ip access-group
		// ipv6 access-group
		// session-limit
		// vrf
	}
	return m
}

func parseShutdown(line []string) bool {
	if line[0] == "no" {
		return false
	}
	return true
}

func contains(xs []string, s string) bool {
	for _, x := range xs {
		if x == s {
			return true
		}
	}
	return false
}

func parseHostname(line []string) {

}
func main() {
	file, err := os.Open("./telnetenabled.config")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	v := make(map[string]interface{})
	// SSH is enabled by default
	v["ssh"] = ManagementSSH{TypeID: "ssh", Shutdown: false}

	m := Management{Management: v}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if line[0] == "management" {
			mgmt := parseManagement(scanner, line)
			if mgmt != nil {
				m.Management[mgmt.Type()] = mgmt
			}
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	b, err := json.MarshalIndent(m, " ", "  ")
	fmt.Printf("%s\n", b)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
