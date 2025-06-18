package fctools

import (
	"context"
	"fmt"
	"log"
	"strconv"

	//"time"
	"crypto/ed25519"
	"encoding/hex"

	pb "github.com/vrypan/farcaster-go/farcaster"
	"github.com/vrypan/fargo/config"
	db "github.com/vrypan/fargo/localdb"
	"github.com/zeebo/blake3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

const FARCASTER_EPOCH int64 = 1609459200

type FarcasterHub struct {
	hubAddr    string
	conn       *grpc.ClientConn
	client     pb.HubServiceClient
	ctx        context.Context
	ctx_cancel context.CancelFunc
}

func apiKeyInterceptor(header, value string) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		md := metadata.Pairs(header, value)
		ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func NewFarcasterHub() *FarcasterHub {
	config.Load()
	hubAddr := config.GetString("hub.host") + ":" + config.GetString("hub.port")
	cred := insecure.NewCredentials()

	if config.GetBool("hub.ssl") {
		cred = credentials.NewClientTLSFromCert(nil, "")
	}
	var interceptor grpc.UnaryClientInterceptor

	if apikey := config.GetString("hub.apikey"); apikey != "" {
		interceptor = apiKeyInterceptor("x-api-key", apikey)
	}

	conn, err := grpc.Dial(
		hubAddr,
		grpc.WithTransportCredentials(cred),
		grpc.WithUnaryInterceptor(interceptor),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(20*1024*1024)),
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

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

/*
Not avaliable on Snapchain
func (hub FarcasterHub) HubInfo() ([]byte, error) {
	res, err := hub.client.GetInfo(hub.ctx, &pb.HubInfoRequest{DbStats: false})
	if err != nil {
		log.Fatalf("could not get HubInfo: %v", err)
		return nil, err
	}
	b, err := json.Marshal(res)
	return b, err
}
*/

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

func (hub FarcasterHub) GetUserData(fid uint64, user_data_type string) (*pb.Message, error) {
	udt := pb.UserDataType(pb.UserDataType_value[user_data_type])
	message, err := hub.client.GetUserData(hub.ctx, &pb.UserDataRequest{Fid: fid, UserDataType: udt})
	if err != nil {
		return nil, err
	}
	return message, nil
}
func (hub FarcasterHub) GetUserDataByFid(fid uint64) (*pb.MessagesResponse, error) {
	messages, err := hub.client.GetUserDataByFid(hub.ctx, &pb.FidRequest{Fid: fid})
	if err != nil {
		return nil, err
	}
	return messages, nil
}
func (hub FarcasterHub) GetUserDataStr(fid uint64, user_data_type string) (string, error) {
	message, err := hub.GetUserData(fid, user_data_type)
	if err != nil {
		return "", err
	}
	s := message.Data.GetUserDataBody().GetValue()
	return string(s), err
}
func (hub FarcasterHub) PrxGetUserDataStr(fid uint64, user_data_type string) (string, error) {
	db.AssertOpen()
	val, err := db.Get("GetUserData/" + strconv.FormatUint(fid, 10) + "/" + user_data_type)
	switch err {
	case nil:
		return string(val), nil
	case db.ERR_NOT_FOUND:
		s, err := hub.GetUserDataStr(fid, user_data_type)
		if err != nil {
			return "", err
		}
		err = db.Set(
			"GetUserData/"+strconv.FormatUint(fid, 10)+"/"+user_data_type,
			[]byte(s),
		)
		if err != nil {
			log.Printf("Could not cache GetUserData result: %v", err)
		}
		return s, nil
	default:
		log.Fatal(err)
		return "", err
	}
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
	message, err := hub.client.GetUsernameProof(hub.ctx, &pb.UsernameProofRequest{Name: []byte(username)})
	if err != nil {
		return 0, fmt.Errorf("failed to get username proof: %w", err)
	}
	return message.Fid, nil
}
func (hub FarcasterHub) PrxGetFidByUsername(username string) (uint64, error) {
	db.AssertOpen()
	fidBytes, err := db.Get("GetFidByUsername/" + username)
	switch err {
	case db.ERR_NOT_FOUND:
		fid, err := hub.GetFidByUsername(username)
		if err != nil {
			return 0, err
		}
		fidBytes = []byte(strconv.FormatUint(fid, 10))
		err = db.Set("GetFidByUsername/"+username, fidBytes)
		if err != nil {
			log.Printf("Could not cache GetFidByUsername result: %v", err)
		}
		return fid, nil
	case nil:
		fid, err := strconv.ParseUint(string(fidBytes), 10, 64)
		if err != nil {
			return 0, err
		}
		return fid, nil
	default:
		log.Fatal(err)
		return 0, err
	}
}

func (hub FarcasterHub) GetCastsByFid(fid uint64, pageSize uint32) ([]*pb.Message, error) {
	reverse := true
	msg, err := hub.client.GetCastsByFid(hub.ctx, &pb.FidRequest{Fid: fid, Reverse: &reverse, PageSize: &pageSize})

	if err != nil {
		return nil, err
	}
	return msg.Messages, nil
}

func (hub FarcasterHub) GetReactionsByFid(fid uint64, reaction string, pageSize uint32) ([]*pb.Message, error) {
	reverse := true
	reactionType := pb.ReactionType(pb.ReactionType_value[reaction])
	msg, err := hub.client.GetReactionsByFid(hub.ctx,
		&pb.ReactionsByFidRequest{Fid: fid, ReactionType: &reactionType, Reverse: &reverse, PageSize: &pageSize},
	)
	if err != nil {
		return nil, err
	}
	return msg.Messages, nil
}

func (hub FarcasterHub) GetCast(fid uint64, hash []byte) (*pb.Message, error) {
	return hub.client.GetCast(hub.ctx, &pb.CastId{Fid: fid, Hash: hash})
}

func (hub FarcasterHub) PrxGetCast(fid uint64, hash []byte) (*pb.Message, error) {
	db.AssertOpen()
	dbKey := "GetCast/" + hex.EncodeToString(hash)
	messageBytes, err := db.Get(dbKey)
	switch err {
	case db.ERR_NOT_FOUND:
		message, err := hub.GetCast(fid, hash)
		if err != nil {
			return nil, err
		}
		messageBytes, err := proto.Marshal(message)
		err = db.Set(dbKey, messageBytes)
		if err != nil {
			log.Printf("Could not cache GetCast result: %v", err)
		}
		return message, nil
	case nil:
		var message pb.Message
		err = proto.Unmarshal(messageBytes, &message)
		if err != nil {
			return nil, err
		}
		return &message, nil
	default:
		log.Fatal(err)
		return nil, err
	}
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
