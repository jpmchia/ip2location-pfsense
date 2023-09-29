package main

import (
	"flag"

	"github.com/jpmchia/ip2location-pfsense/backend/cmd"
	"github.com/jpmchia/ip2location-pfsense/backend/util"
)

var appName string = "IP2Location-pfSense"

func main() {
	util.Verbose = true
	util.Log("[main] Starting %v ...", appName)

	flag.Parse()

	util.Log("[main] Initialising service ...")

	cmd.Execute()
}
