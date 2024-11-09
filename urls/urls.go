package urls

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"
)

//const FARCASTER_EPOCH int64 = 1609459200

type Url struct {
	Link        string
	ContentType string
	Extension   string
}

func NewUrl(link string) *Url {
	return &Url{Link: link}
}

func (u *Url) String() string {
	return u.Filename()
}

func (u *Url) UpdateType() *Url {
	if u.ContentType != "" {
		return u
	}
	if resp, err := http.Head(u.Link); err == nil {
		defer resp.Body.Close()
		u.ContentType = resp.Header.Get("Content-Type")
		return u
	}
	return u
}

func (u *Url) Ext() string {
	if u.Extension != "" {
		return u.Extension
	}
	if u.ContentType != "" {
		p := strings.Split(u.ContentType, ";")
		p = strings.Split(p[0], "/")
		if len(p) > 1 {
			return p[1]
		}
	}
	return ""
}

func (u *Url) UpdateExt() *Url {
	parsed, err := url.Parse(u.Link)
	if err != nil {
		return u
	}
	last := path.Base(parsed.Path)
	if last == "/" || last == "." {
		return u
	}
	ext := filepath.Ext(last)
	if len(ext) > 1 {
		u.Extension = ext[1:]
		return u
	}
	return u
}

func (u *Url) Id() string {
	hash := md5.Sum([]byte(u.Link))
	return hex.EncodeToString(hash[:])
}

func (u *Url) Filename() string {
	basename := u.Id()
	parsedURL, err := url.Parse(u.Link)
	if err != nil {
		if filepath.Ext(u.Link) == "" {
			return basename + "." + u.Ext()
		}
		return basename
	}
	b := filepath.Base(parsedURL.Path)
	if b != "/" {
		basename += "-" + b
	}
	if filepath.Ext(u.Link) == "" {
		return basename + "." + u.Ext()
	}
	return basename
}
