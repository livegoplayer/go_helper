package utils

import (
	"math"
	"reflect"
)

func EncodeValue(recv reflect.Value, send reflect.Value) {
	switch recv.Interface().(type) {
	case int, int64:
		recv.SetInt(AsInt64(send.Interface()))
	case float32, float64:
		recv.SetFloat(AsFloat64(send.Interface()))
	case string:
		recv.SetString(AsString(send.Interface()))
	default:
		recv.Set(send)
	}
}

func Ceil(n int, c int) int {
	if n%c == 0 {
		return n
	}
	return int(math.Ceil(float64(n)/float64(c))) * c
}
