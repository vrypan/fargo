package cmd

import (
	"os"
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fargo",
	Short: "A command line tool to interact with Farcaster",
}

func Execute() {
	log.SetFlags(0)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}


