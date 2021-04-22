package utils

import (
	"math/rand"
	"strings"
	"time"
	"unsafe"
)

func Substring(source string, start int64, end int64) string {
	var r = []rune(source)
	length := int64(len(r))

	if start < 0 || end > length || start > end {
		return ""
	}

	if start == 0 && end == length {
		return source
	}

	return string(r[start:end])
}

//获取两个字符串中间的字符串
func GetSubStringBetween(source string, startString string, endString string) string {
	//先拿到第一个字符串到最后的子字符串
	start := int64(strings.Index(source, startString))
	startIndex := int64(0)
	if start != int64(-1) {
		startIndex = start + int64(strings.Count(startString, "")) - 1
	}
	source = Substring(source, startIndex, int64(strings.Count(source, "")-1))
	if endString == "" {
		return source
	}
	return Substring(source, 0, int64(strings.Index(source, endString)))
}

//生成n位数随机字符串
var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

//获取n位的随机字符串
func RandStringBytesMaskImprSrcUnsafe(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
