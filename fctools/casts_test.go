package fctools

import (
	"encoding/hex"
	"testing"

	pb "github.com/vrypan/fargo/farcaster"
)

func Test_JsonThread(t *testing.T) {
	hash, _ := hex.DecodeString("13d491c6583c9bac177e6f8d76791e7326def624")
	castId := pb.CastId{Fid: 280, Hash: hash}
	s, err := NewCastGroup().FromCast(nil, &castId, true).JsonThread()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("\n" + string(s))
}

func Test_ThreadFromCast(t *testing.T) {
	grp := NewCastGroup()
	hash, _ := hex.DecodeString("da8b15abd2bde4545ecaf548df0192663746a426")
	castId := pb.CastId{Fid: 20396, Hash: hash}
	grp.FromCast(nil, &castId, true)
	// t.Log(grp)
	if hex.EncodeToString(grp.Head[:]) != "9d423ea8741c753eb17832c377e14027e6ea4fbd" {
		t.Fatalf("msg.Head is not 0x9d423ea8741c753eb17832c377e14027e6ea4fbd")
	}
}

func Test_ThreadFromCast2(t *testing.T) {
	grp := NewCastGroup()
	hash, _ := hex.DecodeString("13d491c6583c9bac177e6f8d76791e7326def624")
	castId := pb.CastId{Fid: 280, Hash: hash}
	grp.FromCast(nil, &castId, true)
	if len(grp.Messages) < 10 {
		t.Fatalf("Only %d casts in thread?", len(grp.Messages))
	}
	t.Logf("Total messages in thread: %d", len(grp.Messages))
	t.Log("\n" + grp.PprintThread(nil, 0))
}

func Test_ThreadFromCast_No_Expand(t *testing.T) {
	hash, _ := hex.DecodeString("13d491c6583c9bac177e6f8d76791e7326def624")
	castId := pb.CastId{Fid: 280, Hash: hash}
	grp := NewCastGroup().FromCast(nil, &castId, false)
	if len(grp.Messages) > 1 {
		t.Fatalf("There are %d messages in group. Expected: 1.", len(grp.Messages))
	}
	t.Logf("Total messages in thread: %d", len(grp.Messages))
	t.Log("\n" + grp.PprintThread(nil, 0))
}

func Test_ThreadFromCast_Chain(t *testing.T) {
	hash, _ := hex.DecodeString("13d491c6583c9bac177e6f8d76791e7326def624")
	castId := pb.CastId{Fid: 280, Hash: hash}
	s := NewCastGroup().FromCast(nil, &castId, false).PprintThread(nil, 0)
	t.Log("\n" + s)
}

func Test_ThreadFromFid(t *testing.T) {
	grp := NewCastGroup().FromFid(nil, 280, 20)

	s := grp.PprintList(nil, 0)
	t.Log("\n" + s)
}

func Test_Json(t *testing.T) {
	grp := NewCastGroup().FromFid(nil, 280, 20)
	if len(grp.Messages) != 20 {
		t.Fatalf("Expected 20 messages, got %v\n", len(grp.Messages))
	}
	s, err := grp.JsonList()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("\n" + string(s))
}

func Test_Links(t *testing.T) {
	grp := NewCastGroup().FromFid(nil, 280, 100)

	links := grp.Links()
	t.Log(links)
}
