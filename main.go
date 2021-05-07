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

func main() {
	filePtr := flag.String("input", "eos.config", "config file to convert")
	flag.Parse()

	file, err := os.Open(*filePtr)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	devices := EOSDevices{}
	device := eosparse.NewEOSDevice()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		// TODO: How to handle "no" prefix? contains? -- will this only take place in subconfig?
		switch line[0] {
		case "!": // comments
			{
				if len(line) > 1 && line[1] == "device:" {
					// The first few lines of config have device information embedded
					device = eosparse.ParseDeviceInfo(device, scanner)
				} else {
					// Not device info related comment string so skip
					continue
				}
			}
		case "management":
			{
				device = eosparse.ParseManagement(device, scanner)
			}
		case "hostname":
			{
				device = eosparse.ParseHostname(device, line)
			}
		case "logging":
			{
				device = eosparse.ParseLogging(device, scanner)
			}
		case "no":
			{
				//device = eosparse.NoCommand(device, line)
                continue
			}
        default:
            {
                continue
            }
		}
	}
	devices.Switches = append(devices.Switches, device)
	d, err := json.MarshalIndent(devices, " ", "  ")
	fmt.Printf("%s\n", d)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
