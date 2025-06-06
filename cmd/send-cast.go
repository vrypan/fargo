package cmd

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	pb "github.com/vrypan/farcaster-go/farcaster"
	"github.com/vrypan/fargo/config"
	"github.com/vrypan/fargo/fctools"
	db "github.com/vrypan/fargo/localdb"
)

var sendCastCmd = &cobra.Command{
	Use:   "cast [text]",
	Short: "Posts new cast",
	Long: `[text] is the full cast text.
Any @mentions will be identified automatically and (up to two)
links enclosed in brackets will be converted to embeds.

If the text is longer than 1024 bytes, it will be broken down
to multiple casts posted as a thread.`,
	Run: runSendCast,
}

func runSendCast(cmd *cobra.Command, args []string) {
	//config.Load()
	db.Open()
	defer db.Close()

	var err error
	var privateKey []byte
	// var publicKey []byte
	var s string

	var castMessageBodies []*pb.CastAddBody

	s = config.GetString("cast.privkey")
	if c, _ := cmd.Flags().GetString("privkey"); c != "" {
		s = c
	}
	if len(s) < 2 {
		log.Fatal("Private key error: private key string too short")
	}
	if privateKey, err = hex.DecodeString(s[2:]); err != nil {
		log.Fatalf("Private key error: %v\nUse --help to see options.", err)
	}

	expandedKey := ed25519.NewKeyFromSeed(privateKey)
	publicKey := expandedKey.Public().(ed25519.PublicKey)

	fid := uint64(config.GetInt("cast.fid"))
	if c, _ := cmd.Flags().GetUint64("fid"); c > 0 {
		fid = c
	}
	if fid == 0 {
		log.Fatal("No fid: fid is zero. Use --help to see options.")
	}

	replyToCastFlag, _ := cmd.Flags().GetString("reply-to-cast")
	replyToUrlFlag, _ := cmd.Flags().GetString("reply-to-url")

	prepareFlag, _ := cmd.Flags().GetBool("prepare")

	if len(args) == 0 {
		log.Fatal("Missing arguments: text argument required")
	}

	hub := fctools.NewFarcasterHub()
	defer hub.Close()

	more := args[0]
	var (
		castText          string
		mentionsPositions []uint32
		mentions          []uint64
		embeds            []*pb.Embed
		castType          uint8
	)
	for { // Cast storm!!!
		castText, mentionsPositions, mentions, embeds, more = fctools.ProcessCastBody(more)
		//castText = strings.TrimSpace(castText)
		//more = strings.TrimSpace(more)
		if len(embeds) > 2 {
			embeds = embeds[0:2]
		}
		if len(castText) <= 320 {
			castType = 0
		} else {
			castType = 1
		}
		messageBody := &pb.CastAddBody{
			Mentions:          mentions,
			MentionsPositions: mentionsPositions,
			Text:              castText,
			Type:              pb.CastType(castType),
			Embeds:            embeds,
		}
		castMessageBodies = append(castMessageBodies, messageBody)
		if more == "" {
			break
		}
	}
	for _, messageBody := range castMessageBodies {
		if replyToCastFlag != "" {
			parent, parentHashString := ParseFcURI(replyToCastFlag)
			parentHash := HashToBytes(parentHashString[0])
			parentCast := &pb.CastAddBody_ParentCastId{ParentCastId: &pb.CastId{Fid: parent.Fid, Hash: parentHash}}
			messageBody.Parent = parentCast
		}
		if replyToUrlFlag != "" {
			parentUrl := &pb.CastAddBody_ParentUrl{ParentUrl: replyToUrlFlag}
			messageBody.Parent = parentUrl
		}
		messageData := &pb.MessageData{
			Type:      pb.MessageType(pb.MessageType_value["MESSAGE_TYPE_CAST_ADD"]),
			Fid:       fid,
			Timestamp: uint32(time.Now().Unix() - fctools.FARCASTER_EPOCH),
			Network:   pb.FarcasterNetwork(pb.FarcasterNetwork_value["FARCASTER_NETWORK_MAINNET"]),
			Body: &pb.MessageData_CastAddBody{
				CastAddBody: messageBody,
			},
		}
		message := fctools.CreateMessage(messageData, privateKey, publicKey)
		if prepareFlag {
			jsonData, err := fctools.Marshal(
				message, fctools.MarshalOptions{Bytes2Hash: true, Timestamp2Date: false},
			)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(jsonData))
		} else {
			msg, err := hub.SubmitMessage(message)
			if err != nil {
				log.Fatalf("Error submitting message: %v", err)
			}
			fmt.Printf("Sent: @%d/0x%s\n", msg.Data.Fid, hex.EncodeToString(msg.Hash))
		}
		replyToCastFlag = "@" + strconv.FormatInt(int64(fid), 10) + "/" + "0x" + hex.EncodeToString(message.Hash)
		replyToUrlFlag = "" // When this creates a cast storm, only the first cast should use FIP-2.
	}
}

func init() {
	sendCmd.AddCommand(sendCastCmd)
	sendCastCmd.Flags().Uint64P("fid", "", 0, "Fid who is casting")
	sendCastCmd.Flags().StringP("privkey", "", "", "Application private key. Ex: 0xabc1234....")
	sendCastCmd.Flags().StringP("reply-to-cast", "", "", "Reply to a cast. The expected format is @fid/0xhash")
	sendCastCmd.Flags().StringP("reply-to-url", "", "", "Reply to a URL (aka FIP2)")
	sendCastCmd.Flags().BoolP("prepare", "", false, "Prepare the Message object and print it, but don't send it")
}
