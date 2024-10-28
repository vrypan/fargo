package localdb

import (
//	"fmt"
	"errors"
	"sync"
	"os"
	"io"
	"bytes"
	"encoding/json"
)

var lock sync.Mutex

var db_path = "local.db"

type db_value struct {
	Idx uint64
	Val string
}

type _kv_store struct {
	Max_h 	uint64
	Top		uint64
	Bottom	uint64
	Kv 		map[string]db_value
}

var kv_store _kv_store

var old_kv_store map[string]string

var ERR_NOT_FOUND = errors.New("Not Found")
var ERR_NOT_STORED = errors.New("Not Stored")

func Set(k string, v string) error {
	if db_v, key_exists := kv_store.Kv[k] ; key_exists {
		kv_store.Kv[k] = db_value{Idx: db_v.Idx, Val: v}
	} else {
		kv_store.Top++
		kv_store.Kv[k] = db_value{Idx: kv_store.Top, Val: v}
	}
	return nil
}

func Get(k string) (string, error) {
	if db_v, key_exists := kv_store.Kv[k] ; key_exists {
		return db_v.Val, nil
	} else {
		return "", ERR_NOT_FOUND
	}
}

func save() error {
  f, err := os.Create(db_path); if err != nil {
    return err
  }
  defer f.Close()
  
  var b []byte
  b, err = json.Marshal(kv_store); if err != nil {
    return err
  }
  _, err = io.Copy(f, bytes.NewReader(b))
  return err
}

func load() error {
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

func Open() error {
	lock.Lock()
	kv_store.Kv = make(map[string]db_value)
	err := load()
	if err != nil {
		panic(err)
	}
	return nil
}
func Close() error {
	defer lock.Unlock()
	err := save()
	if err != nil {
		panic(err)
	}
	return nil
}