package fctools

import (
	"regexp"
	"strings"

	pb "github.com/vrypan/fargo/farcaster"
)

func ProcessCastBody(text string) (string, []uint32, []uint64, []*pb.Embed) {
	var mentionPositions []uint32
	var mentions []uint64
	var embeds []*pb.Embed
	var resultText string
	var currentIndex int
	var offset int

	lines := strings.Split(text, "\n")
	for _, line := range lines {
		words := strings.Fields(line)
		for _, word := range words {
			fnameRe := regexp.MustCompile(`^@([a-z0-9][a-z0-9-]{0,15})((\.eth)?)(\S*)`)
			if matched := fnameRe.FindStringSubmatch(word); matched != nil {
				fid, err := GetFidByFname(matched[1] + matched[2])
				if err == nil {
					mentionPositions = append(mentionPositions, uint32(currentIndex))
					mentions = append(mentions, fid)
					// Remove the @fname mention from result text
					resultText += " " + matched[4]
					offset += len(matched[4]) + 1
				}
			} else if strings.HasPrefix(word, "http://") || strings.HasPrefix(word, "https://") {
				embeds = append(embeds, &pb.Embed{
					Embed: &pb.Embed_Url{Url: word},
				})
				if resultText != "" {
					resultText += " "
				}
				resultText += word
				offset += len(word) + 1
			} else {
				if resultText != "" {
					resultText += " "
				}
				resultText += word
				offset += len(word) + 1
			}
			currentIndex = offset
		}
		// Preserve newline at the end of each line
		resultText += "\n"
		offset++
	}
	return resultText, mentionPositions, mentions, embeds
}
