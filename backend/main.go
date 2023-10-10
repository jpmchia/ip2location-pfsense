package main

import (
	"flag"

	"github.com/jpmchia/ip2location-pfsense/cmd"
	"github.com/jpmchia/ip2location-pfsense/util"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var appName string = "ip2location-pfsense"

func main() {
	util.Verbose = true
	util.Log("[main] Starting %v ...", appName)

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	util.Log("[main] Initialising service ...")

	cmd.Execute()
}
