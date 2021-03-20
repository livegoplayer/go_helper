package stl

import (
	"github.com/livegoplayer/go_helper/utils"
	"reflect"
	"sync"
)

// 栈
type Stack struct {
	s []utils.Any
	sync.RWMutex
	itemType reflect.Type
}

func NewStack(instance utils.Any) *Stack {
	return &Stack{
		s:        make([]utils.Any, 0),
		itemType: reflect.TypeOf(instance),
	}
}

// 入栈
func (s *Stack) Add(item utils.Any) *Stack {
	s.Lock()
	defer s.Unlock()
	return s.add(item)
}

func (s *Stack) add(item utils.Any) *Stack {
	if s.itemType != reflect.TypeOf(item) {
		//panic("添加元素类型错误")
		return s
	}
	s.s = append(s.s, item)
	return s
}

// 入栈（list）
func (s *Stack) AddList(list utils.Any) *Stack {
	s.Lock()
	defer s.Unlock()
	return s.addList(list)
}

func (s *Stack) addList(list utils.Any) *Stack {
	if kind := reflect.TypeOf(list).Kind(); kind == reflect.Slice {
		l := reflect.ValueOf(list)
		num := l.Len()
		for i := 0; i < num; i++ {
			s.add(l.Index(i).Interface())
		}
	}
	return s
}

// 出栈，空栈返回nil
func (s *Stack) Pop() utils.Any {
	s.Lock()
	defer s.Unlock()
	return s.pop()
}

func (s *Stack) pop() utils.Any {
	if s.isEmpty() {
		return nil
	}
	item := s.s[s.len()-1]
	s.s = s.s[:s.len()-1]
	return item
}

// 返回栈顶元素，不出栈，空栈返回nil
func (s *Stack) Front() utils.Any {
	s.RLock()
	defer s.RUnlock()
	return s.front()
}

func (s *Stack) front() utils.Any {
	if s.isEmpty() {
		return nil
	}
	return s.s[s.Len()-1]
}

// 栈长度
func (s *Stack) Len() int {
	s.RLock()
	defer s.RUnlock()
	return s.len()
}

func (s *Stack) len() int {
	return len(s.s)
}

// 栈是否为空
func (s *Stack) IsEmpty() bool {
	s.RLock()
	defer s.RUnlock()
	return s.isEmpty()
}

func (s *Stack) isEmpty() bool {
	return s.len() == 0
}

// 栈中是否存在
func (s *Stack) IsExists(item utils.Any) bool {
	s.RLock()
	defer s.RUnlock()
	return s.isExists(item)
}

func (s *Stack) isExists(item utils.Any) bool {
	if s.itemType != reflect.TypeOf(item) {
		return false
	}
	return utils.IsExists(item, s.s)
}

// 重置栈
func (s *Stack) Clean() *Stack {
	s.Lock()
	defer s.Unlock()
	return s.clean()
}

func (s *Stack) clean() *Stack {
	s.s = make([]utils.Any, 0)
	return s
}

func (s *Stack) ToList() []utils.Any {
	s.RLock()
	defer s.RUnlock()
	return s.toList()
}

func (s *Stack) toList() []utils.Any {
	list := make([]utils.Any, s.len())
	copy(list, s.s)
	return list
}
