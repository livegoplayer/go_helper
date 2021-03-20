package stl

import (
	"github.com/livegoplayer/go_helper/utils"
	"reflect"
	"sync"
)

// 队列
type Queue struct {
	q []utils.Any
	sync.RWMutex
	itemType reflect.Type
}

func NewQueue(instance utils.Any) *Queue {
	return &Queue{
		q:        make([]utils.Any, 0),
		itemType: reflect.TypeOf(instance),
	}
}

// 入队列
func (q *Queue) EnQueue(item utils.Any) *Queue {
	q.Lock()
	defer q.Unlock()
	return q.enQueue(item)
}

func (q *Queue) enQueue(item utils.Any) *Queue {
	if q.itemType != reflect.TypeOf(item) {
		//panic("添加元素类型错误")
		return q
	}
	q.q = append(q.q, item)
	return q
}

// 入队列（list）
func (q *Queue) EnQueueList(list utils.Any) *Queue {
	q.Lock()
	defer q.Unlock()
	return q.enQueueList(list)
}

func (q *Queue) enQueueList(list utils.Any) *Queue {
	if kind := reflect.TypeOf(list).Kind(); kind == reflect.Slice {
		l := reflect.ValueOf(list)
		num := l.Len()
		for i := 0; i < num; i++ {
			q.enQueue(l.Index(i).Interface())
		}
	}
	return q
}

// 出队列，空队列返回nil
func (q *Queue) DeQueue() utils.Any {
	q.Lock()
	defer q.Unlock()
	return q.deQueue()
}

func (q *Queue) deQueue() utils.Any {
	if q.isEmpty() {
		return nil
	}
	item := q.q[0]
	q.q = q.q[1:]
	return item
}

// 返回队首元素，空队列返回nil
func (q *Queue) Front() utils.Any {
	q.RLock()
	defer q.RUnlock()
	return q.front()
}

func (q *Queue) front() utils.Any {
	if q.isEmpty() {
		return nil
	}
	return q.q[0]
}

// 队列长度
func (q *Queue) Len() int {
	q.RLock()
	defer q.RUnlock()
	return q.len()
}

func (q *Queue) len() int {
	return len(q.q)
}

// 队列是否为空
func (q *Queue) IsEmpty() bool {
	q.RLock()
	defer q.RUnlock()
	return q.isEmpty()
}

func (q *Queue) isEmpty() bool {
	return q.len() == 0
}

// 队列中是否存在
func (q *Queue) IsExists(item utils.Any) bool {
	q.RLock()
	defer q.RUnlock()
	return q.isExists(item)
}

func (q *Queue) isExists(item utils.Any) bool {
	if q.itemType != reflect.TypeOf(item) {
		return false
	}
	return utils.IsExists(item, q.q)
}

// 重置队列
func (q *Queue) Clean() *Queue {
	q.Lock()
	defer q.Unlock()
	return q.clean()
}

func (q *Queue) clean() *Queue {
	q.q = make([]utils.Any, 0)
	return q
}

func (q *Queue) ToList() []utils.Any {
	q.RLock()
	defer q.RUnlock()
	return q.toList()
}

func (q *Queue) toList() []utils.Any {
	list := make([]utils.Any, q.len())
	copy(list, q.q)
	return list
}
