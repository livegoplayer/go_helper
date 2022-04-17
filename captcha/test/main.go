package main

import (
	"github.com/go-redis/redis"
	"github.com/livegoplayer/go_helper/captcha"
)

func main() {
	captcha.InitCaptcha(redis.Options{Addr: "127.0.0.1:6379", Password: "myredis", DB: 1})
	id, _, _ := captcha.MakeCaptcha("123213")
	res := captcha.VerifyCaptchaWithId(id, "")
	print(res)
}
