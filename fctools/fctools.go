package fctools

import (
        "context"
        "log"
        "time"
        pb "github.com/vrypan/fargo/farcaster"
        "google.golang.org/grpc"
        "google.golang.org/grpc/credentials/insecure"
        "encoding/json"
)

const hubAddr string = "38.242.252.228:2283"

type FarcasterHub struct {
    hubAddr string
    conn *grpc.ClientConn
    client pb.HubServiceClient
    ctx context.Context
    ctx_cancel context.CancelFunc
}

func NewFarcasterHub() *FarcasterHub {
    conn, err := grpc.Dial(hubAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
            log.Fatalf("Did not connect: %v", err)
    }
    client := pb.NewHubServiceClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    return &FarcasterHub {
        hubAddr: hubAddr,
        conn: conn,
        client: client,
        ctx: ctx,
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

func (hub FarcasterHub) GetUserData( fid uint64, user_data_type string, tojson bool) (string, error) {
    _udt := pb.UserDataType(pb.UserDataType_value[user_data_type])
    msg, err := hub.client.GetUserData(hub.ctx, &pb.UserDataRequest{Fid:fid, UserDataType: _udt})
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

func (hub FarcasterHub) GetFidByUsername(username string) (uint64, error){
    msg, err := hub.client.GetUsernameProof(hub.ctx, &pb.UsernameProofRequest{Name: []byte(username)})
    if err != nil {
    	return 0, err
    }
    fid := msg.Fid
    return fid, nil    
}

func (hub FarcasterHub) GetCastsByFid( fid uint64 ) ([]*pb.Message, error) {
    var reverse bool = true
    var page_size uint32 = 10
    msg, err := hub.client.GetCastsByFid(hub.ctx, &pb.FidRequest{Fid: fid, Reverse: &reverse, PageSize: &page_size})
    if err != nil {
    	return nil, err
    }
    return msg.Messages, nil
}

func (hub FarcasterHub) GetCast( fid uint64, hash []byte ) ( *pb.Message, error) {
    msg, err := hub.client.GetCast(hub.ctx, &pb.CastId{Fid: fid, Hash: hash})
    if err != nil {
        return nil, err
    }
    return msg, nil
}
