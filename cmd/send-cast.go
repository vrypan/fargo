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
	Short: "Send a Cast Message to the hub",
	Run:   runSendCast,
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
		log.Fatalf("Private key error: %v", err)
	}

	s = config.GetString("cast.pubkey")
	if c, _ := cmd.Flags().GetString("pubkey"); c != "" {
		s = c
	}
	if len(s) < 2 {
		log.Fatal("Public key error: public key string too short")
	}
	if publicKey, err = hex.DecodeString(s[2:]); err != nil {
		log.Fatalf("Public key error: %v", err)
	}

	fid := uint64(config.GetInt("cast.fid"))
	if c, _ := cmd.Flags().GetUint64("fid"); c > 0 {
		fid = c
	}
	if fid == 0 {
		log.Fatal("No fid: fid is zero")
	}
	fmt.Printf("FID: %#v", fid)

	if len(args) == 0 {
		log.Fatal("Missing arguments: text argument required")
	}

	hub := fctools.NewFarcasterHub()
	defer hub.Close()

	text := args[0]

	castText, mentionsPositions, mentions, embeds := fctools.ProcessCastBody(text)
	fmt.Println(castText)
	fmt.Println(mentionsPositions)
	fmt.Println(mentions)
	fmt.Println(embeds)

	messageData := &pb.MessageData{
		Type:      pb.MessageType(pb.MessageType_value["MESSAGE_TYPE_CAST_ADD"]),
		Fid:       fid,
		Timestamp: uint32(time.Now().Unix() - fctools.FARCASTER_EPOCH),
		Network:   pb.FarcasterNetwork(pb.FarcasterNetwork_value["FARCASTER_NETWORK_MAINNET"]),
		Body: &pb.MessageData_CastAddBody{
			CastAddBody: &pb.CastAddBody{
				Mentions:          mentions,
				MentionsPositions: mentionsPositions,
				Parent:            nil,
				Text:              castText,
				Type:              0,
				Embeds:            embeds,
			}},
	}

	fmt.Printf("MessageData:\n%#v\n", messageData)
	fmt.Println("Keys", privateKey, publicKey)

	msg, err := hub.SubmitMessageData(
		messageData,
		privateKey,
		publicKey,
	)

	fmt.Println(msg)
	if err != nil {
		log.Fatalf("Error submitting message: %v", err)
	}
}

func init() {
	rootCmd.AddCommand(sendCastCmd)
	sendCastCmd.Flags().Uint64P("fid", "", 0, "Fid who is casting")
	sendCastCmd.Flags().StringP("pubkey", "", "", "Application public key")
	sendCastCmd.Flags().StringP("privkey", "", "", "Application private key")
}
