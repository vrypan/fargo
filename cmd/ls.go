/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"strings"
	"os"

	"github.com/spf13/cobra"
	"github.com/vrypan/fargo/config"
)

var rootLs = []string {"hub_info", "fids"}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: lsRun,
}

func lsRun(cmd *cobra.Command, args []string) {
	cfg := config.Load("fargo.txt")
	path := cfg.CurrentPath
	if len(args) == 0 {
		path = cfg.CurrentPath
	} else {
		path = args[0]
	}

	if path == "" {
		return
	}
	fmt.Println("--- Listing contents of ", path)

	parts := strings.Split(path, "/")
	fmt.Println("---", parts)
	if len(parts) == 1 {
		for _, p := range []string{"profile", "casts", "links"} {
			fmt.Println(p)
		}
		return
	}
	if len(parts) == 2 && parts[1] == "profile" {
		for _, p := range []string{"pfp", "display", "bio", "url", "username", ".json"} {
			fmt.Println(p)
		}
		return
	}
	

	fmt.Println("Not found")
	os.Exit(1)
	/*var path = "/"
	if len(args) == 0 {
		path = cfg.CurrentPath
	} else {
		path = args[0]
	}
	fmt.Println("Listing contents of ", path)
	if path == "/" {
		for _, i := range rootLs {
			fmt.Println(i)
		}
	}
	*/
}

func init() {
	rootCmd.AddCommand(lsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

/*

- <user>
- <user>/profile/(pfp|display|bio|url|username|casts|<link_type>)
- <user>/casts/0x<hash>
- <user>/casts/0x<hash>/<reaction_type>s
- <user>/casts/0x<hash>/replies
- <user>/links/like
*/