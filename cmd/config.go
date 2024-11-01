/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vrypan/fargo/config"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Aliases: []string{"c"},
	Short: "Get/Set fargo configuration parameters",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		f := config.Load()
		fmt.Printf("\nConfig file is %s\n", f)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
