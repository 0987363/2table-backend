package models

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

const (
	epoch          = int64(1609459200000)         // 自定义epoch：2021-01-01 00:00:00 UTC
	machineIDBits  = uint(10)                     // 机器ID位数
	sequenceBits   = uint(12)                     // 序列号位数
	machineIDShift = sequenceBits                 // 机器ID左移位数
	timestampShift = sequenceBits + machineIDBits // 时间戳左移位数
	maxMachineID   = -1 ^ (-1 << machineIDBits)   // 最大机器ID
	maxSequence    = -1 ^ (-1 << sequenceBits)    // 最大序列号
)

type Snowflake struct {
	mu        sync.Mutex
	lastTime  int64
	machineID int64
	sequence  int64
}

var snowflake *Snowflake

func init() {
	snowflake, _ = NewSnowflake(0)
}

func GetID() string {
	id := snowflake.Generate()
	return strconv.FormatInt(id, 10)
}

func NewSnowflake(machineID int64) (*Snowflake, error) {
	if machineID < 0 || machineID > maxMachineID {
		return nil, fmt.Errorf("machineID必须在0到%d之间", maxMachineID)
	}
	return &Snowflake{machineID: machineID}, nil
}

func (s *Snowflake) Generate() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixNano() / 1e6 // 毫秒时间戳
	if now < s.lastTime {
		panic("时钟回拨")
	}

	if now == s.lastTime {
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 { // 当前毫秒序列号用尽
			for now <= s.lastTime {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTime = now
	id := (now-epoch)<<timestampShift | (s.machineID << machineIDShift) | s.sequence
	return id
}
