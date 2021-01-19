package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// SWVersion stores EOS version with integer based Major and Minor
type SWVersion struct {
	Major int    `json:"major"`
	Minor int    `json:"minor"`
	Patch string `json:"patch"`
	Meta  string `json:"meta"`
}

// NameServers is an ip name-server
type NameServers struct {
	Vrf       string
	Addresses []string
}

// LoggingHost is a syslog server
type LoggingHost struct {
	Hostname string `json:"hostname"`
}

// Logging configuratoin
type Logging struct {
	Host            LoggingHost
	On              bool
	Level           string
	SourceInterface string
	Vrf             string // TODO need additional fields for host and source interface

	/*
			  buffered            Logging buffer configuration
		  console             Set console logging parameters
		  event               Configure logging events
		  facility            Set logging facility
		  format              Set logging format parameters
		  host                Set syslog server IP address and parameters
		  level               Configure logging severity
		  monitor             Set terminal monitor parameters
		  on                  Turn on logging
		  persistent          Save logging messages to the flash disk
		  policy              Configure logging policies
		  qos                 Configure QoS parameters
		  relogging-interval  Configure relogging-interval for critical log messages
		  repeat-messages     Repeat messages instead of summarizing number of repeats
		  source-interface    Use IP Address of interface as source IP of log messages
		  synchronous         Set synchronizing unsolicited with solicited messages
		  trap                Severity of messages sent to the syslog server
	*/
}

// Hostname is used to hold the hostname for JSON serialization
type Hostname struct {
	Hostname string
}

// EOSDevices captures a list of switches
type EOSDevices struct {
	Switches []EOSDevice `json:"switches"`
}

// EOSDevice is an EOS endpoint
type EOSDevice struct {
	Hostname      string                 `json:"id"` // using id to be consistent with OPA
	Management    map[string]interface{} `json:"management"`
	IPNameServers NameServers            `json:"ipNameServers"`
	Logging       Logging                `json:"logging"`
	SWVersion     SWVersion              `json:"swVersion"`
	HWVersion     string                 `json:"hwVersion"`
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

// TODO : scan through multiple commented lines to get to non-commented line
func parseComments() {

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

func parseHostname(d EOSDevice, line []string) EOSDevice {
	d.Hostname = line[1]
	return d
}

/*
ip name-server vrf default 172.22.22.40
dns domain sjc.aristanetworks.com
!
ntp server 172.22.22.50
ntp server 198.55.111.50
ntp server 216.229.0.17
*/

// parses ip name-server command (max: 3)
// ex: ip name-server vrf default 172.22.22.40
func parseNameServers(d EOSDevice, scanner *bufio.Scanner) EOSDevice {
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
func parseDNSDomain(d EOSDevice, scanner *bufio.Scanner) EOSDevice {
	return d
}

// parses ntp server command
// ex.
// ntp server 172.22.22.50
// ntp server 198.55.111.50
// ntp server 216.229.0.17
func parseNTPServer(d EOSDevice, scanner *bufio.Scanner) EOSDevice {
	return d
}

// parses logging command
// ex: logging host 172.22.22.40
func parseLogging(d EOSDevice, scanner *bufio.Scanner) EOSDevice {
	line := strings.Fields(scanner.Text())
	switch line[0] {
	case "logging":
		d.Logging.On = true
		switch line[1] {
		case "host":
			{
				d = parseLoggingHost(d, scanner)
			}
		}
	case "no":
		{
			d = parseNoLogging(d, scanner)
		}
	}
	return d
}

func parseNoLogging(d EOSDevice, scanner *bufio.Scanner) EOSDevice {
	// line := strings.Fields(scanner.Text())
	return d
}

func parseLoggingHost(d EOSDevice, scanner *bufio.Scanner) EOSDevice {
	line := strings.Fields(scanner.Text())
	d.Logging.Host.Hostname = line[2]
	return d
}

// parseDeviceInfo assumes it is receiving a commented line with device information
// ex. ! device: DMZ-LF18 (DCS-7060SX2-48YC6, EOS-4.24.2.1F)
func parseDeviceInfo(d EOSDevice, scanner *bufio.Scanner) EOSDevice {
	line := strings.Fields(scanner.Text())
	d.SWVersion = parseSWVersion(line[4])
	d.HWVersion = parseHWVersion(line[3])
	return d
}

func parseHWVersion(line string) string {
	cleanLine := strings.TrimPrefix(line, "(")
	cleanLine = strings.TrimSuffix(cleanLine, ",")
	return cleanLine
}

func parseSWVersion(line string) SWVersion {
	cleanLine := strings.TrimPrefix(line, "EOS-")
	cleanLine = strings.TrimSuffix(cleanLine, ")")
	fields := strings.Split(cleanLine, ".")
	var major, minor int
	if i, err := strconv.Atoi(fields[0]); err == nil {
		major = i
	}
	if i, err := strconv.Atoi(fields[1]); err == nil {
		minor = i
	}
	swv := SWVersion{Major: major, Minor: minor, Patch: fields[2]}
	if len(fields) > 3 {
		swv.Meta = fields[3]
	}
	return swv
}

// createManagement() Creates an empty mangagement block with default values
func createManagement() map[string]interface{} {
	m := map[string]interface{}{
		"ssh":    ManagementSSH{TypeID: "ssh", Shutdown: false},
		"telnet": ManagementTelnet{TypeID: "telnet", Shutdown: true},
		"api":    ManagementAPIHTTP{TypeID: "http", Shutdown: true},
	}
	return m
}

// NewEOSDevice creates an empty device with defaults set
func NewEOSDevice() EOSDevice {
	m := createManagement()
	return EOSDevice{Management: m}
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
	device := NewEOSDevice()
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
					device = parseDeviceInfo(device, scanner)
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
				device = parseHostname(device, line)
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
