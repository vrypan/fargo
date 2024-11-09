/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Submit messages to the network",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
}
