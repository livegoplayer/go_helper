package stl

import (
	"github.com/livegoplayer/go_helper/utils"
	"reflect"
	"sync"
)

// 随机队列
type RandomQueue struct {
	l          []utils.Any
	unsafeMode bool
	random     *Random
	sync.RWMutex
	itemType reflect.Type
}

func NewRandomQueue(instance utils.Any) *RandomQueue {
	return &RandomQueue{
		l:        make([]utils.Any, 0),
		itemType: reflect.TypeOf(instance),
	}
}

// 非线程安全
func NewRandomQueueUnsafe(instance utils.Any) *RandomQueue {
	return NewRandomQueue(instance).Unsafe()
}

// 启用非线程安全模式
func (q *RandomQueue) Unsafe() *RandomQueue {
	q.Lock()
	defer q.Unlock()
	return q.unsafe()
}

func (q *RandomQueue) unsafe() *RandomQueue {
	q.unsafeMode = true
	q.random = NewRandom()
	return q
}

// 添加元素
func (q *RandomQueue) Add(item utils.Any) *RandomQueue {
	q.Lock()
	defer q.Unlock()
	return q.add(item)
}

func (q *RandomQueue) add(item utils.Any) *RandomQueue {
	if q.itemType != reflect.TypeOf(item) {
		return q
	}
	q.l = append(q.l, item)
	return q
}

// 添加元素（list）
func (q *RandomQueue) AddList(list utils.Any) *RandomQueue {
	q.Lock()
	defer q.Unlock()
	return q.addList(list)
}

func (q *RandomQueue) addList(list utils.Any) *RandomQueue {
	if kind := reflect.TypeOf(list).Kind(); kind == reflect.Slice {
		l := reflect.ValueOf(list)
		num := l.Len()
		for i := 0; i < num; i++ {
			q.add(l.Index(i).Interface())
		}
	}
	return q
}

// 随机读取队列中的一个元素，不删除，不幂等，空队列返回nil
func (q *RandomQueue) Get() utils.Any {
	q.Lock()
	defer q.Unlock()
	return q.get()
}

func (q *RandomQueue) get() utils.Any {
	return q.getNext(false)
}

// 随机取出队列中的一个元素，删除，空队列返回nil
func (q *RandomQueue) Next() utils.Any {
	q.Lock()
	defer q.Unlock()
	return q.next()
}

func (q *RandomQueue) next() utils.Any {
	return q.getNext(true)
}

func (q *RandomQueue) getRandom() int64 {
	if q.unsafeMode {
		return q.random.Int63n(int64(q.len()))
	}
	return utils.GetRandomN(int64(q.len()))
}

func (q *RandomQueue) getNext(remove bool) utils.Any {
	if q.isEmpty() {
		return nil
	}
	i := q.getRandom()
	res := q.l[i]
	q.l[i] = q.l[0]
	if remove {
		q.l = q.l[1:]
	}
	return res
}

// 队列长度
func (q *RandomQueue) Len() int {
	q.RLock()
	defer q.RUnlock()
	return q.len()
}

func (q *RandomQueue) len() int {
	return len(q.l)
}

// 队列是否为空
func (q *RandomQueue) IsEmpty() bool {
	q.RLock()
	defer q.RUnlock()
	return q.isEmpty()
}

func (q *RandomQueue) isEmpty() bool {
	return q.len() == 0
}

// 队列中是否存在
func (q *RandomQueue) IsExists(item utils.Any) bool {
	q.RLock()
	defer q.RUnlock()
	return q.isExists(item)
}

func (q *RandomQueue) isExists(item utils.Any) bool {
	if q.itemType != reflect.TypeOf(item) {
		return false
	}
	return utils.IsExists(item, q.l)
}

// 重置队列
func (q *RandomQueue) Clean() *RandomQueue {
	q.Lock()
	defer q.Unlock()
	return q.clean()
}

func (q *RandomQueue) clean() *RandomQueue {
	q.l = make([]utils.Any, 0)
	return q
}

func (q *RandomQueue) ToList() []utils.Any {
	q.Lock()
	defer q.Unlock()
	return q.toList()
}

func (q *RandomQueue) toList() []utils.Any {
	list := make([]utils.Any, q.len())
	for i := range q.l {
		list[i] = q.next()
	}
	return list
}
