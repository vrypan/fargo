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

func HubInfo() ([]byte, error) {
	conn, err := grpc.Dial(hubAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
            log.Fatalf("did not connect: %v", err)
            return nil, err
    }
    defer conn.Close()

    client := pb.NewHubServiceClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    hub, err := client.GetInfo(ctx, &pb.HubInfoRequest{DbStats: false})
    if err != nil {
            log.Fatalf("could not get HubInfo: %v", err)
            return nil, err
    }
    b, err := json.Marshal(hub)
    return b, err
}


func GetUserData( fid uint64, user_data_type string, tojson bool) (string, error) {
    conn, err := grpc.Dial(hubAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
            log.Fatalf("did not connect: %v", err)
            return "", err
    }
    defer conn.Close()

    client := pb.NewHubServiceClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    _udt := pb.UserDataType(pb.UserDataType_value[user_data_type])
    msg, err := client.GetUserData(ctx, &pb.UserDataRequest{Fid:fid, UserDataType: _udt})
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

func GetUsernameProofsByFid(fid uint64) ([]string, error) {
    conn, err := grpc.Dial(hubAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return nil, err
    }
    defer conn.Close()

    client := pb.NewHubServiceClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    msg, err := client.GetUserNameProofsByFid(ctx, &pb.FidRequest{Fid: fid})
    if err != nil {
        return nil, err
    }
    var ret []string
    for _, p := range msg.Proofs {
        ret = append(ret, string(p.Name))
    }
    return ret, nil
}

func GetFidByUsername(username string) (uint64, error){
    conn, err := grpc.Dial(hubAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
            log.Fatalf("did not connect: %v", err)
            return uint64(0), err
    }
    defer conn.Close()

    client := pb.NewHubServiceClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    msg, err := client.GetUsernameProof(ctx, &pb.UsernameProofRequest{Name: []byte(username)})
    if err != nil {
    	return 0, err
    }
    fid := msg.Fid
    return fid, nil    
}

func GetCastsByFid( fid uint64 ) ([]*pb.Message, error) {
    conn, err := grpc.Dial(hubAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
            log.Fatalf("did not connect: %v", err)
            return nil, err
    }
    defer conn.Close()

    client := pb.NewHubServiceClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    var reverse bool = true
    var page_size uint32 = 10
    msg, err := client.GetCastsByFid(ctx, &pb.FidRequest{Fid: fid, Reverse: &reverse, PageSize: &page_size})
    if err != nil {
    	return nil, err
    }
    return msg.Messages, nil
}

func GetCast( fid uint64, hash []byte ) ( *pb.Message, error) {
    conn, err := grpc.Dial(hubAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
            log.Fatalf("did not connect: %v", err)
            return nil, err
    }
    defer conn.Close()

    client := pb.NewHubServiceClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    msg, err := client.GetCast(ctx, &pb.CastId{Fid: fid, Hash: hash})
    if err != nil {
        return nil, err
    }
    return msg, nil
}






















