package cmd

import (
	"github.com/spf13/cobra"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:     "post",
	Short:   "Submit messages to the network",
	Aliases: []string{"send"},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
}
