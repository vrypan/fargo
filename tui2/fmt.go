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
