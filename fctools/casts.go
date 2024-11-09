package fctools

import (
	"encoding/hex"
	"encoding/json"
	"strconv"

	pb "github.com/vrypan/fargo/farcaster"
	db "github.com/vrypan/fargo/localdb"
	"google.golang.org/protobuf/encoding/protojson"
)

type Cast struct {
	Message *pb.Message
	Replies []Hash
}

func (c Cast) String() string {
	fid := strconv.FormatUint(c.Message.Data.Fid, 10)
	hash := "0x" + hex.EncodeToString(c.Message.Hash)
	return fid + "/" + hash
}
func (c Cast) Hash() string {
	return "0x" + hex.EncodeToString(c.Message.Hash)
}
func (c Cast) Fid() string {
	return strconv.FormatUint(c.Message.Data.Fid, 10)
}
func (c Cast) Text() string {
	return c.Message.Data.GetCastAddBody().Text
}

type CastGroup struct {
	Head     Hash
	Messages map[Hash]*Cast
	Fnames   map[uint64]string
}

func NewCastGroup() *CastGroup {
	return &CastGroup{Messages: make(map[Hash]*Cast), Fnames: make(map[uint64]string)}
}

/*
Populates a CastGroup with recent casts from an Fid.
Head is set to nil.
*/
func (grp *CastGroup) FromFid(hub *FarcasterHub, fid uint64, count uint32) *CastGroup {
	db.Open()
	defer db.Close()

	if hub == nil {
		hub = NewFarcasterHub()
		defer hub.Close()
	}
	messages, err := hub.GetCastsByFid(fid, count)
	if err != nil {
		return grp
	}

	var hash Hash
	//var cast *pb.Message
	for _, cast := range messages {
		hash = Hash(cast.Hash[:])
		grp.Messages[hash] = &Cast{Message: cast}
	}
	grp.collectFnames(hub)
	return grp
}

/*
Populates a CastGroup with all casts in a thread that a CastId participates.
Head is set to the top cast of the thread.
*/
func (grp *CastGroup) FromCastFidHash(hub *FarcasterHub, fid uint64, hash string, expandTree bool) *CastGroup {
	hash_b, err := hex.DecodeString(hash)
	if err != nil {
		return nil
	}
	castId := &pb.CastId{
		Fid:  fid,
		Hash: hash_b,
	}
	return grp.FromCast(hub, castId, expandTree)
}

func (grp *CastGroup) FromCast(hub *FarcasterHub, castId *pb.CastId, expandTree bool) *CastGroup {
	db.Open()
	defer db.Close()

	if hub == nil {
		hub = NewFarcasterHub()
		defer hub.Close()
	}
	cast, err := hub.PrxGetCast(castId.Fid, castId.Hash)
	if err != nil {
		return grp
	}
	grp.Messages[Hash(cast.Hash)] = &Cast{Message: cast}
	grp.Head = Hash(cast.Hash)
	if expandTree {
		for cast != nil {
			grp.Messages[Hash(cast.Hash)] = &Cast{Message: cast}
			parentCastId := cast.Data.GetCastAddBody().GetParentCastId()
			if parentCastId == nil {
				grp.Head = Hash(cast.Hash)
				break
			}
			cast, err = hub.PrxGetCast(parentCastId.Fid, parentCastId.Hash)
			if err != nil {
				break
			}
		}
		grp.expandReplies(hub, grp.Head)
	}
	grp.collectFnames(hub)
	return grp
}

func (grp *CastGroup) expandReplies(hub *FarcasterHub, hash Hash) {
	replies, err := hub.GetCastReplies(grp.Messages[hash].Message.Data.Fid, grp.Messages[hash].Message.Hash)
	if err != nil {
		return
	}
	for _, r := range replies.Messages {
		parent := grp.Messages[hash]
		parent.Replies = append(parent.Replies, Hash(r.Hash))
		grp.Messages[Hash(r.Hash)] = &Cast{Message: r}
		grp.expandReplies(hub, Hash(r.Hash))
	}
}

func (grp *CastGroup) collectFnames(hub *FarcasterHub) {
	for _, msg := range grp.Messages {
		grp.Fnames[msg.Message.Data.Fid], _ = hub.PrxGetUserDataStr(msg.Message.Data.Fid, "USER_DATA_TYPE_USERNAME")

		for _, mention := range msg.Message.GetData().GetCastAddBody().GetMentions() {
			grp.Fnames[mention], _ = hub.PrxGetUserDataStr(mention, "USER_DATA_TYPE_USERNAME")
		}

		for _, embed := range msg.Message.GetData().GetCastAddBody().GetEmbeds() {
			if cid := embed.GetCastId(); cid != nil {
				grp.Fnames[cid.Fid], _ = hub.PrxGetUserDataStr(cid.Fid, "USER_DATA_TYPE_USERNAME")
			}
		}
		if msg.Message.GetData().GetCastAddBody().GetParentCastId() != nil {
			p_cast_fid := msg.Message.GetData().GetCastAddBody().GetParentCastId().Fid
			p_cast_fname, _ := hub.PrxGetUserDataStr(p_cast_fid, "USER_DATA_TYPE_USERNAME")
			grp.Fnames[p_cast_fid] = p_cast_fname
		}

	}
}

func (grp *CastGroup) PprintThread(hash *Hash, padding int) string {
	if hash == nil {
		hash = &grp.Head
	}
	out := ""
	cast := grp.Messages[*hash].Message
	out += FormatCast(cast, grp.Fnames, padding, (padding == 0), "", "")
	for _, reply := range grp.Messages[*hash].Replies {
		out += grp.PprintThread(&reply, padding+4)
	}
	return out
}
func (grp *CastGroup) PprintList(hash *Hash, padding int) string {
	out := ""
	for _, cast := range grp.Messages {
		out += FormatCast(cast.Message, grp.Fnames, padding, true, "", "") + "\n"
	}
	return out
}

func (grp *CastGroup) JsonList() ([]byte, error) {
	groupData := make([]interface{}, len(grp.Messages))
	var jsonData interface{}
	idx := 0
	for _, message := range grp.Messages {
		json_bytes, err := protojson.Marshal(message.Message)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(json_bytes, &jsonData)
		if err != nil {
			return nil, err
		}
		groupData[idx] = jsonData
		idx++
	}
	updatedJsonBytes, err := json.MarshalIndent(groupData, "", "  ")
	if err != nil {
		return nil, err
	}
	return updatedJsonBytes, nil
}

/*
JsonThread returns a JSON string.
"casts" is a map hash->message
"head" contains the hash of the first message in the thread.
replies[casts[x]["hash"]] contains the hashes of the replies to casts[x]
You can follow the thread by checking
*/
func (grp *CastGroup) JsonThread() ([]byte, error) {
	groupData := struct {
		Head    string                 `json:"head"`
		Casts   map[string]interface{} `json:"casts"`
		Replies map[string][]string    `json:"replies"`
	}{
		Head:    grp.Head.String(),
		Casts:   make(map[string]interface{}),
		Replies: make(map[string][]string),
	}

	for hash, message := range grp.Messages {
		var jsonData interface{}
		jsonBytes, err := protojson.Marshal(message.Message)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(jsonBytes, &jsonData)
		if err != nil {
			return nil, err
		}
		groupData.Casts[hash.String()] = jsonData

		replyHashes := make([]string, len(message.Replies))
		for i, replyHash := range message.Replies {
			replyHashes[i] = replyHash.String()
		}
		groupData.Replies[hash.String()] = replyHashes
	}
	updatedJsonBytes, err := json.MarshalIndent(groupData, "", "  ")
	if err != nil {
		return nil, err
	}
	return updatedJsonBytes, nil
}

func (grp *CastGroup) Links() []string {
	links := []string{}
	for _, message := range grp.Messages {
		if embeds := message.Message.Data.GetCastAddBody().GetEmbeds(); len(embeds) > 0 {
			for _, e := range embeds {
				if e.GetUrl() != "" {
					links = append(links, e.GetUrl())
				}
			}
		}
	}
	return links
}
