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

type FmtCastOpts struct {
	Grep      string
	Highlight string
	Prepend   string
	Append    string
	Width     int
}

type fidNames map[uint64]string

func ppTimestamp(ts uint32) string {
	timestamp := time.Unix(int64(ts)+FARCASTER_EPOCH, 0)
	formattedTime := timestamp.Format("2006-01-02 15:04")
	return coloring.Faint("[" + formattedTime + "]")
}
func PpFname(fname string) string {
	return coloring.Magenta("@" + fname)
}
func ppCastId(fname string, hash []byte) string {
	return PpFname(fname) + coloring.Faint("/"+"0x"+hex.EncodeToString(hash))
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

	if grep != "" && !strings.Contains(builder.String(), grep) {
		return ""
	}
	out = builder.String()
	if grep != "" {
		out = strings.ReplaceAll(out, grep, coloring.Invert(grep))
	}
	return addPadding(out, padding, " ") + "\n"
}

func PprintThread(grp *fctools.CastGroup, hash *fctools.Hash, padding int, hilightHash string, grep string) string {
	if hash == nil {
		hash = &grp.Head
	}
	out := ""
	var cast *pb.Message
	if msg, ok := grp.Messages[*hash]; ok {
		cast = msg.Message
	} else {
		return ""
	}
	out += FormatCast(cast, grp.Fnames, padding, (padding == 0), hilightHash, grep)
	for _, reply := range grp.Messages[*hash].Replies {
		out += PprintThread(grp, &reply, padding+4, hilightHash, grep)
	}
	return out
}
func PprintCastList(grp *fctools.CastGroup, hash *fctools.Hash, padding int, grep string) string {
	out := ""
	for _, cast := range grp.Messages {
		out += FmtCast(cast.Message, grp.Fnames, padding, true, &FmtCastOpts{Grep: grep, Highlight: "", Width: 50})
	}
	return out
}

func FmtCast(
	msg *pb.Message,
	fnames map[uint64]string,
	padding int,
	showInReply bool,
	opts *FmtCastOpts,
) string {

	body := pb.CastAddBody(*msg.Data.GetCastAddBody())

	var builder strings.Builder
	var ptr uint32 = 0
	for i, mention := range body.Mentions {
		builder.WriteString(body.Text[ptr:body.MentionsPositions[i]] + "@" + fnames[mention])
		ptr = body.MentionsPositions[i]
	}
	builder.WriteString(body.Text[ptr:])
	textBody := wordwrap.String(builder.String(), opts.Width)

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

	builder.Reset()
	boldFormatting := "0x"+hex.EncodeToString(msg.Hash) == string(opts.Highlight)
	for n, l := range strings.Split(out, "\n") {
		prefix := "│ "
		if n == 0 {
			if opts.Prepend != "" {
				prefix = opts.Prepend
			} else {
				prefix = "┌─ "
			}
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

	if opts.Grep != "" && !strings.Contains(builder.String(), opts.Grep) {
		return ""
	}
	if opts.Grep != "" {
		out = strings.ReplaceAll(builder.String(), opts.Grep, coloring.Invert(opts.Grep))
	} else {
		out = builder.String()
	}
	out = out + opts.Append
	return addPadding(out, padding, " ") + "\n"
}
