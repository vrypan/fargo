package fctools

import (
        "encoding/hex"
        "log"
        pb "github.com/vrypan/fargo/farcaster"
        ldb "github.com/vrypan/fargo/localdb"
)

func GetFidUrls(fid uint64, count uint32, grep string) []string {
	ldb.Open()
	defer ldb.Close()

	hub := NewFarcasterHub(); defer hub.Close()
	
	casts, err := hub.GetCastsByFid(fid, count)
	if err != nil {
		return nil
	}
    
    var embed_list []string
    for _, m := range casts {
    	cast_body := pb.CastAddBody(*m.Data.GetCastAddBody())
    	for _, embed := range cast_body.Embeds {
			if ebd := embed.GetUrl(); ebd != "" {
				// Should check for mime/type here.
				embed_list = append(embed_list, ebd)
			}
		}
    }
    return embed_list
}

func GetCastUrls(fid uint64, hash string, expand bool, grep string) []string {
	ldb.Open()
	defer ldb.Close()

	hash_bytes, err := hex.DecodeString(hash[2:])
	if err != nil {
		return nil
	}
	hub := NewFarcasterHub(); defer hub.Close()
	return _get_cast_urls( hub, fid, hash_bytes, expand, expand, grep )
}

func _get_cast_urls(hub *FarcasterHub, fid uint64, hash []byte, expand_up bool, expand_down bool, grep string) []string {
	var embed_list []string
	cast, e := hub.GetCast(fid, hash)
	if e != nil {
		log.Fatal(e)
	}

	cast_body := pb.CastAddBody(*cast.Data.GetCastAddBody())
	
	if cast_body.GetParentCastId() != nil && expand_up {
		// If the cast has a parent, start from the partent and keep going up (expand_up==true, expand_down=false)
		return _get_cast_urls(hub, cast_body.GetParentCastId().Fid, cast_body.GetParentCastId().Hash, true, false, grep )
	}
	
	// Should check for grep here
 	for _, embed := range cast_body.Embeds {
		if ebd := embed.GetUrl(); ebd != "" {
			// Should check for mime/type here.
			embed_list = append(embed_list, ebd)
		}
	}
	if expand_down {
		casts, err := hub.GetCastReplies( cast.Data.Fid, cast.Hash )
		if err == nil {
			for _, c := range casts.Messages {
				embed_list = append(embed_list, _get_cast_urls( hub, c.Data.Fid, c.Hash, false, true, grep )...)
			}
		}
	}
	return embed_list
}