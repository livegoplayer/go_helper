package utils

import (
	"bytes"
	"container/heap"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/livegoplayer/go_helper/collect/collection"
	"github.com/livegoplayer/go_helper/mapstructure"
	"github.com/livegoplayer/go_helper/structs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"math"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func VersionThan(va string, vb string) int {
	vas := strings.Split(va, ".")
	vbs := strings.Split(vb, ".")

	count := len(vas)
	if len(vas) < len(vbs) {
		count = len(vbs)
	}
	for i := 0; i < count; i++ {
		a := 0
		if len(vas) < i+1 {
			a = 0
		} else {
			a, _ = strconv.Atoi(vas[i])
		}

		b := 0
		if len(vbs) < i+1 {
			b = 0
		} else {
			b, _ = strconv.Atoi(vbs[i])
		}
		if a > b {
			return 1
		}
		if a < b {
			return -1
		}
	}
	return 0
}

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Base64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func Sha256(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// AppendStuct 将give的结构体中，同名的成员赋值给reciv
func AppendStuct(reciv Any, give Any) {
	gtyp := reflect.TypeOf(give).Elem()
	gval := reflect.ValueOf(give).Elem()

	rtyp := reflect.TypeOf(reciv).Elem()
	rval := reflect.ValueOf(reciv).Elem()

	for i := 0; i < gtyp.NumField(); i++ {
		gname := gtyp.Field(i).Name
		for j := 0; j < rtyp.NumField(); j++ {
			if gname == rtyp.Field(j).Name {
				rval.Field(j).Set(gval.Field(i))
				break
			}
		}
	}
}

// 将m的值，复制给s结构体中同名json-tag的字段，s必须为结构体指针
// 复合结构使用 mapstructure.ToStruct() 代替
func ToStruct(m H, s interface{}, tagNames ...string) {
	styp := reflect.TypeOf(s).Elem()
	sval := reflect.ValueOf(s).Elem()

	tagName := "json"
	if len(tagNames) > 0 {
		tagName = tagNames[0]
	}

	for i := 0; i < styp.NumField(); i++ {
		field := styp.Field(i)
		inputFieldName := field.Tag.Get(tagName)
		if inputFieldName == "" {
			continue
		}
		parses := strings.Split(inputFieldName, ",")
		if len(parses) > 1 {
			inputFieldName = parses[0]
		}
		v, ok := m[inputFieldName]
		if !ok || v == nil {
			continue
		}
		EncodeValue(sval.Field(i), reflect.ValueOf(v))
	}
}

func ToStructV2(m interface{}, s interface{}, tagNames ...string) {
	if str, ok := m.(string); ok {
		JsonToStructV2(str, s, tagNames...)
		return
	}
	// 非字符串内容，由他自己处理
	mapstructure.ToStructV2(m, s, tagNames...)
}

// 将json字符串，复制给s结构体中同名json-tag的字段，s必须为结构体指针
func JsonToStruct(json string, s interface{}, tagNames ...string) {
	m := JsonDecodeSafe(json)
	if m == nil {
		return
	}
	ToStruct(m, s, tagNames...)
}

// 将json字符串，使用mapstruct实现
func JsonToStructV2(json string, s interface{}, tagNames ...string) {
	m := JsonDecodeSafe(json)
	if m == nil {
		return
	}
	mapstructure.ToStruct(m, s, tagNames...)
}

// 将m的值，复制给s结构体中同名json-tag的字段，s必须为结构体指针
func Merge(m H, h H) H {
	for k, v := range h {
		m[k] = v
	}
	return m
}

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

func EncodeUrlWithoutSort(v url.Values, keys []string, withOutEncode bool) string {
	if v == nil {
		return ""
	}
	var buf string
	for _, k := range keys {
		vs, ok := v[k]
		if !ok {
			continue
		}
		keyEscaped := url.QueryEscape(k)
		if withOutEncode {
			keyEscaped = k
		}
		for _, vl := range vs {
			if len(buf) > 0 {
				buf += "&"
			}
			buf += keyEscaped
			buf += "="
			if withOutEncode {
				buf += vl
			} else {
				buf += url.QueryEscape(vl)
			}
		}
	}
	res := buf
	return res
}

func JsonDecodeSafe(jsonStr string) map[string]interface{} {
	var mapResult map[string]interface{}
	ds := json.NewDecoder(strings.NewReader(jsonStr))
	ds.UseNumber()
	_ = ds.Decode(&mapResult)
	return mapResult
}

func JsonDecodes(jsonStr string) (map[string]interface{}, []map[string]interface{}) {
	if jsonStr[0] == '[' {
		mapSlice := make([]map[string]interface{}, 0)
		_ = json.Unmarshal([]byte(jsonStr), &mapSlice)
		return nil, mapSlice
	}
	var mapResult map[string]interface{}
	_ = json.Unmarshal([]byte(jsonStr), &mapResult)
	return mapResult, nil
}

func JsonDecodeAny(jsonStr string) interface{} {
	var any interface{}
	_ = json.Unmarshal([]byte(jsonStr), &any)
	return any
}

func InArrayByFunc(val interface{}, array interface{}, equalFunc func(val, arrayItem interface{}) bool) (exist bool, index int) {
	exist = false
	index = -1
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if equalFunc(val, s.Index(i).Interface()) {
				index = i
				exist = true
				return
			}
		}
	}
	return
}

func IsExists(val interface{}, array interface{}) bool {
	e, _ := InArray(val, array)
	return e
}

func IsExistsByFunc(val interface{}, array interface{}, equalFunc func(val, arrayItem interface{}) bool) bool {
	e, _ := InArrayByFunc(val, array, equalFunc)
	return e
}

func IndexOf(val interface{}, array interface{}) int {
	_, i := InArray(val, array)
	return i
}

func IndexOfByFunc(val interface{}, array interface{}, equalFunc func(val, arrayItem interface{}) bool) int {
	_, i := InArrayByFunc(val, array, equalFunc)
	return i
}

func Deduplication(arr []int64) []int64 {
	helper := make(map[int64]bool)
	res := make([]int64, 0)
	for _, item := range arr {
		if _, ok := helper[item]; !ok {
			res = append(res, item)
			helper[item] = true
		}
	}
	return res
}

// 将下划线风格的单词变为驼峰命名的单词
func UnderLineToCamel(line string) string {
	words := strings.Split(line, "_")
	n := ""
	for _, w := range words {
		n += strings.ToUpper(w[0:1]) + w[1:]
	}
	return n
}

func MaxNum(argus ...int64) int64 {
	max := argus[0]
	for _, i := range argus {
		if i > max {
			max = i
		}
	}
	return max
}

func MinNum(argus ...int64) int64 {
	min := argus[0]
	for _, i := range argus {
		if i < min {
			min = i
		}
	}
	return min
}

func MergeMap(argus ...map[string]interface{}) map[string]interface{} {
	ret := make(map[string]interface{})
	for _, m := range argus {
		for k, v := range m {
			ret[k] = v
		}
	}
	return ret
}

func HttpGet(url string) (map[string]interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	rep := JsonDecodeToMap(string(body))
	return rep, nil
}

func TimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, cTimeout)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(rwTimeout))
		return conn, nil
	}
}

func MapKeys(item map[string]interface{}) []string {
	ks := make([]string, 0)
	for k := range item {
		ks = append(ks, k)
	}
	return ks
}

func HttpPost(url string, params interface{}, seconds int, header ...interface{}) (map[string]interface{}, error) {
	body := JsonEncode(params)

	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	if len(header) > 0 {
		headers := header[0].(map[string]string)
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	connectTimeout := time.Duration(seconds) * time.Second
	readWriteTimeout := time.Duration(seconds) * time.Second
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Dial:            TimeoutDialer(connectTimeout, readWriteTimeout),
		},
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	rep := JsonDecodeToMap(string(content))
	return rep, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - (len(ciphertext) % blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func AesCbcEncrypt(str string, key []byte, IV []byte) ([]byte, error) {
	origData := []byte(str)

	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := cipherBlock.BlockSize()
	origData = PKCS5Padding(origData, blockSize)

	crypted := make([]byte, len(origData))
	cipher.NewCBCEncrypter(cipherBlock, IV).CryptBlocks(crypted, origData)
	return crypted, nil
}

func AesCbcDecrypt(encrypted []byte, key []byte, IV []byte) ([]byte, error) {
	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	decrypt := make([]byte, len(encrypted))

	cipher.NewCBCDecrypter(cipherBlock, IV).CryptBlocks(decrypt, encrypted)
	return decrypt, nil
}

// ToMap converts a struct to a map using the struct's tags.
//
// ToMap uses tags on struct fields to decide which fields to add to the
// returned map.
func ToMap(in interface{}, tag string) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ToMap only accepts structs; got %T", v)
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		if tagv := fi.Tag.Get(tag); tagv != "" {
			// set key of map to value in struct field
			out[tagv] = v.Field(i).Interface()
		}
	}
	return out, nil
}

func DeepGet(m map[string]interface{}, key string) (interface{}, bool) {
	return collection.DeepGet(m, key)
}

func DeepGetMust(m map[string]interface{}, key string) interface{} {
	v, ok := collection.DeepGet(m, key)
	if !ok {
		panic("没有该值")
	}
	return v
}

func DeepGetShould(m map[string]interface{}, key string) interface{} {
	v, _ := collection.DeepGet(m, key)
	return v
}

func DeepMustSet(m map[string]interface{}, key string, val interface{}) {
	succ := collection.DeepSet(m, key, val)
	if !succ {
		panic("设置失败")
	}
}

// 模仿 php 中的 mb_strimwidth
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

func MaxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func MinInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func SliceString2Int(src []string) ([]int, error) {
	dst := make([]int, 0, len(src))
	for i := range src {
		n, err := strconv.Atoi(src[i])
		if err != nil {
			return nil, err
		}
		dst = append(dst, n)
	}
	return dst, nil
}

// 字符串切片转换为 int64 切片
func SliceString2Int64(src []string) ([]int64, error) {
	tmp, err := SliceString2Int(src)
	if err != nil {
		return nil, err
	}

	dst := make([]int64, 0, len(tmp))
	for _, i := range tmp {
		dst = append(dst, int64(i))
	}
	return dst, nil
}

func CopyMapTopLevel(src H) H {
	res := make(H, len(src))
	for k, v := range src {
		res[k] = v
	}
	return res
}

func Diff(arr1 []int64, arr2 []int64) []int64 {
	ret := make([]int64, 0)
	for _, i := range arr1 {
		if !IsExists(i, arr2) {
			ret = append(ret, i)
		}
	}
	return ret
}

type int64list []int64

func (m int64list) Len() int {
	return len(m)
}
func (m int64list) Less(i, j int) bool {
	return m[i] < m[j]
}
func (m int64list) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func Sort(arr []int64) []int64 {
	s := int64list(arr)
	sort.Sort(s)
	return s
}

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

// "s1, s2, s3, s4" => []string{"s1", "s2", "s3", "s4"}
// "[s1, s2, s3, s4]" => []string{"s1", "s2", "s3", "s4"}
// "["s1", "s2", "s3", "s4"]" => []string{"s1", "s2", "s3", "s4"}
func StrToStrSlice(str string) []string {
	if str == "" {
		return nil
	}
	if str[0] == '[' && str[len(str)-1] == ']' {
		str = str[1 : len(str)-1]
	}
	strs := strings.Split(str, ",")
	ids := make([]string, len(strs))
	for i, s := range strs {
		ids[i] = strings.Trim(strings.Trim(s, " "), "\"")
	}
	return ids
}

func SliceToStr(arr []int64) string {
	builder := strings.Builder{}
	for i, num := range arr {
		builder.WriteString(strconv.Itoa(int(num)))
		if i != len(arr)-1 {
			builder.WriteString(",")
		}
	}
	return builder.String()
}

func StrSliceToStr(arr []string) string {
	builder := strings.Builder{}
	for _, str := range arr {
		if str == "" {
			continue
		}
		builder.WriteString(str)
		builder.WriteString(",")
	}

	res := builder.String()
	if res != "" && res[len(res)-1] == ',' {
		res = res[:len(res)-1]
	}
	return res
}

// "1,2,3,4" => []int{1, 2, 3, 4}
// "[1,2,3,4]" => []int{1, 2, 3, 4}
func StrToSlice(str string) []int64 {
	if str == "" {
		return nil
	}
	if str[0] == '[' && str[len(str)-1] == ']' {
		str = str[1 : len(str)-1]
	}
	strs := strings.Split(str, ",")
	ids := make([]int64, len(strs))
	for i, s := range strs {
		ids[i] = AsInt64(s)
	}
	return ids
}

func MaxInSlice(arr []int64) int64 {
	max := int64(math.MinInt64)
	for _, v := range arr {
		if v > max {
			max = v
		}
	}
	return max
}

func RemoveNum(arr []int64, target int64) []int64 {
	index := -1
	for i, v := range arr {
		if v == target {
			index = i
			break
		}
	}
	if index == -1 {
		return arr
	}
	if index == len(arr)-1 {
		return arr[:len(arr)-1]
	}
	return append(arr[:index], arr[index+1:]...)
}

// 必须dest传入指针
func Struct2Struct(src, dest interface{}) {
	mapstructure.ToStruct(structs.Map(src), dest)
}

func GetCurPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	rst := filepath.Dir(path)
	return rst
}

// 删除交集
func RemoveIntersection(arr1, arr2 []int64) []int64 {
	arr2m := make(map[int64]struct{})
	for _, i := range arr2 {
		arr2m[i] = struct{}{}
	}
	result := make([]int64, 0)
	for _, i := range arr1 {
		if _, ok := arr2m[i]; !ok {
			result = append(result, i)
		}
	}
	return result
}

func Ceil(n int, c int) int {
	if n%c == 0 {
		return n
	}
	return int(math.Ceil(float64(n)/float64(c))) * c
}

// 获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

// 获取当前 时分秒 的秒数
func GetTimeSeconds() int64 {
	now := time.Now()
	return int64(now.Hour()*3600 + now.Minute()*60 + now.Second())
}

func GetFileExt(path string) string {
	for i := len(path) - 1; i >= 0 && !os.IsPathSeparator(path[i]); i-- {
		if path[i] == '.' {
			return path[i+1:]
		}
	}
	return ""
}

// 按显示长度截取字符串
func ShowSubstr(s string, l int) string {
	if len(s) <= l {
		return s
	}
	ss, sl, rl, rs := "", 0, 0, []rune(s)
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			rl = 1
		} else {
			rl = 2
		}

		if sl+rl > l {
			break
		}
		sl += rl
		ss += string(r)
	}
	return ss
}

func GetHashCode(str string, funcType ...HashFunc) int64 {
	hashCode := BKDRHashFunc
	if len(funcType) > 0 {
		hashCode = funcType[0]
	}
	return hashCode(str)
}

// 切片中的所有元素是否都满足指定方法
func All(array interface{}, function func(item interface{}) bool) bool {
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if !function(s.Index(i).Interface()) {
				return false
			}
		}
	default:
		return false
	}
	return true
}

// 切片中是否存在满足指定方法的元素，如果不传方法，则判断切片中是否有元素，传多个方法只会取第一个
func AnyOne(array interface{}, function ...func(item interface{}) bool) bool {
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		f := func(_ interface{}) bool { return true }
		if len(function) > 0 {
			f = function[0]
		}
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if f(s.Index(i).Interface()) {
				return true
			}
		}
	}
	return false
}

// 将切片中的每个元素都作为入参传入指定方法中，收集方法返回值并放入切片返回
func Select(array interface{}, function func(item interface{}) interface{}) []interface{} {
	res := make([]interface{}, 0)
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			res = append(res, function(s.Index(i).Interface()))
		}
	}
	return res
}

// 返回切片中满足指定方法的元素个数
func Count(array interface{}, function func(item interface{}) bool) int64 {
	var res int64
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if function(s.Index(i).Interface()) {
				res++
			}
		}
	}
	return res
}

// 返回切片中满足指定方法的元素
func Where(array interface{}, function func(item interface{}) bool) []interface{} {
	res := make([]interface{}, 0)
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if function(s.Index(i).Interface()) {
				res = append(res, s.Index(i).Interface())
			}
		}
	}
	return res
}

// 返回切片中第一个满足指定方法的元素，如不存在则返回nil
func First(array interface{}, function func(item interface{}) bool) interface{} {
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if function(s.Index(i).Interface()) {
				return s.Index(i).Interface()
			}
		}
	}
	return nil
}

// 返回切片中第一个满足指定方法的元素，如不存在则返回入参中的默认值
func FirstOrDefault(array interface{}, function func(item interface{}) bool, def interface{}) interface{} {
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if function(s.Index(i).Interface()) {
				return s.Index(i).Interface()
			}
		}
	}
	return def
}

// 压缩字符串，去除空格、制表符、换页符等字符
func CompressStr(str string) string {
	if str == "" {
		return ""
	}

	// \s 匹配任何空白字符，包括空格、制表符、换页符等等，等价于 \f\n\r\t\v
	// Unicode 正则表达式会匹配全角空格
	re := regexp.MustCompile("\\s+")

	return re.ReplaceAllString(str, "")
}

// You can get a substring of a UTF-8 string without allocating additional memory (you don't have to convert it to a rune slice)
// @from: https://stackoverflow.com/questions/28718682/how-to-get-a-substring-from-a-string-of-runes-in-golang
func Substr(s string, start int, end int) string {
	var a, b int
	for i := range s {
		if b == start {
			a = i
		}
		if b == end {
			return s[a:i]
		}
		b++
	}

	return s[a:]
}

// Range(n) return [0,1,...,n-1,n] with []int;
// Range(n,m) return [n,n+1,...,m-1,m] with []int
func Range(first int, second ...int) []int {
	var start, end int
	if len(second) == 0 {
		end = first
	} else {
		start, end = first, second[0]
	}
	if start > end {
		return []int{}
	}
	res := make([]int, end-start+1)
	for i := start; i <= end; i++ {
		res[i-start] = i
	}
	return res
}

// Range64(n) return [0,1,...,n-1,n] with []int64;
// Range64(n,m) return [n,n+1,...,m-1,m] with []int64
func Range64(first int64, second ...int64) []int64 {
	var start, end int64
	if len(second) == 0 {
		end = first
	} else {
		start, end = first, second[0]
	}
	if start > end {
		return []int64{}
	}
	res := make([]int64, end-start+1)
	for i := start; i <= end; i++ {
		res[i-start] = i
	}
	return res
}

func Bool2Int(b bool) int {
	if b {
		return 1
	}
	return 0
}

// 驼峰转蛇形
func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

func StrInSlice(ss []string, str string) (exist bool) {
	if len(ss) == 0 {
		return
	}

	for _, s := range ss {
		if s == str {
			exist = true
			break
		}
	}

	return
}

// 获取 2 个切片的前后不同之处
func DiffStringSlice(before, after []string) (add, del, keep []string) {
	for _, str := range after {
		// 不在之前的切片中则为新增
		if !StrInSlice(before, str) {
			add = append(add, str)
		}
	}

	for _, str := range before {
		// 不在之后的切片中则为删掉
		if !StrInSlice(after, str) {
			del = append(del, str)
		} else {
			// 保持不变的
			keep = append(keep, str)
		}
	}

	return
}

func IsEmptyStr(str string) bool {
	return str != ""
}

func IsWhiteSpaceStr(str string) bool {
	for _, v := range str {
		// 32是空格
		if v != 32 {
			return false
		}
	}
	return true
}

func GetRandom(seed ...interface{}) int64 {
	str := primitive.NewObjectID().Hex()
	for _, v := range seed {
		str += AsString(v)
	}
	rand.Seed(GetHashCode(str))
	return rand.Int63()
}

func GetRandomN(n int64, seed ...interface{}) int64 {
	str := primitive.NewObjectID().Hex()
	for _, v := range seed {
		str += AsString(v)
	}
	rand.Seed(GetHashCode(str))
	return rand.Int63n(n)
}

// 过滤字符串数组中的空字符串
func FilterEmptyStr(arr []string) []string {
	ret := make([]string, 0)
	for _, i := range arr {
		if strings.TrimSpace(i) == "" {
			continue
		}
		ret = append(ret, strings.TrimSpace(i))
	}
	return ret
}

// 将字符串以逗号分隔开
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
