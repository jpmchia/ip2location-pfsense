package main

import (
	"flag"
	"ip2location-pfsense/cmd"
	. "ip2location-pfsense/util"
	"log"
)

var appName string = "IP2Location-pfSense"

func main() {
	log.Default().Printf("[main] Starting %v ...", appName)

	flag.Parse()
	debugPtr := flag.Bool("debug", false, "verbose output")

	if *debugPtr {
		Debug = true
		LogDebug("[main] Debugging enabled")
	}

	LogDebug("[main] Initialising service ...")

	cmd.Execute()
}
