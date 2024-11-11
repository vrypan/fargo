package cmd

/*
snapshot @fname/cast
-r = recursively fetch the whole thread
--out directory name where to store the snapshot
*/
import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vrypan/fargo/config"
	"github.com/vrypan/fargo/fctools"
	"github.com/vrypan/fargo/localdb"
	"github.com/vrypan/fargo/tui"
	"github.com/vrypan/fargo/urls"
)

var snapshotCmd = &cobra.Command{
	Use:     "snapshot [URI]",
	Aliases: []string{"g"},
	Short:   "Create a cast/thread snapshot",
	Long: `The snapshot includes a text version of the thread as well as json
files containing all the messages in the thread.`,
	Run: getSnapshot,
}

func getSnapshot(cmd *cobra.Command, args []string) {
	config.Load()
	user, parts := parse_url(args)
	if user == nil {
		log.Fatal("User not found")
	}
	if len(parts) != 1 || !strings.HasPrefix(parts[0], "0x") {
		log.Fatal("Expected @fname/0x<hash> path.")
	}
	expandFlag, _ := cmd.Flags().GetBool("recursive")
	outFlag, _ := cmd.Flags().GetString("out")

	/*
		Create the output directory
	*/
	if outFlag == "" {
		outFlag = filepath.Join(config.GetString("download.dir"), "snapshot-"+parts[0])
	}
	log.Printf("Snapshot path: %s", outFlag)
	var err error
	outFlag = os.ExpandEnv(outFlag)
	path, err := filepath.Abs(outFlag)
	if err != nil {
		log.Fatalf("Failed to get absolute path of %s: %v", outFlag, err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, 0755)
		if err != nil {
			log.Fatalf("Failed to create output directory %s: %v", path, err)
		}
	}

	hub := fctools.NewFarcasterHub()
	defer hub.Close()

	log.Println("Fetching casts...")
	casts := fctools.NewCastGroup().FromCastFidHash(hub, user.Fid, parts[0][2:], expandFlag)

	s := tui.PprintThread(casts, nil, 0, "", "")
	err = os.WriteFile(filepath.Join(path, "thread.txt"), []byte(s), 0644)
	if err != nil {
		log.Fatalf("Failed to write thread.txt file: %v", err)
	}
	b, err := casts.JsonThread(true, false)
	if err != nil {
		log.Fatal("Error generating thread.json")
	}
	err = os.WriteFile(filepath.Join(path, "thread.json"), b, 0644)
	if err != nil {
		log.Fatalf("Failed to write thread.txt file: %v", err)
	}

	log.Println("Downloading embeded URLs...")
	allUrls := make([]*urls.Url, len(casts.Links()))
	urlMap := make(map[string]string)
	for i, l := range casts.Links() {
		allUrls[i] = urls.NewUrl(l).UpdateType().UpdateExt()
		urlMap[l] = allUrls[i].Filename()
	}
	for _, u := range allUrls {

		GetFile(u.Link, path, u.Filename(), false)
	}
	urlsJson, err := json.MarshalIndent(urlMap, "", "  ")
	err = os.WriteFile(filepath.Join(path, "embedsmap.json"), urlsJson, 0644)
	if err != nil {
		log.Fatalf("Failed to write embedsmap.json file: %v", err)
	}

	log.Println("Downloading PFPs...")
	localdb.Open()
	defer localdb.Close()
	for fid := range casts.Fnames {
		pfp, _ := hub.PrxGetUserDataStr(fid, "USER_DATA_TYPE_PFP")
		url := urls.NewUrl(pfp).UpdateExt().UpdateExt()
		ext := url.Ext()
		if ext == "" {
			ext = "png"
		}
		fidStr := strconv.FormatUint(fid, 10)
		// GetFile(pfp, path, fidStr+"."+ext, false)
		// log.Printf("Downloading: %s", pfp)
		req, err := http.NewRequest("GET", pfp, nil)
		if err != nil {
			log.Printf("Failed to create request for %s: %v", pfp, err)
		} else {
			req.Header.Set("User-Agent", "curl/8.7.1")
			req.Header.Set("Accept", "*/*")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Failed to download link %s: %v", pfp, err)
			} else {
				defer resp.Body.Close()
				outFile, err := os.Create(filepath.Join(path, fidStr+"."+ext))
				if err != nil {
					log.Fatalf("Failed to create file %s: %v", fidStr+"."+ext, err)
				}
				defer outFile.Close()

				_, err = io.Copy(outFile, resp.Body)
				if err != nil {
					log.Fatalf("Failed to save file %s: %v", fidStr+"."+ext, err)
				}
			}
		}
	}
	log.Println("Done.")
}

func init() {
	rootCmd.AddCommand(snapshotCmd)
	snapshotCmd.Flags().BoolP("recursive", "r", false, "Recursively get parent casts and replies")
	snapshotCmd.Flags().StringP("out", "", "", "Output directory")
}
