package utils

import (
	"sync"
	"time"
)

const (
	workerBits uint8 = 10 // 机器 ID 占用的位数
	seqBits    uint8 = 12 // 序列号占用的位数
)

// Snowflake 结构体
type Snowflake struct {
	mutex       sync.Mutex // 锁，保证并发安全
	startTime   int64      // 启动时间戳（毫秒）
	workerId    uint16     // 机器 ID
	sequenceNum uint16     // 序列号
}

// NewSnowflake construct function
func NewSnowflake(workerId uint16) *Snowflake {
	sf := &Snowflake{}
	sf.startTime = 1693211961606
	sf.workerId = workerId & ((1 << workerBits) - 1)
	return sf
}

func (sf *Snowflake) NextVal() uint64 {
	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	now := getNowMillis() - sf.startTime
	if sf.sequenceNum >= (1<<seqBits)-1 { // 序列号已到达最大值
		time.Sleep(time.Duration(1)) // 等待下一毫秒
		now = getNowMillis() - sf.startTime
		sf.sequenceNum = 0
	}

	sf.sequenceNum++ // 自增序列号
	return uint64(now)<<22 | uint64(sf.workerId)<<12 | uint64(sf.sequenceNum)
}

// 获取当前时间的毫秒表示
func getNowMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
