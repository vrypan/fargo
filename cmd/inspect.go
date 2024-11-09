package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/vrypan/fargo/fctools"
)

var inspectCmd = &cobra.Command{
	Use:     "inspect [cast URI]",
	Aliases: []string{"g"},
	Short:   "Inspect a cast",
	Long: `Returns a json of the corresponding message.
The URI must be in the form: @username/0x<cast hash>`,
	Run: inspectRun,
}

func inspectRun(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		return
	}
	user, parts := parse_url(args)
	// hexFlag, _ := cmd.Flags().GetBool("hex")
	// datesFlag, _ := cmd.Flags().GetBool("dates")

	hub := fctools.NewFarcasterHub()
	defer hub.Close()

	if len(parts) != 1 || parts[0][:2] != "0x" {
		log.Fatal("Not found")
	}

	if user.Fid == 0 {
		log.Fatal("User not found. ", user, parts)
	}
	casts := fctools.NewCastGroup().FromCastFidHash(hub, user.Fid, parts[0][2:], false)
	b, err := casts.JsonList()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", string(b))
}

func init() {
	rootCmd.AddCommand(inspectCmd)
	inspectCmd.Flags().BoolP("hex", "", true, "Show binary fields as hex")
	inspectCmd.Flags().BoolP("dates", "", false, "Convert Farcaster timestamps to dates")
}
