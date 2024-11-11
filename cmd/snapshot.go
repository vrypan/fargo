package cmd

/*
snapshot @fname/cast
-r = recursively fetch the whole thread
--out directory name where to store the snapshot

*/
import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vrypan/fargo/config"
	"github.com/vrypan/fargo/fctools"
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

	// countFlag := uint32(config.GetInt("get.count"))
	// grepFlag, _ := cmd.Flags().GetString("grep")
	outFlag, _ := cmd.Flags().GetString("out")

	/*
		Create the output directory
	*/
	if outFlag == "" {
		log.Fatal("Output path is required. Use --out")
	}
	var err error
	outFlag = os.ExpandEnv(outFlag)
	path, err := filepath.Abs(outFlag)
	if err != nil {
		log.Fatalf("Failed to get absolute path of %s: %v", outFlag, err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		log.Fatalf("Output directory %s already exists", path)
	}
	err = os.Mkdir(path, 0755)
	if err != nil {
		log.Fatalf("Failed to create output directory %s: %v", outFlag, err)
	}

	hub := fctools.NewFarcasterHub()
	defer hub.Close()

	casts := fctools.NewCastGroup().FromCastFidHash(hub, user.Fid, parts[0][2:], expandFlag)
	for _, cast := range casts.Messages {
		s, _ := cast.Json(true, false)
		//fmt.Println(string(s))
		err := os.WriteFile(filepath.Join(path, cast.Hash()+".json"), s, 0644)
		if err != nil {
			log.Fatalf("Failed to write file for cast %s: %v", cast.Hash(), err)
		}
	}
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
	err = os.WriteFile(filepath.Join(path, "urlmap.json"), urlsJson, 0644)
	if err != nil {
		log.Fatalf("Failed to write urlmap.json file: %v", err)
	}
	//fmt.Println(string(urlsJson))

}

func init() {
	rootCmd.AddCommand(snapshotCmd)
	snapshotCmd.Flags().BoolP("recursive", "r", false, "Recursively get parent casts and replies")
	snapshotCmd.Flags().IntP("count", "c", 0, "Number of casts to show when getting @user/casts")
	snapshotCmd.Flags().StringP("grep", "", "", "Only show casts containing a specific string")
	snapshotCmd.Flags().BoolP("json", "", false, "Generate a json object insteead of text")
	snapshotCmd.Flags().BoolP("hex-hashes", "", true, "Used with --json to show hashes in hex")
	snapshotCmd.Flags().BoolP("dates", "", false, "Used with --json to convert fc-timestamps to dates")

	snapshotCmd.Flags().StringP("out", "", "", "Output directory")
}
