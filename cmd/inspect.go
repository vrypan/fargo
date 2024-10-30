/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"encoding/hex"

	"github.com/spf13/cobra"
	"github.com/vrypan/fargo/fctools"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect [cast URI]",
	Aliases: []string{"g"},
	Short: "Inspect a cast",
	Long: `The URI must be in the form:
@username/0x<cast hash>`,
	Run: inspectRun,
}

func inspectRun(cmd *cobra.Command, args []string) {
	fid, parts := parse_url(args)
	hexFlag, _ := cmd.Flags().GetBool("hex")
	datesFlag, _ := cmd.Flags().GetBool("dates")

	hub := fctools.NewFarcasterHub()
	defer hub.Close()

	var jsonData []byte

	if len(parts) == 1 && parts[0][0:2] == "0x" {
		hash_bytes, err1 := hex.DecodeString(parts[0][2:])
		if err1 != nil {
			log.Fatal("Hash is not a hex number")
		}
		msg,err := hub.GetCast(fid, hash_bytes)
		if err != nil {
			log.Fatal(err)
		}

		jsonData, err = fctools.Marshal(msg, fctools.MarshalOptions{Bytes2Hash: hexFlag, Timestamp2Date: datesFlag}) 
		if err != nil {
			log.Fatalf("Error converting message to JSON: %v", err)
		}
		fmt.Println(string(jsonData))
		return
	}
	log.Fatal("Not found")
}

func init() {
	rootCmd.AddCommand(inspectCmd)
	inspectCmd.Flags().BoolP("hex", "", true, "Show binary fields as hex")
	inspectCmd.Flags().BoolP("dates", "", false, "Convert Farcaster timestamps to dates")
}
