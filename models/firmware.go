package models

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/google/uuid"
)

const (
	maxSize                      = 1024 * 1024 * 2
	FirmwareCodeHashInvalidValue = 0xffffffff
	FirmwareCodeHashHexLengthMax = 8
)

// Firmware 结构表示Badger中的固件元数据
type Firmware struct {
	ID string `json:"id"` // Badger使用字符串键

	DeviceType          int32                 `json:"device_type"`
	HardwareVersion     int32                 `json:"hardware_version"`
	FirmwareVersion     int32                 `json:"firmware_version"`
	CodeHash            uint32                `json:"code_hash,omitempty"`
	BuildTime           uint32                `json:"build_time,omitempty"`
	EncryptionSupported int32                 `json:"encryption_supported"`
	Path                string                `json:"path,omitempty"`
	Name                string                `json:"name,omitempty"`
	MD5                 string                `json:"md5,omitempty"`
	Size                int                   `json:"size,omitempty"`
	Data                []byte                `json:"-"` // 实际数据存储分离
	File                *multipart.FileHeader `json:"-"`

	ReleasedAt *time.Time `json:"released_at,omitempty"`
}

// Key 生成BadgerDB的复合键
func (f *Firmware) Key() []byte {
	return []byte(fmt.Sprintf("firmware:%s", f.ID))
}

// IndexKey 创建查询索引键
func (f *Firmware) IndexKey() []byte {
	return []byte(fmt.Sprintf("idx:%d:%d:%d",
		f.DeviceType,
		f.HardwareVersion,
		f.FirmwareVersion))
}

// UploadToBadger 实现BadgerDB存储
func (firmware *Firmware) UploadToBadger(ctx context.Context, db *badger.DB) error {
	// 生成唯一ID
	if firmware.ID == "" {
		firmware.ID = uuid.New().String()
	}

	// 计算MD5
	if len(firmware.Data) > 0 {
		m := md5.Sum(firmware.Data)
		firmware.MD5 = hex.EncodeToString(m[:])
	}

	return db.Update(func(txn *badger.Txn) error {
		// 序列化数据
		data, err := json.Marshal(firmware)
		if err != nil {
			return err
		}

		// 存储主记录
		if err := txn.Set(firmware.Key(), data); err != nil {
			return err
		}

		// 创建索引条目
		return txn.Set(firmware.IndexKey(), []byte(firmware.ID))
	})
}

// GetFirmware 从BadgerDB获取固件
func GetFirmware(db *badger.DB, deviceType, hwVer, fwVer int32) (*Firmware, error) {
	var firmware Firmware

	err := db.View(func(txn *badger.Txn) error {
		// 通过索引查找ID
		indexKey := []byte(fmt.Sprintf("idx:%d:%d:%d",
			deviceType, hwVer, fwVer))

		item, err := txn.Get(indexKey)
		if err != nil {
			return err
		}

		id, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		// 获取实际数据
		item, err = txn.Get([]byte("firmware:" + string(id)))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &firmware)
		})
	})

	return &firmware, err
}

// ParseFirmware 保持原逻辑
func (firmware *Firmware) ParseFirmware(data []byte) error {
	if int(binary.LittleEndian.Uint32(data[8:12])) != len(data) {
		return errors.New("invalid firmware data length")
	}
	return nil
}
