package stl

import (
	"github.com/google/uuid"
	"github.com/livegoplayer/go_helper/utils"
	"math/rand"
)

// rand包使用互斥锁全局共享一个rand.Rand对象来提供随机函数，在多协程情况下因争夺互斥锁，性能消耗大
// 在多协程场景下可为每个协程创建一个rand.Rand对象，提高效率

// 非线程安全
type Random struct {
	*rand.Rand
}

func NewRandom() *Random {
	id, _ := uuid.NewUUID()
	return NewRandomBySeed(utils.GetHashCode(string(id[:])))
}

func NewRandomBySeed(seed int64) *Random {
	return &Random{rand.New(rand.NewSource(seed))}
}
