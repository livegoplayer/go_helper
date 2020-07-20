package helper

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"reflect"
	"strings"
	"time"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

//这个函数会破坏file文件，所以不能瞎用
func FileSha1(file io.Reader) string {
	_sha1 := sha1.New()
	_, _ = io.Copy(_sha1, file)
	return hex.EncodeToString(_sha1.Sum(nil))
}

func StrMd5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}

func JsonEncode(data interface{}) []byte {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonByte, err := json.Marshal(&data)
	if err != nil {
		fmt.Printf("json加密出错:" + err.Error())
	}
	return jsonByte
}

//解析valjson对象
func JsonDecode(val []byte, data interface{}) interface{} {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	//把结果放入val
	err := json.Unmarshal(val, &data)
	if err != nil {
		fmt.Printf("json解析出错:" + err.Error())
	}
	return data
}

//直接获取json中的值
func GetJsonVal(data []byte, path ...interface{}) string {
	return jsoniter.Get(data, path...).ToString()
}

/**字符串->时间对象*/
func Str2Time(formatTimeStr string) time.Time {
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, formatTimeStr, loc) //使用模板在对应时区转化为time.time类型

	return theTime

}

/**字符串->时间戳*/
func Str2Stamp(formatTimeStr string) int64 {
	timeStruct := Str2Time(formatTimeStr)
	millisecond := timeStruct.UnixNano() / 1e6
	return millisecond
}

/**时间对象->字符串*/
func Time2Str() string {
	const shortForm = "2006-01-01 15:04:05"
	t := time.Now()
	temp := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)
	str := temp.Format(shortForm)
	return str
}

/*时间对象->时间戳*/
func Time2Stamp() int64 {
	t := time.Now()
	millisecond := t.UnixNano() / 1e6
	return millisecond
}

/*时间戳->字符串*/
func Stamp2Str(stamp int64) string {
	timeLayout := "2006-01-02 15:04:05"
	str := time.Unix(stamp/1000, 0).Format(timeLayout)
	return str
}

/*时间戳->时间对象*/
func Stamp2Time(stamp int64) time.Time {
	stampStr := Stamp2Str(stamp)
	timer := Str2Time(stampStr)
	return timer
}

//获取file对象
func GetFileByHeader(header *multipart.FileHeader) (multipart.File, error) {
	file, err := header.Open()
	if err != nil {
		return nil, err
	}

	return file, nil
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

//获取source的子串,如果start小于0或者end大于source长度则返回""
//start:开始index，从0开始，包括0
//end:结束index，以end结束，但不包括end
func Substring(source string, start int, end int) string {
	var r = []rune(source)
	length := len(r)

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
	source = Substring(source, strings.Index(source, startString)+strings.Count(startString, "")-1, strings.Count(source, "")-1)
	if endString == "" {
		return source
	}
	return Substring(source, 0, strings.Index(source, endString))
}
