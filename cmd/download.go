package cmd

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/hashicorp/go-getter"
	"github.com/spf13/cobra"
	"github.com/vrypan/fargo/config"
	"github.com/vrypan/fargo/fctools"
	"github.com/vrypan/fargo/fctools/embedurl"
)

var downloadCmd = &cobra.Command{
	Use:     "download [URI]",
	Aliases: []string{"d"},
	Short:   "Download Farcaster-embedded URLs",
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
	countFlag := uint32(config.GetInt("get.count"))
	if c, _ := cmd.Flags().GetInt("count"); c > 0 {
		countFlag = uint32(c)
	}

	grepFlag, _ := cmd.Flags().GetString("grep")
	dirFlag, _ := cmd.Flags().GetString("dir")
	mimetypeFlag, _ := cmd.Flags().GetString("mime-type")
	dryrunFlag, _ := cmd.Flags().GetBool("dry-run")
	skipdownloadedFlag, _ := cmd.Flags().GetBool("skip-downloaded")

	var download_dir string
	if dirFlag == "" {
		config.Load()
		download_dir = config.GetString("download.dir")
	} else {
		download_dir = dirFlag
	}
	if download_dir == "" {
		download_dir = "."
	}
	download_dir = normalizeLocalPath(download_dir)
	fmt.Println(download_dir)

	hub := fctools.NewFarcasterHub()
	defer hub.Close()

	processURLs := func(urls []embedurl.Url) {
		for _, u := range urls {
			m := u.UpdateContentType()
			if mimetypeFlag == "" || (len(m) >= len(mimetypeFlag) && m[:len(mimetypeFlag)] == mimetypeFlag) {
				if !dryrunFlag {
					GetFile(u.Link, download_dir, u.Filename(), skipdownloadedFlag)
				}
				fmt.Printf("%s --> %s\n", u.Link, u.Filename())
			}
		}
	}

	if len(parts) == 1 && parts[0] == "casts" {
		fmt.Println("Count: ", countFlag)
		urls := fctools.GetFidUrls(fid, countFlag, grepFlag)
		processURLs(urls)
		return
	}
	if len(parts) == 1 && parts[0][0:2] == "0x" {
		urls := fctools.GetCastUrls(fid, parts[0], expandFlag, grepFlag)
		processURLs(urls)
		return
	}
	log.Fatal("Not found")
}

func normalizeLocalPath(p string) string {
	if len(p) > 0 && p[0] == '~' {
		usr, err := user.Current()
		if err != nil {
			log.Fatalf("Error getting current user: %v\n", err)
		}
		return filepath.Join(usr.HomeDir, p[1:])
	}
	return p
}

func fileExists(f string) bool {
	_, err := os.Stat(f)
	return err == nil
}

func GetFile(url string, dst_dir string, dst_file string, skipdownloadedFlag bool) string {
	if err := os.MkdirAll(dst_dir, os.ModePerm); err != nil {
		log.Fatalf("Error creating directory: %v\n", err)
	}

	path := filepath.Join(dst_dir, dst_file)
	if skipdownloadedFlag && fileExists(path) {
		return path
	}

	if err := getter.GetFile(path, url); err != nil {
		log.Printf("\n%v: Error downloading file: %v\n", url, err)
		return ""
	}
	return path
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().BoolP("expand", "e", false, "Expand threads")
	downloadCmd.Flags().IntP("count", "c", 0, "Number of casts to show when getting @user/casts")
	downloadCmd.Flags().StringP("grep", "", "", "Only show casts containing a specific string")
	downloadCmd.Flags().StringP("mime-type", "", "", "Download embeds of mime/type")
	downloadCmd.Flags().BoolP("dry-run", "", false, "Do not download the files, just print the URLs and local destination")
	downloadCmd.Flags().BoolP("skip-downloaded", "", true, "If local file exists, do not download")
	downloadCmd.Flags().StringP("dir", "", "", "Download directory. If not specified, the 'downloads.dir' config is used.")
}
