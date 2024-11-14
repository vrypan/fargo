package tui

import (
	"strings"

	"github.com/vrypan/fargo/fctools"
)

var reactionVerb = map[string]string{
	"REACTION_TYPE_LIKE":   "liked",
	"REACTION_TYPE_RECAST": "recasted",
}

func PpReactionsList(reactions *fctools.Reactions, casts *fctools.CastGroup, opts *FmtCastOpts) string {
	var builder strings.Builder
	for _, r := range reactions.Messages {

		castHash := r.Message.Data.GetReactionBody().GetTargetCastId().Hash

		reactionString := reactionVerb[r.Message.Data.GetReactionBody().Type.String()]
		builder.WriteString("\n")

		if cast, ok := casts.Messages[fctools.Hash(castHash)]; ok {
			opts.Prepend = "┌─ " + PpFname(reactions.Fnames[r.Message.Data.Fid]) +
				" " + reactionString + " " + ppTimestamp(r.Message.Data.Timestamp) +
				"\n│\n│─ "
			builder.WriteString(
				FmtCast(cast.Message, casts.Fnames, 0, true, opts),
			)
		}
	}
	return builder.String()
}
