package fctools

import (
	"context"
	"log"

	//"time"
	"crypto/ed25519"
	"encoding/json"

	"github.com/vrypan/fargo/config"
	pb "github.com/vrypan/fargo/farcaster"
	"github.com/zeebo/blake3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

const FARCASTER_EPOCH int64 = 1609459200
const FMT_COLS = 80

type FarcasterHub struct {
	hubAddr    string
	conn       *grpc.ClientConn
	client     pb.HubServiceClient
	ctx        context.Context
	ctx_cancel context.CancelFunc
}

func NewFarcasterHub() *FarcasterHub {
	config.Load()
	hubAddr := config.GetString("hub.host") + ":" + config.GetString("hub.port")
	cred := insecure.NewCredentials()

	if config.GetBool("hub.ssl") {
		cred = credentials.NewClientTLSFromCert(nil, "")
	}

	conn, err := grpc.DialContext(context.Background(), hubAddr, grpc.WithTransportCredentials(cred))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	client := pb.NewHubServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	return &FarcasterHub{
		hubAddr:    hubAddr,
		conn:       conn,
		client:     client,
		ctx:        ctx,
		ctx_cancel: cancel,
	}
}

func (h FarcasterHub) Close() {
	h.conn.Close()
	h.ctx_cancel()
}

func (hub FarcasterHub) HubInfo() ([]byte, error) {
	res, err := hub.client.GetInfo(hub.ctx, &pb.HubInfoRequest{DbStats: false})
	if err != nil {
		log.Fatalf("could not get HubInfo: %v", err)
		return nil, err
	}
	b, err := json.Marshal(res)
	return b, err
}

func (hub FarcasterHub) SubmitMessageData(messageData *pb.MessageData, signerPrivate, signerPublic []byte) (*pb.Message, error) {
	const hashLen = 20

	dataBytes, err := proto.Marshal(messageData)
	if err != nil {
		return nil, err
	}

	fullHash := blake3.Sum256(dataBytes)
	hash := fullHash[:hashLen]

	signature := ed25519.Sign(append(signerPrivate, signerPublic...), hash)

	message := pb.Message{
		Data:            messageData,
		Hash:            hash,
		HashScheme:      pb.HashScheme_HASH_SCHEME_BLAKE3,
		Signature:       signature,
		SignatureScheme: pb.SignatureScheme_SIGNATURE_SCHEME_ED25519,
		Signer:          signerPublic,
		DataBytes:       dataBytes,
	}

	return hub.client.SubmitMessage(hub.ctx, &message)
}

func (hub FarcasterHub) SubmitMessage(message *pb.Message) (*pb.Message, error) {
	msg, err := hub.client.SubmitMessage(hub.ctx, message)
	return msg, err
}

func (hub FarcasterHub) GetUserData(fid uint64, user_data_type string, tojson bool) (string, error) {
	_udt := pb.UserDataType(pb.UserDataType_value[user_data_type])
	msg, err := hub.client.GetUserData(hub.ctx, &pb.UserDataRequest{Fid: fid, UserDataType: _udt})
	if err != nil {
		return "", err
	}
	if tojson {
		b, err := json.Marshal(msg)
		return string(b), err
	}
	return pb.UserDataBody(*msg.Data.GetUserDataBody()).Value, nil
}

func (hub FarcasterHub) GetUsernameProofsByFid(fid uint64) ([]string, error) {
	msg, err := hub.client.GetUserNameProofsByFid(hub.ctx, &pb.FidRequest{Fid: fid})
	if err != nil {
		return nil, err
	}
	ret := make([]string, len(msg.Proofs))
	for i, p := range msg.Proofs {
		ret[i] = string(p.Name)
	}
	return ret, nil
}

func (hub FarcasterHub) GetFidByUsername(username string) (uint64, error) {
	msg, err := hub.client.GetUsernameProof(hub.ctx, &pb.UsernameProofRequest{Name: []byte(username)})
	if err != nil {
		return 0, err
	}
	return msg.Fid, nil
}

func (hub FarcasterHub) GetCastsByFid(fid uint64, pageSize uint32) ([]*pb.Message, error) {
	reverse := true
	msg, err := hub.client.GetCastsByFid(hub.ctx, &pb.FidRequest{Fid: fid, Reverse: &reverse, PageSize: &pageSize})
	if err != nil {
		return nil, err
	}
	return msg.Messages, nil
}

func (hub FarcasterHub) GetCast(fid uint64, hash []byte) (*pb.Message, error) {
	return hub.client.GetCast(hub.ctx, &pb.CastId{Fid: fid, Hash: hash})
}

func (hub FarcasterHub) GetCastReplies(fid uint64, hash []byte) (*pb.MessagesResponse, error) {
	return hub.client.GetCastsByParent(
		hub.ctx,
		&pb.CastsByParentRequest{
			Parent: &pb.CastsByParentRequest_ParentCastId{
				ParentCastId: &pb.CastId{Fid: fid, Hash: hash},
			},
		},
	)
}
