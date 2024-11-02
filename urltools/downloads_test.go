package urltools

import (
	"testing"
	"github.com/vrypan/fargo/config"
)

func Test_Download(t *testing.T) {
	config.Load()
	localDir := config.GetString("downloads.dir")

    url := "https://github.com/vrypan/fargo/releases/download/0.1.7/fargo_0.1.7_checksums.txt"
    
    
    local := GetFile(url, localDir)
    t.Logf("Remote: %v\n", url)
    t.Logf("Local: %v\n", local) 
}

func Test_Mimetypes(t *testing.T) {
	var url string
	var mimetype string
	var err error

	url = "https://vrypan.net"
	mimetype, err = GetMimeType(url)
	t.Logf("%s \t %s\t%v", url, mimetype, err)

	url = "https://traffic.libsyn.com/atpfm/atp611.mp3"
	mimetype, err = GetMimeType(url)
	t.Logf("%s \t %s\t%v", url, mimetype, err)

	url = "https://github.com/vrypan/fargo/releases/download/0.1.7/fargo_0.1.7_checksums.txt"
	mimetype, err = GetMimeType(url)
	t.Logf("%s \t %s\t%v", url, mimetype, err)
}