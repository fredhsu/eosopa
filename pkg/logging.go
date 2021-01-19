package eosparse

import (
	"bufio"
	"strings"
)

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
