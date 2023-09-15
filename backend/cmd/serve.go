package cmd

import (
	"pfSense/service"

	"github.com/spf13/cobra"
)

// lookupCmd represents the lookup command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the IP2Location backend service for pfSense",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		service.Service(args)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
