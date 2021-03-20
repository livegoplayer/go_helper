package utils

import (
	"math/rand"
	"time"
)

func RandInt64(min, max int64) int64 {
	if min >= max {
		return 0
	}
	rand.Seed(time.Now().UnixNano() + time.Now().Unix())
	return rand.Int63n(max-min) + min
}

func RandSliceIdx(size, len int) []int {
	idxs := make([]int, len)
	for i := range idxs {
		idxs[i] = i
	}
	if size <= len {
		return idxs
	}
	ret := make([]int, size)
	for i := range ret {
		ret[i], idxs = RandSlice(idxs)
	}
	return ret
}

func RandSlice(s []int) (int, []int) {
	idx := rand.Int() % len(s)
	ret := s[idx]
	return ret, append(s[:idx], s[idx+1:]...)
}
