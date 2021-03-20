package stl

import (
	"github.com/livegoplayer/go_helper/utils"
	"reflect"
	"sort"
	"sync"
)

// 集合
type Set struct {
	m map[utils.Any]bool
	sync.RWMutex
	itemType reflect.Type
}

// 传入不可比较类型时返回nil
func NewSet(instance utils.Any) *Set {
	canEquals := false
	func() {
		defer func() {
			recover()
		}()
		canEquals = instance == instance
	}()
	if !canEquals {
		return nil
	}
	return &Set{
		m:        map[utils.Any]bool{},
		itemType: reflect.TypeOf(instance),
	}
}

// 添加
func (s *Set) Add(item utils.Any) *Set {
	if s == nil {
		return s
	}
	s.Lock()
	defer s.Unlock()
	s.add(item)
	return s
}

func (s *Set) add(item utils.Any) bool {
	if s.itemType != reflect.TypeOf(item) {
		return false
	}
	s.m[item] = true
	return true
}

// 添加（list）
func (s *Set) AddList(list utils.Any) *Set {
	if s == nil {
		return s
	}
	s.Lock()
	defer s.Unlock()
	return s.addList(list)
}

func (s *Set) addList(list utils.Any) *Set {
	if kind := reflect.TypeOf(list).Kind(); kind == reflect.Slice {
		l := reflect.ValueOf(list)
		num := l.Len()
		for i := 0; i < num; i++ {
			s.add(l.Index(i).Interface())
		}
	}
	return s
}

// 尝试添加
func (s *Set) TryAdd(item utils.Any) bool {
	if s == nil {
		return false
	}
	s.Lock()
	defer s.Unlock()
	return s.tryAdd(item)
}

func (s *Set) tryAdd(item utils.Any) bool {
	if s.isExists(item) {
		return false
	}
	s.add(item)
	return true
}

// 移除
func (s *Set) Remove(item utils.Any) *Set {
	if s == nil {
		return s
	}
	s.Lock()
	defer s.Unlock()
	return s.remove(item)
}

func (s *Set) remove(item utils.Any) *Set {
	delete(s.m, item)
	return s
}

// 集合中是否存在
func (s *Set) IsExists(item utils.Any) bool {
	if s == nil {
		return false
	}
	s.RLock()
	defer s.RUnlock()
	return s.isExists(item)
}

func (s *Set) isExists(item utils.Any) bool {
	if s.itemType != reflect.TypeOf(item) {
		return false
	}
	_, ok := s.m[item]
	return ok
}

// 集合长度
func (s *Set) Len() int {
	if s == nil {
		return 0
	}
	s.RLock()
	defer s.RUnlock()
	return s.len()
}

func (s *Set) len() int {
	return len(s.m)
}

// 集合是否为空
func (s *Set) IsEmpty() bool {
	if s == nil {
		return true
	}
	s.RLock()
	defer s.RUnlock()
	return s.isEmpty()
}

func (s *Set) isEmpty() bool {
	return s.len() == 0
}

// 重置集合
func (s *Set) Clean() *Set {
	if s == nil {
		return s
	}
	s.Lock()
	defer s.Unlock()
	return s.clean()
}

func (s *Set) clean() *Set {
	s.m = make(map[utils.Any]bool)
	return s
}

// 不保证以传入顺序返回,不幂等
func (s *Set) ToList() []utils.Any {
	if s == nil {
		return make([]utils.Any, 0)
	}
	s.RLock()
	defer s.RUnlock()
	return s.toList()
}

func (s *Set) toList() []utils.Any {
	list := make([]utils.Any, 0, s.len())
	for k := range s.m {
		list = append(list, k)
	}
	return list
}

func (s *Set) ToSortList(less func(i, j int) bool) []utils.Any {
	list := s.ToList()
	sort.Slice(list, less)
	return list
}
