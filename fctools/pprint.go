package fctools

import (
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	"github.com/muesli/reflow/wordwrap"
	pb "github.com/vrypan/fargo/farcaster"
)

/*
import (

	"encoding/hex"
	"strconv"
	"strings"
	"time"

	"github.com/go-color-term/go-color-term/coloring"
	"github.com/muesli/reflow/wordwrap"
	pb "github.com/vrypan/fargo/farcaster"
	ldb "github.com/vrypan/fargo/localdb"

)

	func addPadding(s string, padding int) string {
		padding_s := strings.Repeat(" ", padding)
		lines := strings.Split(strings.TrimSpace(s), "\n")
		for i, line := range lines {
			lines[i] = padding_s + line
		}
		return strings.Join(lines, "\n")
	}

	func GetFidByFname(fname string) (uint64, error) {
		ldb.Open()
		defer ldb.Close()

		hub := NewFarcasterHub()
		defer hub.Close()

		return hub.GetFidByUsername(fname)
	}

	func _print_fid(fid uint64) string {
		fid_s := strconv.FormatUint(fid, 10)
		fname, err := ldb.Get("FidName:" + fid_s)
		if err == ldb.ERR_NOT_FOUND {
			hub := NewFarcasterHub()
			defer hub.Close()
			if fname, err = hub.GetUserData(fid, "USER_DATA_TYPE_USERNAME"); err == nil {
				ldb.Set("FidName:"+fid_s, fname)
			}
		}
		if fname != "" {
			return coloring.Magenta("@" + fname)
		}
		return coloring.Magenta("@" + fid_s)
	}

	func _print_timestamp(ts uint32) string {
		timestamp := time.Unix(int64(ts)+FARCASTER_EPOCH, 0)
		formattedTime := timestamp.Format("2006-01-02 15:04")
		return coloring.For("[" + formattedTime + "]").Color(8).String()
	}

	func _print_url(s string) string {
		// pp := color.New(color.FgBlue).Add(color.Underline).SprintFunc()
		return coloring.For(s).Green().Underline().String()
	}

	func FormatCastId(fid uint64, hash []byte, highlight string) string {
		hash_s := "0x" + hex.EncodeToString(hash)
		out := _print_fid(fid)
		colorFunc := coloring.For("/" + hash_s).Color(8).String
		if hash_s == highlight {
			colorFunc = coloring.For("/" + hash_s).Red().String
		}
		return out + colorFunc()
	}

	func FormatCast(msg *pb.Message, padding int, showInReply bool, highlight string, grep string) string {
		var out string

		body := pb.CastAddBody(*msg.Data.GetCastAddBody())

		var ptr uint32 = 0
		for i, mention := range body.Mentions {
			out += body.Text[ptr:body.MentionsPositions[i]] + _print_fid(mention)
			ptr = body.MentionsPositions[i]
		}
		out += body.Text[ptr:]
		out = wordwrap.String(out, 79)

		if showInReply {
			switch body.GetParent().(type) {
			case *pb.CastAddBody_ParentCastId:
				out = "↳ In reply to " + FormatCastId(body.GetParentCastId().Fid, body.GetParentCastId().Hash, highlight) + "\n\n" + out
			case *pb.CastAddBody_ParentUrl:
				out = "↳ In reply to " + _print_url(body.GetParentUrl()) + "\n\n" + out
			}
		}

		out = " " + _print_timestamp(msg.Data.Timestamp) + "\n" + out
		// out = " (" + time.Unix( int64(msg.Data.Timestamp) + FARCASTER_EPOCH, 0).String() + ")\n" + out
		out = FormatCastId(msg.Data.Fid, msg.Hash, highlight) + out

		if len(body.Embeds) > 0 {
			out += "\n----"
		}
		for i, embed := range body.Embeds {
			switch embed.GetEmbed().(type) {
			case *pb.Embed_CastId:
				out += "\n[" + strconv.Itoa(i+1) + "] " + FormatCastId(embed.GetCastId().Fid, embed.GetCastId().Hash, highlight)
			case *pb.Embed_Url:
				out += "\n[" + strconv.Itoa(i+1) + "] " + _print_url(embed.GetUrl())
			}
		}

		var out2 string = ""
		for n, l := range strings.Split(out, "\n") {
			if n == 0 {
				out2 = "┌─ " + l + "\n"
			} else {
				out2 += "│ " + l + "\n"
			}
		}
		out2 += "└───\n"

		if grep == "" {
			return addPadding(out2, padding) + "\n"
		} else {
			if strings.Contains(out2, grep) {
				out2 = strings.ReplaceAll(out2, grep, coloring.Invert(grep))
				return addPadding(out2, padding) + "\n"
			} else {
				return ""
			}
		}
	}

	func PrintCastsByFid(fid uint64, count uint32, grep string) (string, error) {
		ldb.Open()
		defer ldb.Close()
		hub := NewFarcasterHub()
		defer hub.Close()

		casts, err := hub.GetCastsByFid(fid, count)
		if err != nil {
			return "", err
		}

		var builder strings.Builder
		for _, m := range casts {
			builder.WriteString(FormatCast(m, 0, true, "", grep))
		}
		return builder.String(), nil
	}

	func _print_cast(hub *FarcasterHub, fid uint64, hash []byte, expand bool, padding int, highlight string, grep string) string {
		cast, err := hub.GetCast(fid, hash)
		if err != nil {
			panic(err)
		}

		castBody := pb.CastAddBody(*cast.Data.GetCastAddBody())

		// If there's a parent cast and we're expanding from the root
		if castBody.GetParentCastId() != nil && expand && padding == 0 {
			return _print_cast(hub, castBody.GetParentCastId().Fid, castBody.GetParentCastId().Hash, expand, padding, highlight, grep)
		}

		showInReply := padding == 0
		out := FormatCast(cast, padding, showInReply, highlight, grep)

		if expand {
			if casts, err := hub.GetCastReplies(cast.Data.Fid, cast.Hash); err == nil {
				for _, reply := range casts.Messages {
					out += _print_cast(hub, reply.Data.Fid, reply.Hash, true, padding+4, highlight, grep)
				}
			}
		}

		return out
	}
*/
type fidNames map[uint64]string

func ppTimestamp(ts uint32) string {
	timestamp := time.Unix(int64(ts)+FARCASTER_EPOCH, 0)
	formattedTime := timestamp.Format("2006-01-02 15:04")
	return "[" + formattedTime + "]"
}
func ppCastId(fname string, hash []byte) string {
	return "@" + fname + "/" + "0x" + hex.EncodeToString(hash)
}
func ppUrl(url string) string {
	return url
}
func addPadding(s string, padding int) string {
	padding_s := strings.Repeat(" ", padding)
	lines := strings.Split(strings.TrimSpace(s), "\n")
	for i, line := range lines {
		lines[i] = padding_s + line
	}
	return strings.Join(lines, "\n")
}

func FormatCast(msg *pb.Message, fnames map[uint64]string, padding int, showInReply bool, highlight string, grep string) string {
	var out string

	body := pb.CastAddBody(*msg.Data.GetCastAddBody())

	var ptr uint32 = 0
	for i, mention := range body.Mentions {
		out += body.Text[ptr:body.MentionsPositions[i]] + "@" + fnames[mention]
		ptr = body.MentionsPositions[i]
	}
	out += body.Text[ptr:]
	out = wordwrap.String(out, 79)

	if showInReply {
		switch body.GetParent().(type) {
		case *pb.CastAddBody_ParentCastId:
			h := "0x" + hex.EncodeToString(body.GetParentCastId().Hash)
			id := fnames[body.GetParentCastId().Fid]
			out = "↳ In reply to @" + id + "/" + h + "\n\n" + out
		case *pb.CastAddBody_ParentUrl:
			out = "↳ In reply to " + body.GetParentUrl() + "\n\n" + out
		}
	}

	out = " " + ppTimestamp(msg.Data.Timestamp) + "\n" + out
	out = ppCastId(fnames[msg.Data.Fid], msg.Hash) + out

	if len(body.Embeds) > 0 {
		out += "\n----"
	}
	for i, embed := range body.Embeds {
		switch embed.GetEmbed().(type) {
		case *pb.Embed_CastId:
			out += "\n[" + strconv.Itoa(i+1) + "] " + ppCastId(fnames[embed.GetCastId().Fid], embed.GetCastId().Hash)
		case *pb.Embed_Url:
			out += "\n[" + strconv.Itoa(i+1) + "] " + ppUrl(embed.GetUrl())
		}
	}

	var out2 string = ""
	for n, l := range strings.Split(out, "\n") {
		if n == 0 {
			out2 = "┌─ " + l + "\n"
		} else {
			out2 += "│ " + l + "\n"
		}
	}
	out2 += "└───\n"

	return addPadding(out2, padding) + "\n"
}

/*
func PrintCast(fid uint64, hash string, expand bool, grep string) string {
	db.Open()
	defer db.Close()

	hash_bytes, err := hex.DecodeString(hash[2:])
	if err != nil {
		return ""
	}
	hub := NewFarcasterHub()
	defer hub.Close()

	castMessage, err := hub.PrxGetCast(fid, hash_bytes)
	if err != nil {
		return ""
	}

	fnames := make(fidNames)
	fnames[fid], _ = hub.PrxGetUserData(fid, "USER_DATA_TYPE_USERNAME")
	for _, mention := range castMessage.Data.GetCastAddBody().Mentions {
		fnames[mention], _ = hub.PrxGetUserData(mention, "USER_DATA_TYPE_USERNAME")
	}
	if castMessage.Data.GetCastAddBody().GetParentCastId() != nil {
		p_cast_fid := castMessage.Data.GetCastAddBody().GetParentCastId().Fid
		p_cast_fname, _ := hub.PrxGetUserData(p_cast_fid, "USER_DATA_TYPE_USERNAME")
		fnames[p_cast_fid] = p_cast_fname
	}
	return FormatCast(castMessage, fnames, 0, true, "", "")
}
*/
