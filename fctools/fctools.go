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

func (hub *FarcasterHub) SubmitMessageData(messageData *pb.MessageData, signerPrivate []byte, signerPublic []byte) (*pb.Message, error) {
	hash_scheme := pb.HashScheme(pb.HashScheme_value["HASH_SCHEME_BLAKE3"])
	signature_scheme := pb.SignatureScheme(pb.SignatureScheme_value["SIGNATURE_SCHEME_ED25519"])
	data_bytes, err := proto.Marshal(messageData)
	if err != nil {
		return nil, err
	}
	signerPublic_ := append(signerPrivate, signerPublic...) // required by ed25519 Go implementation

	hasher := blake3.New()
	hasher.Write(data_bytes)
	hash := hasher.Sum(nil)[0:20]

	signature := ed25519.Sign(signerPublic_, hash)

	message := pb.Message{
		Data:            messageData,
		Hash:            hash,
		HashScheme:      hash_scheme,
		Signature:       signature,
		SignatureScheme: signature_scheme,
		Signer:          signerPublic,
		DataBytes:       data_bytes,
	}
	msg, err := hub.client.SubmitMessage(hub.ctx, &message)
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
	} else {
		ud := pb.UserDataBody(*msg.Data.GetUserDataBody())
		return ud.Value, nil
	}
}

func (hub FarcasterHub) GetUsernameProofsByFid(fid uint64) ([]string, error) {
	msg, err := hub.client.GetUserNameProofsByFid(hub.ctx, &pb.FidRequest{Fid: fid})
	if err != nil {
		return nil, err
	}
	var ret []string
	for _, p := range msg.Proofs {
		ret = append(ret, string(p.Name))
	}
	return ret, nil
}

func (hub FarcasterHub) GetFidByUsername(username string) (uint64, error) {
	msg, err := hub.client.GetUsernameProof(hub.ctx, &pb.UsernameProofRequest{Name: []byte(username)})
	if err != nil {
		return 0, err
	}
	fid := msg.Fid
	return fid, nil
}

func (hub FarcasterHub) GetCastsByFid(fid uint64, page_size uint32) ([]*pb.Message, error) {
	var reverse bool = true
	//var page_size uint32 = 10
	msg, err := hub.client.GetCastsByFid(hub.ctx, &pb.FidRequest{Fid: fid, Reverse: &reverse, PageSize: &page_size})
	if err != nil {
		return nil, err
	}
	return msg.Messages, nil
}

func (hub FarcasterHub) GetCast(fid uint64, hash []byte) (*pb.Message, error) {
	msg, err := hub.client.GetCast(hub.ctx, &pb.CastId{Fid: fid, Hash: hash})
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (hub FarcasterHub) GetCastReplies(fid uint64, hash []byte) (*pb.MessagesResponse, error) {
	// GetCastsByParent(CastsByParentRequest) returns (MessagesResponse);
	response, err := hub.client.GetCastsByParent(
		hub.ctx,
		&pb.CastsByParentRequest{
			Parent: &pb.CastsByParentRequest_ParentCastId{
				ParentCastId: &pb.CastId{Fid: fid, Hash: hash},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	return response, nil
}
