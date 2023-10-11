package main

import (
	"github.com/jpmchia/ip2location-pfsense/cmd"
	"github.com/jpmchia/ip2location-pfsense/util"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var appName string = "ip2location-pfsense"

func init() {
	util.Log("[main] Initialising service ...")
}

func main() {

	pflag.BoolVarP(&util.Verbose, "verbose", "v", false, "Enable verbose logging to stdout.")
	pflag.BoolVarP(&util.Debug, "debug", "d", false, "Enable dutput verbose debugging information.")
	pflag.Parse()

	err := viper.BindPFlag("verbose", pflag.Lookup("verbose"))
	if err != nil {
		util.Log("[main] Error binding verbose flag: %v", err)
	}
	err = viper.BindPFlag("debug", pflag.Lookup("debug"))
	if err != nil {
		util.Log("[main] Error binding debug flag: %v", err)
	}

	util.Verbose = viper.GetBool("verbose")
	util.Debug = viper.GetBool("debug")

	util.Log("[main] Starting %v ...", appName)

	cmd.Execute()
}
