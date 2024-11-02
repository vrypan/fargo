package cmd

import (
	"fmt"
	"log"
	"github.com/spf13/cobra"
	"github.com/vrypan/fargo/fctools"
	"github.com/vrypan/fargo/urltools"
)

var downloadCmd = &cobra.Command{
	Use:   "download [URI]",
	Aliases: []string{"g"},
	Short: "Download Farcaster-embedded URLs",
	Long: `This command works like "get", but instead
of displaying casts, it downloads the URLs embedded in
these casts.

Use the --mime-type flag to indicate the type of embedded
URLs you want to download.`,
	Run: downloadRun,
}

func downloadRun(cmd *cobra.Command, args []string) {
	fid, parts := parse_url(args)
	expandFlag, _ := cmd.Flags().GetBool("expand")
	countFlag, _ := cmd.Flags().GetUint32("count")
	grepFlag, _ := cmd.Flags().GetString("grep")
	//mimetypeFlag, _ := cmd.Flags().GetString("mime-type")

	hub := fctools.NewFarcasterHub()
	defer hub.Close()

	if len(parts) == 1 && parts[0] == "casts" {
		urls := fctools.GetFidUrls(fid, countFlag, grepFlag)
		for _, u := range urls {
			m, err := urltools.GetMimeType(u)
			if err != nil {
				m = "???"
			}
			fmt.Println(m, u)
		}
		return
	}
	if len(parts) == 1 && parts[0][0:2] == "0x" {
		s := fctools.PrintCast( fid, parts[0], expandFlag, grepFlag )
		fmt.Println(s)
		return
	}
	log.Fatal("Not found")
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().BoolP("expand", "e", false, "Expand threads")
	downloadCmd.Flags().Uint32P("count", "c", 20, "Number of casts to show when getting @user/casts")
	downloadCmd.Flags().StringP("grep", "", "", "Only show casts containing a specific string")
	downloadCmd.Flags().StringP("mime-type", "", "", "Download embeds of mime/type")
}
