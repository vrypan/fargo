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
	s = strings.TrimSpace(s)
	lines := strings.Split(s, "\n")
	for i := range lines {
		lines[i] = padding_s + lines[i]
	}
	return strings.Join(lines, "\n")
}

func GetFidByFname(fname string) (uint64, error) {
	ldb.Open()
	defer ldb.Close()
	var fid uint64

	fid_s, err := ldb.Get("FnameFid:" + fname)
	if err == ldb.ERR_NOT_FOUND {
		hub := NewFarcasterHub()
		defer hub.Close()
		fid, err = hub.GetFidByUsername(fname)
		if err == nil {
			ldb.Set("FnameFid:"+fname, strconv.FormatUint(fid, 10))
		}
		return fid, nil
	} else {
		fid, _ = strconv.ParseUint(fid_s, 10, 64)
		return fid, nil
	}
}
func _print_fid(fid uint64) string {
	fid_s := strconv.FormatUint(fid, 10)
	fname, err := ldb.Get("FidName:" + fid_s)
	if err == ldb.ERR_NOT_FOUND {
		hub := NewFarcasterHub()
		defer hub.Close()
		fname, err = hub.GetUserData(fid, "USER_DATA_TYPE_USERNAME", false)
		if err == nil {
			ldb.Set("FidName:"+fid_s, fname)
		}
	}
	if len(fname) > 0 {
		return coloring.Magenta("@" + fname)
	} else {
		return coloring.Magenta("@" + fid_s)
	}
}
func _print_timestamp(ts uint32) string {
	ret := "[" + time.Unix(int64(ts)+FARCASTER_EPOCH, 0).Format("2006-01-02 15:04") + "]"
	return coloring.For(ret).Color(8).String()
}
func _print_url(s string) string {
	// pp := color.New(color.FgBlue).Add(color.Underline).SprintFunc()
	return coloring.For(s).Green().Underline().String()
}

func formatCastId(fid uint64, hash []byte, highlight string) string {
	var out string = ""
	hash_s := "0x" + hex.EncodeToString(hash)
	out += _print_fid(fid)
	if hash_s == highlight {
		out += coloring.For("/" + hash_s).Red().String()
	} else {
		out += coloring.For("/" + hash_s).Color(8).String()
	}
	return out
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
			out = "↳ In reply to " + formatCastId(body.GetParentCastId().Fid, body.GetParentCastId().Hash, highlight) + "\n\n" + out
		case *pb.CastAddBody_ParentUrl:
			out = "↳ In reply to " + _print_url(body.GetParentUrl()) + "\n\n" + out
		}
	}

	out = " " + _print_timestamp(msg.Data.Timestamp) + "\n" + out
	// out = " (" + time.Unix( int64(msg.Data.Timestamp) + FARCASTER_EPOCH, 0).String() + ")\n" + out
	out = formatCastId(msg.Data.Fid, msg.Hash, highlight) + out

	if len(body.Embeds) > 0 {
		out += "\n----"
	}
	for _, embed := range body.Embeds {
		switch embed.GetEmbed().(type) {
		case *pb.Embed_CastId:
			out += "\n* " + formatCastId(embed.GetCastId().Fid, embed.GetCastId().Hash, highlight)
		case *pb.Embed_Url:
			out += "\n* " + _print_url(embed.GetUrl())
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
	var out string = ""
	for _, m := range casts {
		out += FormatCast(m, 0, true, "", grep)
	}
	return out, nil
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
	cast, e := hub.GetCast(fid, hash)
	if e != nil {
		panic(e)
	}

	cast_body := pb.CastAddBody(*cast.Data.GetCastAddBody())

	if cast_body.GetParentCastId() != nil && expand && padding == 0 {
		return _print_cast(hub, cast_body.GetParentCastId().Fid, cast_body.GetParentCastId().Hash, true, 0, highlight, grep)
	}

	showInReply := false
	if padding == 0 {
		showInReply = true
	}
	out := FormatCast(cast, padding, showInReply, highlight, grep)
	if expand {
		casts, err := hub.GetCastReplies(cast.Data.Fid, cast.Hash)
		if err == nil {
			for _, c := range casts.Messages {
				out += _print_cast(hub, c.Data.Fid, c.Hash, true, padding+4, highlight, grep)
			}
		}
	}
	return out
}
