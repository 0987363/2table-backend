package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
)

// 基础模型结构（Badger使用字符串ID）
type Base struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 键前缀常量（替代MongoDB的集合概念）
const (
	// 固件相关
	FirmwarePrefix        = "firmware:"
	FirmwareReleasePrefix = "firmware_release:"

	// 用户与组织
	UserPrefix    = "user:"
	CompanyPrefix = "company:"
	OrderPrefix   = "order:"
)

// 键生成工具
func CreateKey(prefix string, id string) string {
	return fmt.Sprintf("%s%s", prefix, id)
}

// 基础操作方法
type BaseModel struct {
	db *badger.DB
}

func NewBaseModel(db *badger.DB) *BaseModel {
	return &BaseModel{db: db}
}

// 通用保存方法
func (m *BaseModel) Save(key string, data interface{}) error {
	return m.db.Update(func(txn *badger.Txn) error {
		bytes, err := json.Marshal(data)
		if err != nil {
			return err
		}
		return txn.Set([]byte(key), bytes)
	})
}

// 通用获取方法
func (m *BaseModel) Get(key string, result interface{}) error {
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

// 事务操作方法示例
func (m *BaseModel) UpdateWithTransaction(key string, updateFn func(current []byte) ([]byte, error)) error {
	return m.db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}

		var current []byte
		if err == nil {
			current, err = item.ValueCopy(nil)
			if err != nil {
				return err
			}
		}

		newValue, err := updateFn(current)
		if err != nil {
			return err
		}

		return txn.Set([]byte(key), newValue)
	})
}
