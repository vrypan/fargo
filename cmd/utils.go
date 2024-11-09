package cmd

import (
	"encoding/hex"
	"log"
	"strings"

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
	user := fctools.NewUser().FromFname(nil, parts[0][1:])

	/*
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
	*/
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
