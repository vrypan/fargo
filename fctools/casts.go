package fctools

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"time"

	pb "github.com/vrypan/fargo/farcaster"
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
func (c Cast) Json(hexHashes bool, realTimestamps bool) ([]byte, error) {
	//data := interface
	var jsonData interface{}
	jsonBytes, err := protojson.Marshal(c.Message)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonBytes, &jsonData)
	jsonPretty(jsonData, hexHashes, realTimestamps)
	if err != nil {
		return nil, err
	}
	updatedJsonBytes, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return nil, err
	}
	return updatedJsonBytes, nil
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
	grp.CollectFnames(hub)
	return grp
}

/*
Populates a CastGroup with recent likes from an Fid.
Head is set to nil.
*/
func (grp *CastGroup) FromCastIds(hub *FarcasterHub, castIds []*pb.CastId) *CastGroup {
	if hub == nil {
		hub = NewFarcasterHub()
		defer hub.Close()
	}
	for _, cid := range castIds {
		grp.AppendCast(hub, cid)
	}
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
	grp.CollectFnames(hub)
	return grp
}
func (grp *CastGroup) AppendCast(hub *FarcasterHub, castId *pb.CastId) *CastGroup {
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

func (grp *CastGroup) CollectFnames(hub *FarcasterHub) *CastGroup {
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
	return grp
}

func (grp *CastGroup) JsonList(hexHashes bool, realTimestamps bool) ([]byte, error) {
	groupData := make([]interface{}, len(grp.Messages))
	var jsonData interface{}
	idx := 0
	for _, message := range grp.Messages {
		json_bytes, err := protojson.Marshal(message.Message)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(json_bytes, &jsonData)
		jsonPretty(jsonData, hexHashes, realTimestamps)
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
func (grp *CastGroup) JsonThread(hexHashes bool, realTimestamps bool) ([]byte, error) {
	groupData := struct {
		Head    string                 `json:"head"`
		Casts   map[string]interface{} `json:"casts"`
		Replies map[string][]string    `json:"replies"`
		Fnames  map[uint64]string      `json:"fnames"`
	}{
		Head:    grp.Head.String(),
		Casts:   make(map[string]interface{}),
		Replies: make(map[string][]string),
		Fnames:  grp.Fnames,
	}

	for hash, message := range grp.Messages {
		var jsonData interface{}
		jsonBytes, err := protojson.Marshal(message.Message)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(jsonBytes, &jsonData)
		jsonPretty(jsonData, hexHashes, realTimestamps)
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

func jsonPretty(data interface{}, hexHashes bool, realTimestamps bool) {
	switch value := data.(type) {
	case map[string]interface{}:
		for k, v := range value {
			if hexHashes && (k == "hash" || k == "signature" || k == "signer") {
				bytes, err := base64.StdEncoding.DecodeString(v.(string))
				if err != nil {
					panic("Field value is not base64")
				}
				value[k] = "0x" + hex.EncodeToString(bytes)
			} else if realTimestamps && (k == "timestamp") {
				value[k] = time.Unix(int64(v.(float64))+FARCASTER_EPOCH, 0)
			} else {
				jsonPretty(v, hexHashes, realTimestamps)
			}
		}
	case []interface{}:
		for _, v := range value {
			jsonPretty(v, hexHashes, realTimestamps)
		}
	}
}
