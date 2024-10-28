/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (	
	"github.com/spf13/cobra"
	"github.com/vrypan/fargo/config"
)

// cdCmd represents the cd command
var cdCmd = &cobra.Command{
	Use:   "cd",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: cdRun,
}

func cdRun(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		if args[0][0] == '@' {
			// Abs path starts with @
			c := config.Config{File: "fargo.txt", CurrentPath: args[0]}	
			c.Save()
		} else {
			// Rel path
			c := config.Load("fargo.txt")
			new_path := c.CurrentPath + "/" + args[0]	
			c.CurrentPath = new_path
			c.Save()
		}
	}
	
}

func init() {
	rootCmd.AddCommand(cdCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cdCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cdCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
