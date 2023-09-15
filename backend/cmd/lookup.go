package cmd

import (
	"github.com/spf13/cobra"
)

// lookupCmd represents the lookup command
var lookupCmd = &cobra.Command{
	Use:   "lookup [IP address]",
	Short: "Retrieve information for the specified IP address",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// ip2location.LookupIP(args)
	},
}

func init() {

	rootCmd.AddCommand(lookupCmd)

}
