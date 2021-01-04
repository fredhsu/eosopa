package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// ManagementIntf gives common management functions
type ManagementIntf interface {
	Type() string
	Enabled() bool
}

// ManagementAPIHTTP stores data about eAPI config
type ManagementAPIHTTP struct {
	typeID    string
	protocols []string
	Shutdown  bool
}

// Enabled indicates if the feature is enabled
func (m ManagementAPIHTTP) Enabled() bool {
	return !m.Shutdown
}

// Type gives a string value of the management type
func (m ManagementAPIHTTP) Type() string {
	return m.typeID
}

// ManagementTelnet stores telnet settings
type ManagementTelnet struct {
	typeID   string
	Shutdown bool
}

// Type gives a string value of the management type
func (m ManagementTelnet) Type() string {
	return m.typeID
}

// Enabled indicates if the feature is enabled
func (m ManagementTelnet) Enabled() bool {
	return !m.Shutdown
}

// ManagementSSH stores ssh settings
type ManagementSSH struct {
	typeID     string
	Shutdown   bool
	ServerPort int
}

// Enabled indicates if the feature is enabled
func (m ManagementSSH) Enabled() bool {
	return !m.Shutdown
}

// Type gives a string value of the management type
func (m ManagementSSH) Type() string {
	return m.typeID
}

func parseManagement(scanner *bufio.Scanner, line []string) ManagementIntf {
	switch mgmt := line[1]; mgmt {
	case "api":
		fmt.Println("api")
		return ManagementAPIHTTP{}
	case "telnet":
		return parseTelnet(scanner)
	case "ssh":
		return parseSSH(scanner)
	default:
		fmt.Println("Not a recognized management type")
		return nil
	}
}

func parseSSH(scanner *bufio.Scanner) ManagementSSH {
	m := ManagementSSH{typeID: "ssh", Shutdown: false}
	line := strings.Fields(scanner.Text())
	if contains(line, "shutdown") {
		m.Shutdown = parseShutdown(line)
	}
	return m
}

func parseTelnet(scanner *bufio.Scanner) ManagementTelnet {
	m := ManagementTelnet{typeID: "telnet"}
	line := strings.Fields(scanner.Text())
	if contains(line, "shutdown") {
		m.Shutdown = parseShutdown(line)
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
func main() {
	file, err := os.Open("./dmz-lf18.config")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if line[0] == "management" {
			mgmt := parseManagement(scanner, line)
			fmt.Printf("got a %+v\n", mgmt)
		}
		i++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
