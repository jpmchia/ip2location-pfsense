package cmd

import (
	"log"

	"github.com/jpmchia/IP2Location-pfSense/ip2location"
	"github.com/jpmchia/IP2Location-pfSense/service"

	"github.com/spf13/cobra"
)

// lookupCmd represents the lookup command
var countersCmd = &cobra.Command{
	Use:   "counters",
	Short: "Manage the IP2Location API call counters",
	Long:  `Display / set the current value of the counters.`,
	Run: func(cmd *cobra.Command, args []string) {
		service.Start(args)
	},
}

var createCountersCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new counters file and set the counters to zero",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Executing counters create command")
		ip2location.CreateCountersFile(args)
	},
}

var showCountersCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the current value of the counters",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Executing counters show command")
		ip2location.ShowCounters()
	},
}

var resetCountersCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the counters",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Executing reset counters command")
		ip2location.ResetCounters(args)
	},
}

func init() {
	rootCmd.AddCommand(countersCmd)
	countersCmd.AddCommand(createCountersCmd)
	countersCmd.AddCommand(showCountersCmd)
	countersCmd.AddCommand(resetCountersCmd)
}
