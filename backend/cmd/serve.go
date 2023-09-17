package cmd

import (
	"ip2location-pfsense/service"

	"github.com/spf13/cobra"
)

// lookupCmd represents the lookup command
var serveCmd = &cobra.Command{
	Use:   "service",
	Short: "Start the IP2Location service",
	Long: `Starts the IP2Location service. This service will listen for incoming requests and 
		respond with the appropriate information.`,
	Run: func(cmd *cobra.Command, args []string) {
		service.Start(args)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
