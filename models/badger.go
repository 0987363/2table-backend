package models

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/dgraph-io/badger/v4"
)

const (
	FileCollection = "file"
	UserCollection = "user"
)

func GenKey(collection, id string) string {
	return fmt.Sprintf("%s:%s", collection, id)
}

type Badger struct {
	db *badger.DB
}

func NewBadger(dataPath string) (*Badger, error) {
	opts := badger.DefaultOptions(dataPath)
	opts.IndexCacheSize = 10 << 20    // 10MB
	opts.ValueLogFileSize = 128 << 20 // 128MB

	if err := os.MkdirAll(dataPath, 0755); err != nil {
		return nil, Error("Create badger dirctory failed: ", err)
	}

	db, err := badger.Open(opts)
	if err != nil {
		return nil, Error("Open badger failed: ", err)
	}

	return &Badger{db}, nil
}

func (m *Badger) InsertFile(collection, key string, data *File) error {
	return m.db.Update(func(txn *badger.Txn) error {
		bytes, err := json.Marshal(data)
		if err != nil {
			return err
		}
		return txn.Set([]byte(GenKey(collection, key)), bytes)
	})
}

func (m *Badger) ListFile(collection string, result *[]*File) error {
	return m.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(collection)

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			res := File{}
			item := it.Item()
			fmt.Println("key: ", string(item.Key()))
			if err := item.Value(func(val []byte) error {
				return json.Unmarshal(val, &res)
			}); err != nil {
				return err
			}
			*result = append(*result, &res)
		}
		return nil
	})
}

func (m *Badger) ListAll(collection string, result *[]*File) error {
	return m.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(nil); it.Valid(); it.Next() {
			res := File{}
			item := it.Item()
			if err := item.Value(func(val []byte) error {
				return json.Unmarshal(val, &res)
			}); err != nil {
				return err
			}
			*result = append(*result, &res)
		}
		return nil
	})
}

func (m *Badger) GetFile(collection, key string, result *File) error {
	return m.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(GenKey(collection, key)))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, result)
		})
	})
}

func (m *Badger) DeleteFile(collection, key string) error {
	return m.db.Update(func(txn *badger.Txn) error {
		if err := txn.Delete([]byte(GenKey(collection, key))); err != nil {
			return err
		}
		return nil
	})
}

// 二级索引工具方法示例（创建时间索引）
func CreateTimeIndexKey(prefix string, t time.Time, id string) string {
	return fmt.Sprintf("%s@created:%d:%s", prefix, t.UnixNano(), id)
}
