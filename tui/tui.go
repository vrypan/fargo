package tui

/*
Functions that display fargo data nicely
*/
import (
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	"github.com/go-color-term/go-color-term/coloring"
	"github.com/muesli/reflow/wordwrap"
	pb "github.com/vrypan/fargo/farcaster"
	"github.com/vrypan/fargo/fctools"
)

const FARCASTER_EPOCH int64 = 1609459200

type fidNames map[uint64]string

func ppTimestamp(ts uint32) string {
	timestamp := time.Unix(int64(ts)+FARCASTER_EPOCH, 0)
	formattedTime := timestamp.Format("2006-01-02 15:04")
	return coloring.Faint("[" + formattedTime + "]")
}
func ppCastId(fname string, hash []byte) string {
	return coloring.Magenta("@"+fname) + coloring.Faint("/"+"0x"+hex.EncodeToString(hash))
}
func ppUrl(url string) string {
	return coloring.Green(url)
}
func addPadding(s string, padding int, paddingString string) string {
	padding_s := strings.Repeat(paddingString, padding)
	lines := strings.Split(strings.TrimSpace(s), "\n")
	for i, line := range lines {
		lines[i] = padding_s + line
	}
	return strings.Join(lines, "\n")
}
func boldBlock(s string) string {
	sb := strings.Builder{}
	lines := strings.Split(strings.TrimSpace(s), "\n")
	for _, line := range lines {
		sb.WriteString(coloring.Bold(line))
		sb.WriteString("\n")
	}
	return strings.TrimSuffix(sb.String(), "\n")
}

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

func FormatCast(msg *pb.Message, fnames map[uint64]string, padding int, showInReply bool, highlight string, grep string) string {

	body := pb.CastAddBody(*msg.Data.GetCastAddBody())

	var builder strings.Builder
	var ptr uint32 = 0
	for i, mention := range body.Mentions {
		builder.WriteString(body.Text[ptr:body.MentionsPositions[i]] + "@" + fnames[mention])
		ptr = body.MentionsPositions[i]
	}
	builder.WriteString(body.Text[ptr:])
	textBody := wordwrap.String(builder.String(), 79)

	builder.Reset()
	builder.WriteString(ppCastId(fnames[msg.Data.Fid], msg.Hash))
	builder.WriteString(" ")
	builder.WriteString(ppTimestamp(msg.Data.Timestamp))
	builder.WriteString("\n")

	if showInReply {
		switch body.GetParent().(type) {
		case *pb.CastAddBody_ParentCastId:
			h := "0x" + hex.EncodeToString(body.GetParentCastId().Hash)
			id := fnames[body.GetParentCastId().Fid]
			builder.WriteString("↳ In reply to @")
			builder.WriteString(id)
			builder.WriteString("/")
			builder.WriteString(h)
			builder.WriteString("\n\n")
		case *pb.CastAddBody_ParentUrl:
			builder.WriteString("↳ In reply to ")
			builder.WriteString(body.GetParentUrl())
			builder.WriteString("\n\n")
		}
	}

	builder.WriteString(textBody)

	if len(body.Embeds) > 0 {
		builder.WriteString("\n----")
	}
	for i, embed := range body.Embeds {
		switch embed.GetEmbed().(type) {
		case *pb.Embed_CastId:
			builder.WriteString("\n[")
			builder.WriteString(strconv.Itoa(i + 1))
			builder.WriteString("] ")
			builder.WriteString(ppCastId(fnames[embed.GetCastId().Fid], embed.GetCastId().Hash))
		case *pb.Embed_Url:
			builder.WriteString("\n[")
			builder.WriteString(strconv.Itoa(i + 1))
			builder.WriteString("] ")
			builder.WriteString(ppUrl(embed.GetUrl()))
		}
	}
	out := builder.String()

	if grep != "" {
		if strings.Contains(out, grep) {
			out = strings.ReplaceAll(out, grep, coloring.Invert(grep))
		} else {
			return ""
		}
	}

	builder.Reset()
	boldFormatting := hex.EncodeToString(msg.Hash) == string(highlight)
	for n, l := range strings.Split(out, "\n") {
		prefix := "│ "
		if n == 0 {
			prefix = "┌─ "
		}
		if boldFormatting {
			builder.WriteString(coloring.Bold(prefix + l + "\n"))
		} else {
			builder.WriteString(prefix + l + "\n")
		}
	}
	if boldFormatting {
		builder.WriteString(coloring.Bold("└───") + "\n")
	} else {
		builder.WriteString("└───\n")
	}
	return addPadding(builder.String(), padding, " ") + "\n"
}

func PprintThread(grp *fctools.CastGroup, hash *fctools.Hash, padding int, hilightHash string, grep string) string {
	if hash == nil {
		hash = &grp.Head
	}
	out := ""
	cast := grp.Messages[*hash].Message
	out += FormatCast(cast, grp.Fnames, padding, (padding == 0), hilightHash, grep)
	for _, reply := range grp.Messages[*hash].Replies {
		out += PprintThread(grp, &reply, padding+4, hilightHash, grep)
	}
	return out
}
func PprintList(grp *fctools.CastGroup, hash *fctools.Hash, padding int, grep string) string {
	out := ""
	for _, cast := range grp.Messages {
		out += FormatCast(cast.Message, grp.Fnames, padding, true, "", grep)
	}
	return out
}
