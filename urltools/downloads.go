package urltools

import (
    "log"
    "os"
    "path/filepath"
    "github.com/hashicorp/go-getter"
    "os/user"
    "net/http"
)

func normalizeLocalPath(p string) string {
    if p[0:1] == "~" {
        usr, err := user.Current()
        if err != nil {
            log.Fatalf("%v\n",err)
        }
        home := usr.HomeDir
        return filepath.Join(home, p[1:])
    }
    return p
}

func GetFile(url string, dst string) string {
    dst = normalizeLocalPath(dst)
    if err := os.MkdirAll(dst, os.ModePerm); err != nil {
        log.Fatalf("Error creating directory: %v\n", err)
    }

    local_filename := filepath.Base(url)
    if local_filename == "" || local_filename == "/" {
        log.Fatalf("invalid URL path: %s", url)
    }

    if err := getter.GetAny(dst,url); err != nil {
        log.Fatalf("Error downloading file: %v", err)
    }

    file_path := filepath.Join(dst, local_filename)
    return file_path

}

func GetMimeType(url string) (string, error) {
    resp, err := http.Head(url)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    contentType := resp.Header.Get("Content-Type")
    return contentType, nil
}