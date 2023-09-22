package cmd

import (
	"ip2location-pfsense/ip2location"
	"log"

	"github.com/spf13/cobra"
)

// lookupCmd represents the lookup command
var lookupCmd = &cobra.Command{
	Use:   "lookup [IP address]",
	Short: "Retrieve information for the specified IP address",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := ip2location.RetrieveIpLocation(args[0], args[0])
		if err != nil {
			log.Fatal("Unable to add lookup command")
		}
	},
}

func init() {

	rootCmd.AddCommand(lookupCmd)

}
