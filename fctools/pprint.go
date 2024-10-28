package fctools

import (
        "time"
        "strconv"
        "strings"
        "encoding/hex"

        pb "github.com/vrypan/fargo/farcaster"
        ldb "github.com/vrypan/fargo/localdb"
        "github.com/muesli/reflow/wordwrap"
        "github.com/fatih/color"
)

const FMT_COLS = 80

func _print_fid(fid uint64) string {
	fid_s := strconv.FormatUint(fid, 10)
	fname, err := ldb.Get("FidName:"+fid_s)
	if err == ldb.ERR_NOT_FOUND {
		fname, err = GetUserData(fid, "USER_DATA_TYPE_USERNAME", false)
		if err == nil {
			ldb.Set("FidName:"+fid_s, fname)
		}
	}
	pp := color.New(color.FgMagenta).SprintFunc()
	if len(fname) > 0 {
		return pp("@"+ fname)
	} else {
		return pp("@"+ fid_s)
	}
}
func _print_url(s string) string {
	pp := color.New(color.FgBlue).Add(color.Underline).SprintFunc()
	return pp(s)
}

func formatCastId(fid uint64, hash []byte) string {
	var out string ="" 
	out += _print_fid(fid)
	out += "/0x" + hex.EncodeToString(hash) 
	return out
}

func FormatCast( msg pb.Message ) string {
	var out string

	body := pb.CastAddBody(*msg.Data.GetCastAddBody())

	var ptr uint32 = 0
	for i, mention := range body.Mentions {
		out += body.Text[ptr : body.MentionsPositions[i]] + _print_fid(mention)
		ptr = body.MentionsPositions[i]
	}
	out += body.Text[ptr:] 
 	out = wordwrap.String(out, 79)

	switch body.GetParent().(type) {
		case *pb.CastAddBody_ParentCastId:
			out = "↳ In reply to " + formatCastId(body.GetParentCastId().Fid, body.GetParentCastId().Hash ) + "\n\n" + out
		case *pb.CastAddBody_ParentUrl:
			out = "↳ In reply to " + _print_url(body.GetParentUrl()) + "\n\n" + out
	}

	out = " (" + time.Unix( int64(msg.Data.Timestamp) + FARCASTER_EPOCH, 0).String() + ")\n" + out
	out = formatCastId(msg.Data.Fid, msg.Hash ) + out
	
 	if len(body.Embeds) > 0 {
 		out += "\n----"
 	}
 	for _, embed := range body.Embeds {
		switch embed.GetEmbed().(type) {
			case *pb.Embed_CastId:
				out += "\n* " + formatCastId(embed.GetCastId().Fid, embed.GetCastId().Hash )
			case *pb.Embed_Url:
				out += "\n* " + _print_url(embed.GetUrl())
		}
	}

 	var out2 string = ""
 	for n, l := range strings.Split(out,"\n") {
 		if n == 0 {
 			out2 = "┌─ " + l + "\n"
 		} else {
 			out2 += "│ " + l + "\n"
 		}
 	}
 	out2 += "└───\n"
 	//"➞"
    return out2
}

func PpCastsByFid(fid uint64) (string, error) {
	ldb.Open()
	defer ldb.Close()

	casts, err := GetCastsByFid(fid)
	if err != nil {
		return "", err
	}
    var out string = ""
    for _, m := range casts {
    	out += FormatCast(*m)
    }
    return out, nil
}