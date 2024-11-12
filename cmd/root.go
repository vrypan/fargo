package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/vrypan/fargo/config"
)

var rootCmd = &cobra.Command{
	Use:   "fargo",
	Short: "A command line tool to interact with Farcaster",
}

func Execute() {
	log.SetFlags(0)
	config.Load()
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if cmd.Name() != "set" {
			warnHoyt()
		}
	}
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func warnHoyt() {
	if hub := config.GetString("hub.host"); hub == "hoyt.farcaster.xyz" && config.GetString("warn.hoyt") != "off" {
		log.Println("===========================================================")
		log.Println(" WARNING: You are using the default hub: hoyt.farcaster.xyz")
		log.Println(" For better performance, consider using a dedicated hub.")
		log.Println(" To turn off this warning use the following command:\n\n fargo config set warn.hoyt off")
		log.Println("===========================================================")
	}
}
