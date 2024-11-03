package cmd

import (
	"fmt"
	"strings"
	"strconv"
	"log"
	"github.com/spf13/cobra"
	"github.com/vrypan/fargo/fctools"
	"github.com/vrypan/fargo/config"
)

var getCmd = &cobra.Command{
	Use:   "get [URI]",
	Aliases: []string{"g"},
	Short: "Get Farcaster data",
	Long: `URI formats supported:
- @username/casts
- @username/0x<hash>
- @username/0x<hash>/embed
- @username/0x<hash>/embed/<index>
- @username/profile/[pfp|display|url|bio|username|location]`,
	Run: getRun,
}

func getRun(cmd *cobra.Command, args []string) {
	config.Load()

	// config.BindPFlag("count", getCmd.Flags().Lookup("get.count"))
	
	fid, parts := parse_url(args)
	expandFlag, _ := cmd.Flags().GetBool("expand")
	countFlag := uint32( config.GetInt("get.count") )
	if c, _:= cmd.Flags().GetInt("count") ; c >0 {
		countFlag = uint32(c)
	}

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

	// @fname/profile
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
		fmt.Println(s)
		return
	}

	// @fname/casts
	if len(parts) == 1 && parts[0] == "casts" {
		s, err := fctools.PrintCastsByFid(fid, countFlag, grepFlag)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(s)
		return
	}

	// @fname/0x<hash>
	if len(parts) == 1 && parts[0][0:2] == "0x" {
		s := fctools.PrintCast( fid, parts[0], expandFlag, grepFlag )
		fmt.Println(s)
		return
	}

	// @fname/0x<hash>/embeds/<index>
	if len(parts) >= 2 && parts[0][0:2] == "0x" && parts[1] == "embed" {
		urls := fctools.GetCastUrls(fid, parts[0], false, "")
		if len(parts) == 2 {
			for _, u := range urls {
				fmt.Println(u.Link)		
			}
		}
		if len(parts) == 3 {
			idx, err := strconv.Atoi(parts[2])
			if err != nil {
				log.Fatal("Not found")			
			}
			if idx < len(urls) {
				fmt.Println(urls[idx].Link)
			}
		}
		return
	}
	log.Fatal("Not found")
}

func init() {
	rootCmd.AddCommand(getCmd)
	config.Load()
	
	getCmd.Flags().BoolP("expand", "e", false, "Expand threads")
	getCmd.Flags().IntP("count", "c", 0, "Number of casts to show when getting @user/casts")
	getCmd.Flags().StringP("grep", "", "", "Only show casts containing a specific string")
}
