/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"strings"
	"strconv"
	"os"
	"log"

	"github.com/spf13/cobra"
	//"github.com/vrypan/fargo/config"
	"github.com/vrypan/fargo/fctools"
	"github.com/vrypan/fargo/ipfs"
)

var getCmd = &cobra.Command{
	Use:   "get [URI]",
	Aliases: []string{"g"},
	Short: "Get Farcaster data",
	Long: `URI formats supported:
- @username/casts
- @username/0x<cast hash>
- @username/profile/[pfp|display|url|bio|username|location]`,
	Run: getRun,
}

func parse_url(args []string) (uint64, []string) {
	if len(args) == 0 {
		log.Fatal("No path")
	}

	parts := strings.Split(args[0], "/")
	
	if parts[0][0:1] != "@" {
		log.Fatal("Path should start with @")
	}

	var fid uint64
	var err error

	if parts[0] == "@" {
		return 0, nil
	}
	if fid, err = strconv.ParseUint(parts[0][1:], 10, 64); err != nil {
    	fid, err = fctools.GetFidByFname(parts[0][1:])
    	if err != nil {
    		log.Fatalf("Error looking up %v [%v]", parts[0], err)
    	}
	}
	return fid, parts[1:]
}

func getRun(cmd *cobra.Command, args []string) {
	fid, parts := parse_url(args)
	expandFlag, _ := cmd.Flags().GetBool("expand")
	ipfsFlag, _ := cmd.Flags().GetBool("ipfs")
	countFlag, _ := cmd.Flags().GetUint32("count")
	grepFlag, _ := cmd.Flags().GetString("grep")

	hub := fctools.NewFarcasterHub()
	defer hub.Close()

	if fid==0 {
		b,e := hub.HubInfo()
		if e != nil {
			log.Fatalf("Error! %v", e)
		}
		fmt.Println(string(b))
		return
	}

	if len(parts) == 2 && parts[0] == "profile" {
		//fid, _ := strconv.ParseUint(parts[0], 10, 64)
		req := strings.Split(parts[1], ".")
		var json = false
		if len(req)>1 && req[1] == "message" { 
			json = true 
		}
		s, err := hub.GetUserData(fid, "USER_DATA_TYPE_"+strings.ToUpper(req[0]), json)
		if err != nil {
			log.Fatal(err)
		}
		if ipfsFlag && s[0:7]=="ipfs://" {
			download_path := ipfs.GetCid(s[0:7])
			fmt.Println(download_path)
			return
		} 
		fmt.Println(s)
		return
	}
	if len(parts) == 1 && parts[0] == "casts" {
		s, err := fctools.PrintCastsByFid(fid, countFlag, grepFlag)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(s)
		return
	}
	if len(parts) == 1 && parts[0][0:2] == "0x" {
		s := fctools.PrintCast( fid, parts[0], expandFlag, grepFlag )
		fmt.Println(s)
		return
	}
	fmt.Println("Not found")
	os.Exit(1)
}

func init() {

	rootCmd.AddCommand(getCmd)
	getCmd.Flags().BoolP("expand", "e", false, "Expand threads")
	getCmd.Flags().BoolP("ipfs", "", false, "If URI is @user/profile/url and it points to an ipfs;// URL, download it.")
	getCmd.Flags().Uint32P("count", "c", 20, "Number of casts to show when getting @user/casts")
	getCmd.Flags().StringP("grep", "", "", "Only show casts containing a specific string")
	//getCmd.Flags().StringVarP("source", "s", "", "Source directory to read from")


	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
