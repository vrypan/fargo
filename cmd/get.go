package cmd

/*
fargo
	-X GET/POST
	-r recursive
	@username
	@username/<hash>
	@username/profile/<USER_DATA_TYPE>



*/
import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vrypan/fargo/config"
	"github.com/vrypan/fargo/fctools"
	"github.com/vrypan/fargo/tui"
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
	user, parts := parse_url(args)
	if user == nil {
		log.Fatal("User not found")
	}
	expandFlag, _ := cmd.Flags().GetBool("expand")
	countFlag := uint32(config.GetInt("get.count"))
	if c, _ := cmd.Flags().GetInt("count"); c > 0 {
		countFlag = uint32(c)
	}

	//grepFlag, _ := cmd.Flags().GetString("grep")

	hub := fctools.NewFarcasterHub()
	defer hub.Close()

	switch {
	case len(parts) == 1 && parts[0] == "profile":
		s := user.FetchUserData(hub, nil).String()
		fmt.Println(s)
	case len(parts) == 2 && parts[0] == "profile":
		t := strings.ToUpper("USER_DATA_TYPE_" + parts[1])
		s := user.FetchUserData(hub, []string{t}).Value(t)
		fmt.Println(s)
	case len(parts) == 1 && parts[0] == "casts":
		// TBA: grepFlag
		casts := fctools.NewCastGroup().FromFid(hub, user.Fid, countFlag)
		s := tui.PprintList(casts, nil, 0)
		fmt.Println(s)
	case len(parts) == 1 && strings.HasPrefix(parts[0], "0x"):
		// TBA: grepFlag
		casts := fctools.NewCastGroup().FromCastFidHash(hub, user.Fid, parts[0][2:], expandFlag)
		fmt.Println(tui.PprintThread(casts, nil, 0))
	case len(parts) >= 2 && strings.HasPrefix(parts[0], "0x") && parts[1] == "embed":
		casts := fctools.NewCastGroup().FromCastFidHash(hub, user.Fid, parts[0][2:], false)
		embeds := casts.Messages[casts.Head].Message.Data.GetCastAddBody().GetEmbeds()
		if len(parts) == 2 {
			for _, u := range embeds {
				if u.GetUrl() != "" {
					fmt.Println(u.GetUrl())
				}
			}
		} else if len(parts) == 3 {
			if idx, err := strconv.Atoi(parts[2]); err == nil && idx < len(embeds) {
				fmt.Println(embeds[idx].GetUrl())
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
