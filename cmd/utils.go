package cmd

import (
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"os/exec"
	"runtime"

	"github.com/vrypan/fargo/fctools"
)

func parse_url(args []string) (*fctools.User, []string) {
	if len(args) == 0 {
		log.Fatal("No path")
	}
	return ParseFcURI(args[0])
}

func ParseFcURI(uri string) (*fctools.User, []string) {
	parts := strings.Split(uri, "/")

	if parts[0][0:1] != "@" {
		log.Fatal("Path should start with @")
	}
	user, err := fctools.NewUser().FromFname(nil, parts[0][1:])
	if err != nil {
		log.Fatal(err)
	}
	return user, parts[1:]
}

// Convert "0xhash" to []byte
func HashToBytes(hash string) []byte {
	if hash_bytes, err := hex.DecodeString(hash[2:]); err != nil {
		return nil
	} else {
		return hash_bytes
	}
}
func OpenUrl(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}
