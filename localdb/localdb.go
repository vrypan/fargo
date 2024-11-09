package localdb

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/dgraph-io/badger/v3"
)

var db *badger.DB
var db_path = ""

const dot_dir = ".fargo"

var ERR_NOT_FOUND = errors.New("Not Found")

func IsOpen() bool {
	return db != nil
}

func AssertOpen() {
	if db == nil {
		panic("DB not open")
	}
}

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

func Set(k string, v []byte) error {
	/*
	err := db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(k), v).WithTTL(24 * time.Hour)
	})
	return err
	*/
	/*
	err := db.Update(func(txn *badger.Txn) error {
	  e := badger.NewEntry([]byte([]byte(k)), v).WithTTL(time.Hour)
	  err := txn.SetEntry(e)
	  return err
	})
	*/
	err := db.Update(func(txn *badger.Txn) error {
        e := badger.NewEntry([]byte(k), v).WithTTL(24 * time.Hour)
        return txn.SetEntry(e)
    })
	return err
}

func Get(k string) ([]byte, error) {
	var val []byte
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(k))
		if err != nil {
			return ERR_NOT_FOUND
		}
		err = item.Value(func(v []byte) error {
			val = v
			return nil
		})
		return err
	})
	if err != nil {
		return nil, err
	}
	return val, nil
}

func Open() error {
	if db_path == "" {
		if dotDir, err := createDotDir(); err != nil {
			return err
		} else {
			db_path = filepath.Join(dotDir, "local2.db")
		}
	}

	var err error
	db, err = badger.Open(badger.DefaultOptions(db_path).WithLoggingLevel(badger.ERROR))
	return err
}

func Close() error {
	return db.Close()
}
