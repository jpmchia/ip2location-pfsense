/*
Author: Jean-Paul Chia
Copyright © 2023 TerraNet UK <info@terranet.uk>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
Author: Jean-Paul Chia
Copyright © 2023 TerraNet UK <info@terranet.uk>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"os"

	"ip2location-pfsense/config"
	"ip2location-pfsense/util"
	. "ip2location-pfsense/util"

	"ip2location-pfsense/config"
	"ip2location-pfsense/util"
	. "ip2location-pfsense/util"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var redis bool

var redis bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ip2location-pfsense",
	Short: "IP2Location backend service for pfSense",
	Long: `
p2location_pfsense is a backend service and CLI tool for retrieving and 
displaying IP geolocation information on pfSense devices. This service 
is designed for use with the IP2Location pfSense dashboard widget. 

Download the widget from: https://github.com/jpmchia/IP2Location-pfSense
Download the widget from: https://github.com/jpmchia/IP2Location-pfSense

The service facilitates the retrieval of geolocation and other auxiliary
information assocated with a specified IPv4 or IPv6 address from the API
provided by IP2Location.io. 

Register for a free API account at: https://www.ip2location.io/dashboard`,
Register for a free API account at: https://www.ip2location.io/dashboard`,
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

// init is called prior to any command execution.
// init is called prior to any command execution.
func init() {
	// var cfgFile string
	cfgFile := "config.yaml"

	// Define global persistent flags
	rootCmd.PersistentFlags().StringVarP(&config.CfgFile, "config", "c", cfgFile, "specifiy the filename and path of the configiration file")
	rootCmd.PersistentFlags().BoolVarP(&util.Debug, "debug", "v", false, "output verbose debugging information")
	// rootCmd.PersistentFlags().BoolVar(&redis, "redis", true, "use Redis cache")

	// Bind flags to viper
	//viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	if Debug {
		LogDebug("Debugging verbose mode enabled")
	}

	// Load configuration
	config.LoadConfigProvider(cfgFile)
}
