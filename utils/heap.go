package utils

import "container/heap"

type MinHeap interface {
	heap.Interface
	Top() Value
}

type Value interface {
	GetValue() float64
}

func TopK(source []Value, target MinHeap, k int) {
	heap.Init(target)
	for _, item := range source {
		if target.Len() < k {
			heap.Push(target, item)
			continue
		}
		// 最小堆
		if target.Top().GetValue() < item.GetValue() {
			heap.Pop(target)
			heap.Push(target, item)
		}
	}
}
