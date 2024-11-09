package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/vrypan/fargo/localdb"
)

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Cache management",
	Long: `By default, all records written in the local database are pruned
after 24 hours. Set db.ttlhours to a value that suits you better
using "fargo set".

Note that changing db.ttlhours will only affect new entries.

You can also delete the database path, and it will be re-created
next time you run fargo.`,
	Run: cacheManager,
}

func cacheManager(cmd *cobra.Command, args []string) {
	bytesOnDisk, err := localdb.GetSize()
	if err != nil {
		log.Fatal("Error reading files.", err)
	}
	humanReadableSize := fmt.Sprintf("%.2f MB", float64(bytesOnDisk)/1024/1024)

	if err != nil {
		log.Fatal("Error reading files.", err)
	}
	localdb.Open()
	defer localdb.Close()
	entriesCount, err := localdb.CountEntries()
	if err != nil {
		log.Fatal("Error accessing the database.", err)
	}
	fmt.Printf("DB path: %s\n", localdb.Path())
	fmt.Printf("Database size on disk: %s\n", humanReadableSize)
	fmt.Printf("Number of entries: %d\n", entriesCount)
}

func init() {
	rootCmd.AddCommand(cacheCmd)
}
