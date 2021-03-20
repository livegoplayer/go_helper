package utils

import "time"

//自动失败重试
//入参：
//t:重试次数
//sleepTime：失败后等待时间
//f:调用方法，interface{}为返回值，bool为是否成功
//callbacks:失败时的回调函数,i为调用次数
//出参：
//interface{}:调用方法的返回值
//bool:最终是否成功
func Retry(t int, sleepTime time.Duration, f func() (interface{}, bool), callbacks ...func(i int)) (interface{}, bool) {
	if t <= 0 {
		t = 3
	}
	if sleepTime < 0 {
		sleepTime = 500 * time.Millisecond
	}
	var res interface{}
	ok := false
	for i := 1; i <= t; i++ {
		res, ok = f()
		if ok {
			return res, true
		}
		for _, callBack := range callbacks {
			callBack(i)
		}
		if i != t {
			time.Sleep(sleepTime)
		}
	}
	res, ok = f()
	return res, ok
}
