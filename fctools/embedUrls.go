package fctools

import (
	"encoding/hex"
	"log"

	pb "github.com/vrypan/fargo/farcaster"
	"github.com/vrypan/fargo/fctools/embedurl"
	ldb "github.com/vrypan/fargo/localdb"
)

func GetFidUrls(fid uint64, count uint32, grep string) []embedurl.Url {
	ldb.Open()
	defer ldb.Close()

	hub := NewFarcasterHub()
	defer hub.Close()

	casts, err := hub.GetCastsByFid(fid, count)
	if err != nil {
		return nil
	}

	var embed_list []embedurl.Url
	for _, cast := range casts {
		embed_list = append(embed_list, embedurl.FromMessage(cast)...)
	}
	return embed_list
}

func GetCastUrls(fid uint64, hash string, expand bool, grep string) []embedurl.Url {
	ldb.Open()
	defer ldb.Close()

	hash_bytes, err := hex.DecodeString(hash[2:])
	if err != nil {
		return nil
	}
	hub := NewFarcasterHub()
	defer hub.Close()
	return _get_cast_urls(hub, fid, hash_bytes, expand, expand, grep)
}

func _get_cast_urls(hub *FarcasterHub, fid uint64, hash []byte, expand_up bool, expand_down bool, grep string) []embedurl.Url {
	var embed_list []embedurl.Url
	cast, e := hub.GetCast(fid, hash)
	if e != nil {
		log.Fatal(e)
	}

	cast_body := pb.CastAddBody(*cast.Data.GetCastAddBody())

	if cast_body.GetParentCastId() != nil && expand_up {
		// If the cast has a parent, start from the partent and keep going up (expand_up==true, expand_down=false)
		return _get_cast_urls(hub, cast_body.GetParentCastId().Fid, cast_body.GetParentCastId().Hash, true, false, grep)
	}

	embed_list = append(embed_list, embedurl.FromMessage(cast)...)

	if expand_down {
		casts, err := hub.GetCastReplies(cast.Data.Fid, cast.Hash)
		if err == nil {
			for _, c := range casts.Messages {
				embed_list = append(embed_list, _get_cast_urls(hub, c.Data.Fid, c.Hash, false, true, grep)...)
			}
		}
	}
	return embed_list
}
