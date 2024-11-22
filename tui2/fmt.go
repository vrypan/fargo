package tui2

import (
	"fmt"
	"strings"
	"time"

	gloss "github.com/charmbracelet/lipgloss"
	pb "github.com/vrypan/fargo/farcaster"
)

const FARCASTER_EPOCH int64 = 1609459200

var styleFid = gloss.NewStyle().Foreground(gloss.Color("#b93ec1"))
var styleSecondary = gloss.NewStyle().Foreground(gloss.Color("#777777"))
var styleLink = gloss.NewStyle().Foreground(gloss.Color("#51bf37"))
var styleTextBlock = gloss.NewStyle().Width(80)
var styleFocus = gloss.NewStyle().Bold(true)
var styleNormal = gloss.NewStyle()

var myCuteBorder = gloss.Border{
	Bottom:     "───" + strings.Repeat(" ", 1000),
	Left:       "┌" + strings.Repeat("│", 100),
	BottomLeft: "└───",
}

func tsToDate(ts uint32) string {
	timestamp := time.Unix(int64(ts)+FARCASTER_EPOCH, 0)
	return timestamp.Format("2006-01-02 15:04")
}

func selected(s gloss.Style, flag bool) gloss.Style {
	if flag {
		return gloss.NewStyle().
			Foreground(gloss.Color("#000")).
			Background(s.GetForeground())
	}
	return s
}
func (m *CastsModel) fmtCast(idx int, margin int) string {
	const width = 80
	field := 1

	castHash := m.hashIdx[idx]
	message := m.casts.Messages[castHash].Message
	castAddBody := message.Data.GetCastAddBody()
	fid := message.Data.Fid
	fname := m.casts.Fnames[fid]
	timestamp := message.Data.Timestamp

	maxField := 1 + len(castAddBody.Mentions) + len(castAddBody.Embeds)
	if m.activeField > maxField {
		m.activeField = maxField
	}

	var builder strings.Builder
	builder.WriteString(selected(styleFid, field == m.activeField).Render(fmt.Sprintf("@%s", fname)))
	field++
	builder.WriteString(styleSecondary.Render(fmt.Sprintf("/0x%s [%s]", castHash.String(), tsToDate(timestamp))) + "\n")

	// Check if the cast has a Parent (castId or URL)
	if parent := castAddBody.GetParent(); parent != nil {
		switch parent.(type) {
		case *pb.CastAddBody_ParentCastId:
			id := m.casts.Fnames[castAddBody.GetParentCastId().Fid]
			builder.WriteString(fmt.Sprintf("↳ In reply to @%s/0x%x\n\n", id, castAddBody.GetParentCastId().Hash))
		case *pb.CastAddBody_ParentUrl:
			builder.WriteString(fmt.Sprintf("↳ In reply to %s\n\n", castAddBody.GetParentUrl()))
		}
	}

	// Expand mentions in castText
	var tmp strings.Builder
	ptr := uint32(0)
	for i, mention := range castAddBody.Mentions {
		tmp.WriteString(castAddBody.Text[ptr:castAddBody.MentionsPositions[i]] + selected(styleFid, field == m.activeField).Render("@"+m.casts.Fnames[mention]))
		field++
		ptr = castAddBody.MentionsPositions[i]
	}
	tmp.WriteString(castAddBody.Text[ptr:])

	builder.WriteString(styleTextBlock.Render(tmp.String()))

	// Check if the cast has embeds (castId or URL)
	if len(castAddBody.Embeds) > 0 {
		builder.WriteString("\n\n----")
		for i, embed := range castAddBody.Embeds {
			switch embed.GetEmbed().(type) {
			case *pb.Embed_CastId:
				builder.WriteString("\n" +
					selected(styleLink, field == m.activeField).
						Render(fmt.Sprintf("[%d] @%s/%x", i+1, m.casts.Fnames[embed.GetCastId().Fid], embed.GetCastId().Hash)))
			case *pb.Embed_Url:
				builder.WriteString("\n" +
					selected(styleLink, field == m.activeField).Render(fmt.Sprintf("[%d] %s", i+1, embed.GetUrl())))
			}
			field++
		}
	}

	// Wrap the text, add borders, paddings, etc
	return gloss.NewStyle().MarginLeft(margin).
		PaddingLeft(1).PaddingBottom(0).
		BorderStyle(myCuteBorder).
		BorderLeft(true).BorderBottom(true).
		Render(builder.String())
}

/*
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



func (m *CastsModel) expandText(cast *pb.Message) string {
	var builder strings.Builder
	var ptr uint32 = 0
	castBody := cast.Data.GetCastAddBody()
	for i, mention := range castBody.Mentions {
		builder.WriteString(castBody.Text[ptr:castBody.MentionsPositions[i]] + "@" + m.casts.Fnames[mention])
		ptr = castBody.MentionsPositions[i]
	}
	builder.WriteString(castBody.Text[ptr:])
	return builder.String()
}

func (m *CastsModel) fmtCast(cast *pb.Message, padding int, highlight bool, showInReply bool) string {

	body := cast.Data.GetCastAddBody()

	var builder strings.Builder
	var ptr uint32 = 0
	for i, mention := range body.Mentions {
		builder.WriteString(body.Text[ptr:body.MentionsPositions[i]] + "@" + fnames[mention])
		ptr = body.MentionsPositions[i]
	}
	builder.WriteString(body.Text[ptr:])
	textBody := wordwrap.String(builder.String(), opts.Width)

	builder.Reset()
	fname := m.casts.Fnames[cast.Data.Fid]
	hashString := fmt.Sprintf("0x%x", cast.Hash)
	timestamp := time.Unix(int64(cast.Data.Timestamp)+FARCASTER_EPOCH, 0)
	formattedTime := timestamp.Format("2006-01-02 15:04")

	builder.WriteString("┌─ " + "@" + fname + "/" + hashString)
	builder.WriteString(" ")
	builder.WriteString(formattedTime)
	builder.WriteString("\n")

	if showInReply {
		switch body.GetParent().(type) {
		case *pb.CastAddBody_ParentCastId:
			h := "0x" + hex.EncodeToString(body.GetParentCastId().Hash)
			id := m.casts.Fnames[body.GetParentCastId().Fid]
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
			builder.WriteString(
				ppCastId(m.casts.Fnames[embed.GetCastId().Fid], embed.GetCastId().Hash),
			)
		case *pb.Embed_Url:
			builder.WriteString("\n[")
			builder.WriteString(strconv.Itoa(i + 1))
			builder.WriteString("] ")
			builder.WriteString(ppUrl(embed.GetUrl()))
		}
	}
	out := builder.String()
	builder.Reset()
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

func (m *CastsModel) fmtCast2(cast *pb.Message, padding int, highlight bool) string {
	colorSpecial := lipgloss.Color("#900")
	colorSpecial2 := lipgloss.Color("#090")
	colorDim := lipgloss.Color("#777")
	colorNormal := lipgloss.Color("#aaa")

	textBorder := lipgloss.Border{Left: "│", BottomLeft: "└", Bottom: "───", TopLeft: "┌", Top: " "}

	fnameStyle := lipgloss.NewStyle().Foreground(colorSpecial).PaddingLeft(1)
	hashStyle := lipgloss.NewStyle().Foreground(colorDim)
	//textStyle := lipgloss.NewStyle().Foreground(colorNormal).BorderLeft(true).BorderStyle()
	linkStyle := lipgloss.NewStyle().Foreground(colorSpecial2)

	/*
		textBlockStyle := lipgloss.NewStyle().
			MarginLeft(padding).
			Border(textBorder, false, false, false, true).
			PaddingLeft(1).
			Foreground(colorNormal)
		footerBlockStyle := lipgloss.NewStyle().
			MarginLeft(padding).
			Border(textBorder, false, false, true, true).
			PaddingLeft(1).
			Foreground(colorNormal)
*/
//cast := m.casts.Messages[m.hashIdx[index]].Message
/*
	fname := m.casts.Fnames[cast.Data.Fid]
	body := cast.Data.GetCastAddBody()

	hashString := fmt.Sprintf("0x%x", cast.Hash) // Convert []byte to hex string
	contentHead := "┌─" + fnameStyle.Render("@"+fname) + "/" + hashStyle.Render(hashString)
	contentBody := lipgloss.NewStyle().Foreground(colorNormal).MarginLeft(1).Render(
		wordwrap.String(m.expandText(cast), 80),
	)

	contentEmbeds := ""
	if len(body.Embeds) > 0 {
		contentEmbeds += lipgloss.NewStyle().Foreground(colorNormal).Render("----") + "\n"
		for i, embed := range body.Embeds {
			switch embed.GetEmbed().(type) {
			case *pb.Embed_CastId:
				contentEmbeds += "\n" + lipgloss.NewStyle().Foreground(colorNormal).Render(
					fmt.Sprintf("│ [%d] @%s/%x",
						i+1,
						m.casts.Fnames[embed.GetCastId().Fid], embed.GetCastId().Hash),
				)
			case *pb.Embed_Url:
				contentEmbeds += "\n" + linkStyle.Render(fmt.Sprintf("│ [%d] %s", i+1, embed.GetUrl()))
			}
		}
	}

	return contentHead + "\n" +
		contentBody + "\n" + contentEmbeds + "\n└───\n\n"
}
*/
