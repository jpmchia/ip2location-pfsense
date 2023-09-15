// File: cmd/root.go
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "IP2LOCATION",
	Short: "IP2Location backend service for pfSense",
	Long: `
p2location_pfsense is a backend service and CLI tool for retrieving and 
displaying IP geolocation information on pfSense devices. This service 
is designed for use with the IP2Location pfSense dashboard widget. 

Download the widget from: https://github.com/jpmchia/pfsense_ip2location

The service facilitates the retrieval of geolocation and other auxiliary
information assocated with a specified IPv4 or IPv6 address from the API
provided by IP2Location.io. 

Register for a free API account at: https://www.ip2location.io/dashboard

Optionally, lookup and query results may cached locally in a Redis store
to improve response times and reduce the number of calls to the API.`,
	Run: func(cmd *cobra.Command, args []string) {},
}

// Adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// init is called by Cobra prior to any command execution.
// Define flags and configuration settings.
// Cobra supports persistent flags, which, if defined here,
// are global for the application.
func init() {
	// var cfgFile string
	cfgFile := "config.yaml"
	// var useRedis bool
	useRedis := true
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.IP2LOCATION.yaml)")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "/usr/local/ip2location/config.yaml", "specify the location of the configuration file")
	rootCmd.PersistentFlags().BoolVar(&useRedis, "use-cache", true, "enable the use of Redis to cache results ")
}
