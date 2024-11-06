package cmd

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vrypan/fargo/config"
	"github.com/vrypan/fargo/fctools"
)

var getCmd = &cobra.Command{
	Use:     "get [URI]",
	Aliases: []string{"g"},
	Short:   "Get Farcaster data",
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

	fid, parts := parse_url(args)
	expandFlag, _ := cmd.Flags().GetBool("expand")
	countFlag := uint32(config.GetInt("get.count"))
	if c, _ := cmd.Flags().GetInt("count"); c > 0 {
		countFlag = uint32(c)
	}

	grepFlag, _ := cmd.Flags().GetString("grep")

	hub := fctools.NewFarcasterHub()
	defer hub.Close()

	if fid == 0 {
		if b, e := hub.HubInfo(); e != nil {
			log.Fatalf("Error! %v", e)
		} else {
			fmt.Println(string(b))
		}
		return
	}

	switch {
	case len(parts) == 2 && parts[0] == "profile":
		req := strings.Split(parts[1], ".")
		json := len(req) > 1 && req[1] == "message"
		if s, err := hub.GetUserData(fid, "USER_DATA_TYPE_"+strings.ToUpper(req[0]), json); err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(s)
		}
	case len(parts) == 1 && parts[0] == "casts":
		if s, err := fctools.PrintCastsByFid(fid, countFlag, grepFlag); err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(s)
		}
	case len(parts) == 1 && strings.HasPrefix(parts[0], "0x"):
		fmt.Println(fctools.PrintCast(fid, parts[0], expandFlag, grepFlag))
	case len(parts) >= 2 && strings.HasPrefix(parts[0], "0x") && parts[1] == "embed":
		urls := fctools.GetCastUrls(fid, parts[0], false, "")
		if len(parts) == 2 {
			for _, u := range urls {
				fmt.Println(u.Link)
			}
		} else if len(parts) == 3 {
			if idx, err := strconv.Atoi(parts[2]); err == nil && idx < len(urls) {
				fmt.Println(urls[idx].Link)
			} else {
				log.Fatal("Not found")
			}
		}
	default:
		log.Fatal("Not found")
	}
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().BoolP("expand", "e", false, "Expand threads")
	getCmd.Flags().IntP("count", "c", 0, "Number of casts to show when getting @user/casts")
	getCmd.Flags().StringP("grep", "", "", "Only show casts containing a specific string")
}
