package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	eosparse "github.com/fredhsu/eosopa/pkg"
)

// EOSDevices captures a list of switches
type EOSDevices struct {
	Switches []eosparse.EOSDevice `json:"switches"`
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

func parseManagement(scanner *bufio.Scanner, line []string) ManagementIntf {
	// m := Management{}
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

/*
ip name-server vrf default 172.22.22.40
dns domain sjc.aristanetworks.com
!
ntp server 172.22.22.50
ntp server 198.55.111.50
ntp server 216.229.0.17
*/

// createManagement() Creates an empty mangagement block with default values
func createManagement() map[string]interface{} {
	m := map[string]interface{}{
		"ssh":    ManagementSSH{TypeID: "ssh", Shutdown: false},
		"telnet": ManagementTelnet{TypeID: "telnet", Shutdown: true},
		"api":    ManagementAPIHTTP{TypeID: "http", Shutdown: true},
	}
	return m
}

func main() {
	filePtr := flag.String("input", "eos.config", "config file to convert")
	flag.Parse()

	file, err := os.Open(*filePtr)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// v := make(map[string]interface{})
	// SSH is enabled by default
	// v["ssh"] = ManagementSSH{TypeID: "ssh", Shutdown: false}

	devices := EOSDevices{}
	device := eosparse.NewEOSDevice()
	// m := Management{Management: v}
	// var m = map[string]interface{}
	// m := map[string]interface{}{}
	m := createManagement()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		// TODO: How to handle "no" prefix? contains? -- will this only take place in subconfig?
		switch line[0] {
		case "!":
			// Skipping comments unless the provide device info
			{
				if len(line) > 1 && line[1] == "device:" {
					device = eosparse.ParseDeviceInfo(device, scanner)
				} else {
					continue
				}
			}
		case "management":
			{
				mgmt := parseManagement(scanner, line)
				// device := parseManagement(device, scanner)
				if mgmt != nil {
					m[mgmt.Type()] = mgmt
				}
				if err != nil {
					log.Fatal(err)
				}
			}
		case "hostname":
			{
				device = eosparse.ParseHostname(device, line)
			}
		}
	}
	device.Management = m
	devices.Switches = append(devices.Switches, device)
	d, err := json.MarshalIndent(devices, " ", "  ")
	fmt.Printf("%s\n", d)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
