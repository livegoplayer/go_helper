package helper

import (
	"strconv"
	"unsafe"
)

func String(unknownVal interface{}) string {

	//如果是数字
	if intVal, ok := unknownVal.(int); ok {
		return strconv.Itoa(intVal)
	}

	var stringVal, ok = unknownVal.(string)

	if ok {
		return stringVal
	}
	return ""
}

func Int32(unknownVal interface{}) int32 {
	intVal, ok := unknownVal.(int32)

	if ok {
		return intVal
	}

	return 0
}

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}
