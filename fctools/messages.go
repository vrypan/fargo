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
	var mentionPositions []uint32
	var mentions []uint64
	var embeds []*pb.Embed
	var resultText string
	var currentIndex int
	var offset int

	urlRe := regexp.MustCompile(`^\[(http[s]?://(?:[a-zA-Z]|[0-9]|[$-_@.&+]|[!*\\(\\),]|(?:%[0-9a-fA-F][0-9a-fA-F]))+)\]$`)
	fnameRe := regexp.MustCompile(`^@([a-z0-9][a-z0-9-]{0,15})((\.eth)?)(\S*)`)
	var embedCount int

	lines := strings.Split(text, "\n")
	for lIdx, line := range lines {
		words := strings.Fields(line)
		for wIdx, word := range words {
			switch {
			case fnameRe.MatchString(word):
				if matched := fnameRe.FindStringSubmatch(word); matched != nil {
					if fid, err := GetFidByFname(matched[1] + matched[2]); err == nil {
						if len(resultText+" "+matched[4]) > 1024 {
							more := word
							for _, w := range words[wIdx+1:] {
								more += " " + w
							}
							for _, l := range lines[lIdx+1:] {
								more += "\n" + l
							}
							return resultText, mentionPositions, mentions, embeds, more
						}
						mentionPositions = append(mentionPositions, uint32(currentIndex))
						mentions = append(mentions, fid)
						offset += len(matched[4]) + 1
					}
				}
			case urlRe.MatchString(word):
				if matched := urlRe.FindStringSubmatch(word); matched != nil {
					if len(resultText+"["+strconv.Itoa(embedCount+1)+"]") > 1024 {
						more := word
						for _, w := range words[wIdx+1:] {
							more += " " + w
						}
						for _, l := range lines[lIdx+1:] {
							more += "\n" + l
						}
						return resultText, mentionPositions, mentions, embeds, more
					}
					embeds = append(embeds, &pb.Embed{
						Embed: &pb.Embed_Url{Url: matched[1]},
					})
					if resultText != "" {
						resultText += " " + "[" + strconv.Itoa(embedCount+1) + "]"
					} else {
						resultText += "[" + strconv.Itoa(embedCount+1) + "]"
					}
					offset += 4
					embedCount++
				}
			default:
				if len(resultText+" "+word) > 1024 {
					more := word
					for _, w := range words[wIdx+1:] {
						more += " " + w
					}
					for _, l := range lines[lIdx+1:] {
						more += "\n" + l
					}
					return resultText, mentionPositions, mentions, embeds, more
				}
				if resultText != "" && wIdx > 0 {
					resultText += " " + word
				} else {
					resultText += word
				}
				offset += len(word) + 1
			}
			currentIndex = offset
		}
		resultText += "\n"
		offset++
	}
	return resultText, mentionPositions, mentions, embeds, ""
}

func CreateMessage(messageData *pb.MessageData, signerPrivate []byte, signerPublic []byte) *pb.Message {
	hash_scheme := pb.HashScheme(pb.HashScheme_value["HASH_SCHEME_BLAKE3"])
	signature_scheme := pb.SignatureScheme(pb.SignatureScheme_value["SIGNATURE_SCHEME_ED25519"])
	data_bytes, _ := proto.Marshal(messageData)
	signerPublic_ := append(signerPrivate, signerPublic...) // required by ed25519 Go implementation

	hasher := blake3.New()
	hasher.Write(data_bytes)
	hash := hasher.Sum(nil)[0:20]

	signature := ed25519.Sign(signerPublic_, hash)

	return &pb.Message{
		Data:            messageData,
		Hash:            hash,
		HashScheme:      hash_scheme,
		Signature:       signature,
		SignatureScheme: signature_scheme,
		Signer:          signerPublic,
		DataBytes:       data_bytes,
	}
}
