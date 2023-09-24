package main

import (
	"flag"

	"github.com/jpmchia/ip2location-pfsense/backend/cmd"
	"github.com/jpmchia/ip2location-pfsense/backend/util"
)

var appName string = "IP2Location-pfSense"

func main() {
	util.LogDebug("[main] Starting %v ...", appName)

	flag.Parse()

	util.LogDebug("[main] Initialising service ...")

	cmd.Execute()
}
