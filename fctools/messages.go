package fctools

import (
	"crypto/ed25519"
	"regexp"
	"strconv"
	"strings"

	pb "github.com/vrypan/fargo/farcaster"
	"github.com/zeebo/blake3"
	"google.golang.org/protobuf/proto"
)

func ProcessCastBody(text string) (string, []uint32, []uint64, []*pb.Embed, string) {
	var (
		mentionPositions []uint32
		mentions         []uint64
		embeds           []*pb.Embed
		resultText       string
		currentIndex     int
		offset           int
		embedCount       int
	)

	urlRe := regexp.MustCompile(`^\[(http[s]?://(?:[a-zA-Z]|[0-9]|[$-_@.&+]|[!*\\(\\),]|(?:%[0-9a-fA-F][0-9a-fA-F]))+)\]$`)
	fnameRe := regexp.MustCompile(`^@([a-z0-9][a-z0-9-]{0,15})((\.eth)?)(\S*)`)

	lines := strings.Split(text, "\n")
	for lIdx, line := range lines {
		words := strings.Fields(line)
		for wIdx, word := range words {
			switch {
			case fnameRe.MatchString(word):
				if matched := fnameRe.FindStringSubmatch(word); matched != nil {
					if fid, err := GetFidByFname(matched[1] + matched[2]); err == nil {
						if len(resultText+" "+matched[4]) > 1024 {
							return formMoreText(word, words[wIdx+1:], lines[lIdx+1:]), mentionPositions, mentions, embeds, ""
						}
						mentionPositions = append(mentionPositions, uint32(currentIndex))
						mentions = append(mentions, fid)
						offset += len(matched[4]) + 1
					}
				}
			case urlRe.MatchString(word):
				if matched := urlRe.FindStringSubmatch(word); matched != nil {
					if len(resultText+"["+strconv.Itoa(embedCount+1)+"]") > 1024 {
						return formMoreText(word, words[wIdx+1:], lines[lIdx+1:]), mentionPositions, mentions, embeds, ""
					}
					embeds = append(embeds, &pb.Embed{
						Embed: &pb.Embed_Url{Url: matched[1]},
					})
					if resultText != "" {
						resultText += " "
					}
					resultText += "[" + strconv.Itoa(embedCount+1) + "]"
					offset += 4
					embedCount++
				}
			default:
				if len(resultText+" "+word) > 1024 {
					return formMoreText(word, words[wIdx+1:], lines[lIdx+1:]), mentionPositions, mentions, embeds, ""
				}
				if resultText != "" && wIdx > 0 {
					resultText += " "
				}
				resultText += word
				offset += len(word) + 1
			}
			currentIndex = offset
		}
		resultText += "\n"
		offset++
	}
	return resultText, mentionPositions, mentions, embeds, ""
}

func formMoreText(word string, remainingWords []string, remainingLines []string) string {
	more := word
	for _, w := range remainingWords {
		more += " " + w
	}
	for _, l := range remainingLines {
		more += "\n" + l
	}
	return more
}

func CreateMessage(messageData *pb.MessageData, signerPrivate []byte, signerPublic []byte) *pb.Message {
	hashScheme := pb.HashScheme(pb.HashScheme_value["HASH_SCHEME_BLAKE3"])
	signatureScheme := pb.SignatureScheme(pb.SignatureScheme_value["SIGNATURE_SCHEME_ED25519"])
	dataBytes, _ := proto.Marshal(messageData)
	signerCombined := append(signerPrivate, signerPublic...)

	hasher := blake3.New()
	hasher.Write(dataBytes)
	hash := hasher.Sum(nil)[:20]

	signature := ed25519.Sign(signerCombined, hash)

	return &pb.Message{
		Data:            messageData,
		Hash:            hash,
		HashScheme:      hashScheme,
		Signature:       signature,
		SignatureScheme: signatureScheme,
		Signer:          signerPublic,
		DataBytes:       dataBytes,
	}
}
