package fctools

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

	fid_s, err := ldb.Get("FnameFid:" + fname)
	if err == nil {
		return strconv.ParseUint(fid_s, 10, 64)
	} else if err == ldb.ERR_NOT_FOUND {
		hub := NewFarcasterHub()
		defer hub.Close()

		fid, err := hub.GetFidByUsername(fname)
		if err == nil {
			ldb.Set("FnameFid:"+fname, strconv.FormatUint(fid, 10))
		}
		return fid, err
	}

	return 0, err
}
func _print_fid(fid uint64) string {
	fid_s := strconv.FormatUint(fid, 10)
	fname, err := ldb.Get("FidName:" + fid_s)
	if err == ldb.ERR_NOT_FOUND {
		hub := NewFarcasterHub()
		defer hub.Close()
		if fname, err = hub.GetUserData(fid, "USER_DATA_TYPE_USERNAME", false); err == nil {
			ldb.Open()
			defer ldb.Close() // Added defer to ensure the database closes properly
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

func PrintCast(fid uint64, hash string, expand bool, grep string) string {
	ldb.Open()
	defer ldb.Close()

	hash_bytes, err := hex.DecodeString(hash[2:])
	if err != nil {
		return ""
	}
	hub := NewFarcasterHub()
	defer hub.Close()
	return _print_cast(hub, fid, hash_bytes, expand, 0, hash, grep)
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
