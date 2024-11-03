package embedurl

import (
	"encoding/hex"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	pb "github.com/vrypan/fargo/farcaster"
)

const FARCASTER_EPOCH int64 = 1609459200

type Url struct {
	Link        string
	Fid         uint64
	Hash        []byte
	Timestamp   uint32
	Pos         uint8
	ContentType string
	Extension   string
}

func FromMessage(cast pb.Message) []Url {
	var ret []Url
	cast_body := pb.CastAddBody(*cast.Data.GetCastAddBody())
	for i, embed := range cast_body.Embeds {
		if link := embed.GetUrl(); link != "" {
			newUrl := Url{Link: link, Fid: cast.Data.Fid, Hash: cast.Hash, Timestamp: cast.Data.Timestamp, Pos: uint8(i)}
			newUrl.UpdateExtension()
			ret = append(ret, newUrl)
		}
	}
	return ret
}

func (u *Url) UnixTimestamp() int64 {
	return int64(u.Timestamp) + FARCASTER_EPOCH
}
func (u *Url) Id() string {
	return strconv.FormatUint(u.Fid, 10) +
		"-0x" + hex.EncodeToString(u.Hash) +
		"-" + strconv.Itoa(int(u.Pos))
}
func (u *Url) String() string {
	return u.Id() +
		"-" + u.ContentType +
		"-" + u.Extension

}

func (u *Url) UpdateContentType() string {
	if u.ContentType != "" {
		return u.ContentType
	}
	if resp, err := http.Head(u.Link); err == nil {
		defer resp.Body.Close()
		u.ContentType = resp.Header.Get("Content-Type")
		return u.ContentType
	}
	return "error"
}

func (u *Url) UpdateExtension() string {
	parsed, err := url.Parse(u.Link)
	if err != nil {
		return ""
	}
	last := path.Base(parsed.Path)
	if last == "/" || last == "." {
		return ""
	}
	ext := filepath.Ext(last)
	if len(ext) > 1 {
		u.Extension = ext[1:]
		return ext
	}
	return ""
}

func (u *Url) Filename() string {
	if u.Extension != "" {
		return u.Id() + "." + u.Extension
	}
	if u.ContentType != "" {
		p := strings.Split(u.ContentType, ";")
		p = strings.Split(p[0], "/")
		if len(p) > 1 {
			u.Extension = p[1]
			return u.Id() + "." + u.Extension
		}
	}
	return u.Id()
}
