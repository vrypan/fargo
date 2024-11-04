package cmd

import (
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/vrypan/fargo/config"
	pb "github.com/vrypan/fargo/farcaster"
	"github.com/vrypan/fargo/fctools"
)

var sendCastCmd = &cobra.Command{
	Use:   "cast [text]",
	Short: "Send new cast",
	Long: `[text] is the full cast text.
Any @mentions will be identified automatically and
the first two links will be converted to embeds.`,
	Run: runSendCast,
}

func runSendCast(cmd *cobra.Command, args []string) {
	config.Load()

	var err error
	var privateKey []byte
	var publicKey []byte
	var s string

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

	s = config.GetString("cast.pubkey")
	if c, _ := cmd.Flags().GetString("pubkey"); c != "" {
		s = c
	}
	if len(s) < 2 {
		log.Fatal("Public key error: public key string too short")
	}
	if publicKey, err = hex.DecodeString(s[2:]); err != nil {
		log.Fatalf("Public key error: %v\nUse --help to see options.", err)
	}

	fid := uint64(config.GetInt("cast.fid"))
	if c, _ := cmd.Flags().GetUint64("fid"); c > 0 {
		fid = c
	}
	if fid == 0 {
		log.Fatal("No fid: fid is zero. Use --help to see options.")
	}

	replyToFlag, _ := cmd.Flags().GetString("reply-to")

	prepareFlag, _ := cmd.Flags().GetBool("prepare")

	if len(args) == 0 {
		log.Fatal("Missing arguments: text argument required")
	}

	hub := fctools.NewFarcasterHub()
	defer hub.Close()

	text := args[0]
	castText, mentionsPositions, mentions, embeds := fctools.ProcessCastBody(text)
	if len(embeds) > 2 {
		embeds = embeds[0:2]
	}

	messageBody := &pb.CastAddBody{
		Mentions:          mentions,
		MentionsPositions: mentionsPositions,
		//Parent:            nil,
		Text:   castText,
		Type:   0,
		Embeds: embeds,
	}

	if replyToFlag != "" {
		parentFid, parentHashString := ParseFcURI(replyToFlag)
		parentHash := HashToBytes(parentHashString[0])
		parentCast := &pb.CastAddBody_ParentCastId{ParentCastId: &pb.CastId{Fid: parentFid, Hash: parentHash}}
		messageBody.Parent = parentCast
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

	if prepareFlag {
		jsonData, err := fctools.Marshal(
			messageData, fctools.MarshalOptions{Bytes2Hash: true, Timestamp2Date: false},
		)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(jsonData))
		return
	}
	msg, err := hub.SubmitMessageData(
		messageData,
		privateKey,
		publicKey,
	)
	if err != nil {
		log.Fatalf("Error submitting message: %v", err)
	}
	fmt.Printf("Sent: %s\n", fctools.FormatCastId(msg.Data.Fid, msg.Hash, ""))
}

func init() {
	sendCmd.AddCommand(sendCastCmd)
	sendCastCmd.Flags().Uint64P("fid", "", 0, "Fid who is casting")
	sendCastCmd.Flags().StringP("pubkey", "", "", "Application public key. Ex: 0xdef1234....")
	sendCastCmd.Flags().StringP("privkey", "", "", "Application private key. Ex: 0xabc1234....")
	sendCastCmd.Flags().StringP("reply-to", "", "", "Reply to a cast. The expected format is @fid/0xhash")
	sendCastCmd.Flags().BoolP("prepare", "", false, "Prepare the Message object and print it, but don't send it")
}
