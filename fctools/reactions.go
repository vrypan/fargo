package fctools

import (
	"encoding/hex"
	"encoding/json"
	"strconv"

	pb "github.com/vrypan/fargo/farcaster"
	"google.golang.org/protobuf/encoding/protojson"
)

type Reaction struct {
	Message *pb.Message
}

type Reactions struct {
	Messages []*Reaction
	Fnames   map[uint64]string
}

func (reaction *Reaction) String() string {
	var target string
	user_fid := strconv.FormatUint(reaction.Message.Data.Fid, 10)
	reactionType := reaction.Message.Data.GetReactionBody().Type.String()
	if reaction.Message.Data.GetReactionBody().GetTargetCastId() != nil {
		castId := reaction.Message.Data.GetReactionBody().GetTargetCastId()
		cast_fid := strconv.FormatUint(castId.Fid, 10)
		cast_hash := "0x" + hex.EncodeToString(castId.Hash)
		target = cast_fid + "/" + cast_hash
	}
	if reaction.Message.Data.GetReactionBody().GetTargetUrl() != "" {
		target = reaction.Message.Data.GetReactionBody().GetTargetUrl()
	}
	return user_fid + " " + reactionType + " --> " + target
}

func (reaction *Reaction) Json(hexHashes bool, realTimestamps bool) ([]byte, error) {
	//data := interface
	var jsonData interface{}
	jsonBytes, err := protojson.Marshal(reaction.Message)
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

func NewReactions() *Reactions {
	return &Reactions{
		Messages: make([]*Reaction, 0),
		Fnames:   make(map[uint64]string),
	}
}

/*
Populates a CastGroup with recent likes from an Fid.
Head is set to nil.
*/
func (reactions *Reactions) FromFid(hub *FarcasterHub, fid uint64, reactionType string, count uint32) *Reactions {
	if hub == nil {
		hub = NewFarcasterHub()
		defer hub.Close()
	}

	if messages, err := hub.GetReactionsByFid(fid, reactionType, count); err == nil {
		for _, reaction := range messages {
			reactions.Messages = append(reactions.Messages, &Reaction{Message: reaction})
		}
	}

	return reactions
}
func (reactions *Reactions) CollectFnames(hub *FarcasterHub) *Reactions {
	for _, msg := range reactions.Messages {
		reactions.Fnames[msg.Message.Data.Fid], _ = hub.PrxGetUserDataStr(msg.Message.Data.Fid, "USER_DATA_TYPE_USERNAME")
	}
	return reactions
}
func (reactions *Reactions) CastIds() []*pb.CastId {
	ids := make([]*pb.CastId, 0, len(reactions.Messages))
	for _, m := range reactions.Messages {
		if castId := m.Message.Data.GetReactionBody().GetTargetCastId(); castId != nil {
			ids = append(ids, castId)
		}
	}
	return ids
}
func (reactions *Reactions) JsonList(hexHashes bool, realTimestamps bool) ([]byte, error) {
	groupData := make([]interface{}, len(reactions.Messages))
	var jsonData interface{}
	idx := 0
	for _, message := range reactions.Messages {
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
