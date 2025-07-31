package utility

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"
)

var (
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
	tmp    = uint32(random.Intn(100))
)

// GetSeqNum 获取一个seq_num
func GetSeqNum() string {
	now := time.Now()
	return fmt.Sprintf("%d@%d@%d", now.UnixNano(), random.Uint32(), atomic.AddUint32(&tmp, 1))
}