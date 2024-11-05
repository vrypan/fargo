package localdb

import (
	//	"fmt"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/user"
	"path/filepath"
	//"github.com/vrypan/fargo/config"
)

// var lock sync.Mutex

var db_path = ""

const dot_dir = ".fargo"

type db_value struct {
	Idx uint64
	Val string
}

type _kv_store struct {
	Max_h  uint64
	Top    uint64
	Bottom uint64
	Kv     map[string]db_value
}

var kv_store _kv_store

var ERR_NOT_FOUND = errors.New("Not Found")
var ERR_NOT_STORED = errors.New("Not Stored")

func createDotDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	dotDir := filepath.Join(usr.HomeDir, dot_dir)
	if _, err := os.Stat(dotDir); os.IsNotExist(err) {
		err = os.Mkdir(dotDir, 0755)
		if err != nil {
			return "", err
		}
	}
	return dotDir, nil
}

func Set(k string, v string) error {
	if db_v, key_exists := kv_store.Kv[k]; key_exists {
		kv_store.Kv[k] = db_value{Idx: db_v.Idx, Val: v}
	} else {
		kv_store.Top++
		kv_store.Kv[k] = db_value{Idx: kv_store.Top, Val: v}
	}
	return nil
}

func Get(k string) (string, error) {
	if db_v, key_exists := kv_store.Kv[k]; key_exists {
		return db_v.Val, nil
	} else {
		return "", ERR_NOT_FOUND
	}
}

func save() error {
	f, err := os.Create(db_path)
	if err != nil {
		return err
	}
	var b []byte
	b, err = json.Marshal(kv_store)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, bytes.NewReader(b))
	return err
}

func load() error {
	if db_path == "" {
		if dot_dir, e := createDotDir(); e != nil {
			panic(e)
		} else {
			db_path = filepath.Join(dot_dir, "local.db")
		}
	}
	b, err := os.ReadFile(db_path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			panic(err)
		}
	}
	if err = json.Unmarshal(b, &kv_store); err != nil {
		return err
	}
	return nil
}

func Stats() (uint64, uint64, uint64) {
	return kv_store.Max_h, kv_store.Top, kv_store.Bottom
}

func Open() error {
	if kv_store.Kv == nil {
		//lock.Lock()
		kv_store.Kv = make(map[string]db_value)
		err := load()
		if err != nil {
			panic(err)
		}
	}
	return nil
}
func Close() error {
	//defer lock.Unlock()
	err := save()
	if err != nil {
		panic(err)
	}
	return nil
}
