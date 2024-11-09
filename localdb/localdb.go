package localdb

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/vrypan/fargo/config"
)

var db *badger.DB
var db_path = ""
var ttl = 24

const dot_dir = ".fargo"

func init() {
	if db_path == "" {
		configDir, err := config.ConfigDir()
		if err != nil {
			panic(err)
		}
		db_path = filepath.Join(configDir, "local.db")
	}
}
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

var (
	ERR_NOT_FOUND  = errors.New("Not Found")
	ERR_NOT_STORED = errors.New("Not Stored")
)

func Set(k string, v []byte) error {
	err := db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(k), v).WithTTL(time.Duration(ttl) * time.Hour)
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
	config.Load()
	ttl = config.GetInt("db.ttlhours")

	var err error
	db, err = badger.Open(badger.DefaultOptions(db_path).WithLoggingLevel(badger.ERROR))
	return err
}

func Close() error {
	return db.Close()
}

func GetSize() (int64, error) {
	var size int64
	err := filepath.Walk(db_path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

func CountEntries() (int, error) {
	AssertOpen()
	count := 0
	err := db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			count++
		}
		return nil
	})
	return count, err
}

func Path() string {
	return db_path
}
