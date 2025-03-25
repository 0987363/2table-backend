package models

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/dgraph-io/badger/v4"
)

func CreateKey(prefix string, id string) string {
	return fmt.Sprintf("%s%s", prefix, id)
}

type Badger struct {
	db *badger.DB
}

func NewBadger(dataPath string) (*Badger, error) {
	opt := badger.DefaultOptions(dataPath)
	opt.IndexCacheSize = 100 << 20

	if err := os.MkdirAll(dataPath, 0755); err != nil {
		return nil, Error("Create badger dirctory failed: ", err)
	}

	db, err := badger.Open(opt)
	if err != nil {
		return nil, Error("Open badger failed: ", err)
	}

	return &Badger{db}, nil
}

func (m *Badger) Insert(collection, key string, data interface{}) error {
	return m.db.Update(func(txn *badger.Txn) error {
		bytes, err := json.Marshal(data)
		if err != nil {
			return err
		}
		return txn.Set([]byte(key), bytes)
	})
}

func (m *Badger) Get(collection, key string, result interface{}) error {
	return m.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, result)
		})
	})
}

// 二级索引工具方法示例（创建时间索引）
func CreateTimeIndexKey(prefix string, t time.Time, id string) string {
	return fmt.Sprintf("%s@created:%d:%s", prefix, t.UnixNano(), id)
}
