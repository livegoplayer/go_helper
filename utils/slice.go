package utils

import (
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// 定义一个新的类型：element
type element interface {
	// element 支持如下类型
	string | int8 | int16 | int32 | int64 | int | float32 | float64 | uint | uint8 | uint16 | uint32 | uint64
}

type numElement interface {
	int8 | int16 | int32 | int64 | int | float32 | float64 | uint | uint8 | uint16 | uint32 | uint64
}

// Unique 数组或slice边去重
func Unique[T element](arr []T) []T {
	tmp := make(map[T]struct{})
	l := len(arr)
	if l == 0 {
		return arr
	}

	rel := make([]T, 0, l)
	for _, item := range arr {
		_, ok := tmp[item]
		if ok {
			continue
		}
		tmp[item] = struct{}{}
		rel = append(rel, item)
	}

	return rel[:len(tmp)]
}

type MyMAP[KEY element, VALUE element | struct{} | bool | byte] map[KEY]VALUE

// UniqueByMap RemoveRepByMap 通过map主键唯一的特性过滤重复元素
func UniqueByMap[T element](slc []T) []T {
	var result []T
	var tempMap MyMAP[T, byte] = map[T]byte{}
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l {
			result = append(result, e)
		}
	}
	return result
}

func Diff[T int | string](a []T, b []T) []T {
	var diffArray []T
	var temp MyMAP[T, struct{}] = map[T]struct{}{}
	for _, val := range b {
		if _, ok := temp[val]; !ok {
			temp[val] = struct{}{}
		}
	}

	for _, val := range a {
		if _, ok := temp[val]; !ok {
			diffArray = append(diffArray, val)
		}
	}

	return diffArray
}

// Intersect 求两个切片的交集
func Intersect[T int | string](a []T, b []T) []T {
	var inter []T
	var mp MyMAP[T, bool] = map[T]bool{}

	for _, s := range a {
		if _, ok := mp[s]; !ok {
			mp[s] = true
		}
	}
	for _, s := range b {
		if _, ok := mp[s]; ok {
			inter = append(inter, s)
		}
	}

	return inter
}

// RemoveIntersection 删除交集
func RemoveIntersection[T int | string](arr1, arr2 []T) []T {
	arr2m := make(map[T]struct{})
	for _, i := range arr2 {
		arr2m[i] = struct{}{}
	}
	result := make([]T, 0)
	for _, i := range arr1 {
		if _, ok := arr2m[i]; !ok {
			result = append(result, i)
		}
	}
	return result
}

func MaxNum[T numElement](argus ...T) T {
	max := argus[0]
	for _, i := range argus {
		if i > max {
			max = i
		}
	}
	return max
}

func MinNum[T numElement](argus ...T) T {
	min := argus[0]
	for _, i := range argus {
		if i < min {
			min = i
		}
	}
	return min
}

func SliceString2Num[T numElement](src []string) ([]T, error) {
	dst := make([]T, 0, len(src))
	for i := range src {
		n, err := strconv.Atoi(src[i])
		if err != nil {
			return nil, err
		}
		dst = append(dst, T(n))
	}
	return dst, nil
}

func SortByFunc[T any](s []T, f func(i, j int) bool) []T {
	sort.SliceStable(s, f)
	return s
}

func FilterByFunc[T any](s []T, f func(T) bool) []T {
	m := make([]T, 0)
	for _, v := range s {
		if f(v) {
			m = append(m, v)
		}
	}
	return m
}

// StrToStrSlice "s1, s2, s3, s4" => []string{"s1", "s2", "s3", "s4"}
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

func SliceToStr[T numElement](arr []T) string {
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

// StrToNumSlice StrToSlice "1,2,3,4" => []int{1, 2, 3, 4}
// "[1,2,3,4]" => []int{1, 2, 3, 4}
func StrToNumSlice(str string) []int64 {
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

// All 切片中的所有元素是否都满足指定方法
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

// AnyOne 切片中是否存在满足指定方法的元素，如果不传方法，则判断切片中是否有元素，传多个方法只会取第一个
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

// Select 将切片中的每个元素都作为入参传入指定方法中，收集方法返回值并放入切片返回
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

// Count 返回切片中满足指定方法的元素个数
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

// Where 返回切片中满足指定方法的元素
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

// First 返回切片中第一个满足指定方法的元素，如不存在则返回nil
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

// FirstOrDefault 返回切片中第一个满足指定方法的元素，如不存在则返回入参中的默认值
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

// FilterEmptyStr 过滤字符串数组中的空字符串
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
