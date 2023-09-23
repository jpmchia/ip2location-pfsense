/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/jpmchia/ip2location-pfsense/backend/config"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration commands",
	Long:  `Use this command to create configuration files and show the current configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Specify a subcommand to use or -h to display a list of valid commands.")
	},
}

var createConfigCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new configuration file with default values",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Executing config create command")
		config.CreateConfigFile(args)
	},
}

var showConfigCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the current configuration",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		config.ShowConfig()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(createConfigCmd)
	configCmd.AddCommand(showConfigCmd)
}
