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

var DEBUG = 1

func Execute() {
	if DEBUG != 0 {
		logFile, err := os.OpenFile("fargo.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("failed to open log file: %v", err)
		}
		defer logFile.Close()
		log.SetOutput(logFile)
	} else {
		log.SetFlags(0)
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	config.Load()
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if cmd.Name() != "set" {
			warnHoyt()
		}
	}
	err2 := rootCmd.Execute()
	if err2 != nil {
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
