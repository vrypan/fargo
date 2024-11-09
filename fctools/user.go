package fctools

import (
	"encoding/json"
	"strconv"

	pb "github.com/vrypan/fargo/farcaster"
	"google.golang.org/protobuf/encoding/protojson"
)

type User struct {
	Fid      uint64
	UserData map[string]*pb.Message
}

func NewUser() *User {
	return &User{UserData: make(map[string]*pb.Message)}
}

func (u *User) FromFid(fid uint64) *User {
	u.Fid = fid
	return u
}

func (u *User) FromFname(hub *FarcasterHub, fname string) *User {
	fid, err := strconv.ParseUint(fname, 10, 64)
	if err != nil {
		if hub == nil {
			hub = NewFarcasterHub()
			defer hub.Close()
		}
		fid, err = hub.GetFidByUsername(fname)
		if err != nil {
			return nil
		}
	}
	u.Fid = fid
	return u
}

/*
hub == nil --> Create new hub connection
types == nil --> Fetch all USER_DATA_TYPE_*
*/
func (u *User) FetchUserData(hub *FarcasterHub, types []string) *User {
	// db.Open()
	// defer db.Close()
	if hub == nil {
		hub = NewFarcasterHub()
		defer hub.Close()
	}
	if types == nil {
		types = make([]string, len(pb.UserDataType_name))
		for _, tn := range pb.UserDataType_name {
			types = append(types, tn)
		}
	}
	for _, t := range types {
		message, err := hub.GetUserData(u.Fid, t)
		if err == nil {
			u.UserData[t] = message
		}
	}
	return u
}

func (u *User) Value(t string) string {
	return u.UserData[t].Data.GetUserDataBody().GetValue()
}

func (u *User) Json() ([]byte, error) {
	UserData := make(map[string]json.RawMessage)
	for _, message := range u.UserData {
		jsonBytes, err := protojson.Marshal(message)
		if err != nil {
			return nil, err
		}
		UserData[message.Data.GetUserDataBody().Type.String()] = jsonBytes
	}
	updatedJsonBytes, err := json.MarshalIndent(UserData, "", "  ")
	if err != nil {
		return nil, err
	}
	return updatedJsonBytes, nil
}

func (u *User) String() string {
	var out string
	for _, message := range u.UserData {
		out += message.Data.GetUserDataBody().Type.String() +
			": " +
			message.Data.GetUserDataBody().GetValue() +
			"\n"
	}
	return out
}
