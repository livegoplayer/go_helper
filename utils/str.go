package utils

import (
	"fmt"
	"math/rand"
	"regexp"
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

// SubStrByShowLen 模仿 php 中的 mb_strimwidth
// 字显示宽度为 1 或 2
func SubStrByShowLen(s, suffix string, l int) string {
	if len(s) <= l {
		return s
	}
	ss, sl, rl, rs := "", 0, 0, []rune(s)
	suffixl := len(suffix)
	for _, r := range rs {
		rint := int(r) // 获取字节值
		if rint < 128 {
			rl = 1
		} else {
			rl = 2
		}

		if sl+rl+suffixl > l {
			break
		}
		sl += rl
		ss += string(r)
	}

	if sl < suffixl {
		return ss
	}
	return ss + suffix
}

// GetSubStringBetween 获取两个字符串中间的字符串
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

// 生成n位数随机字符串
var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// RandStringBytesMaskImprSrcUnsafe 获取n位的随机字符串
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

// CompressStr 压缩字符串，去除空格、制表符、换页符等字符
func CompressStr(str string) string {
	if str == "" {
		return ""
	}

	// \s 匹配任何空白字符，包括空格、制表符、换页符等等，等价于 \f\n\r\t\v
	// Unicode 正则表达式会匹配全角空格
	re := regexp.MustCompile("\\s+")

	return re.ReplaceAllString(str, "")
}

func HasWhiteSpaceStr(str string) bool {
	for _, v := range str {
		// 32是空格
		if v != 32 {
			return false
		}
	}
	return true
}

// SplitComma 将字符串以逗号分隔开
func SplitComma(str string) []string {
	str = strings.ReplaceAll(str, "，", ",")
	arr := strings.Split(str, ",")
	return FilterEmptyStr(arr)
}

func Number2Chinese(number int64, money ...bool) (chinese string) {
	isMoney := len(money) > 0 && money[0]
	if number == 0 {
		if isMoney {
			return "零圆整"
		}
		return "零"
	}
	if number < 0 {
		number *= -1
		defer func() {
			chinese = fmt.Sprintf("负%s", chinese)
		}()
	}
	chineseMap := []string{"%", "十", "百", "千", "万", "十", "百", "千", "亿", "十", "百", "千"}
	chineseNum := []string{"零", "一", "二", "三", "四", "五", "六", "七", "八", "九"}
	if isMoney {
		chineseMap = []string{"圆整", "拾", "佰", "仟", "万", "拾", "佰", "仟", "亿", "拾", "佰", "仟"}
		chineseNum = []string{"零", "壹", "贰", "叁", "肆", "伍", "陆", "柒", "捌", "玖"}
	} else {
		defer func() {
			chinese = strings.TrimSuffix(chinese, "%")
			if strings.HasPrefix(chinese, "一十") {
				chinese = strings.Replace(chinese, "一十", "十", 1)
			}
		}()
	}
	listNum := make([]int64, 0)
	for ; number > 0; number = number / 10 {
		listNum = append(listNum, number%10)
	}
	for i := len(listNum) - 1; i >= 0; i-- {
		chinese = fmt.Sprintf("%s%s%s", chinese, chineseNum[listNum[i]], chineseMap[i])
	}
	str := ""
	for {
		str = chinese
		str = strings.Replace(str, "零亿", "亿", 1)
		str = strings.Replace(str, "零万", "万", 1)
		if isMoney {
			str = strings.Replace(str, "零拾", "零", 1)
			str = strings.Replace(str, "零佰", "零", 1)
			str = strings.Replace(str, "零仟", "零", 1)
			str = strings.Replace(str, "零零", "零", 1)
			str = strings.Replace(str, "零圆", "圆", 1)
		} else {
			str = strings.Replace(str, "零十", "零", 1)
			str = strings.Replace(str, "零百", "零", 1)
			str = strings.Replace(str, "零千", "零", 1)
			str = strings.Replace(str, "零零", "零", 1)
			str = strings.Replace(str, "零%", "%", 1)
		}
		if str == chinese {
			chinese = str
			break
		}
		chinese = str
	}
	return
}
