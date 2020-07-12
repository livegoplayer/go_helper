package helper

//随机验证码图片，以后可以拓展成随机短信验证码 todo 以后可以抽成微服务

import (
	"image/color"
	"strings"
	"time"

	"github.com/go-redis/redis"

	captcha "github.com/mojocn/base64Captcha"
)

var captchaMaker *captcha.Captcha

func init() {

	//设置回值二维码的size
	driver := captcha.NewDriverString(60, 240, 10, captcha.OptionShowHollowLine, 6, "123123", &color.RGBA{R: 255, G: 255, B: 255, A: 255}, nil).ConvertFonts()
	redisStore := NewRedisStore(time.Hour * 24)
	captchaMaker = &captcha.Captcha{
		Driver: driver,
		Store:  redisStore,
	}
}

func MakeCaptcha(s string) (captchaId string, captchaImg string, err error) {
	driver := captchaMaker.Driver.(*captcha.DriverString)
	if driver == nil {
		panic("二维码替换字符串失败")
	}

	driver.Source = s
	captchaMaker.Driver = driver

	captchaId, captchaImg, err = captchaMaker.Generate()
	if err != nil {
		panic("二维码生成失败")
	}

	return
}

func VerifyCaptchaWithId(captchaId string, answer string) bool {
	return captchaMaker.Store.Verify(captchaId, answer, true)
}

// memoryStore is an internal store for captcha ids and their values.
type redisStore struct {
	expiration  time.Duration
	redisClient *redis.Client
}

var prefix = "go_us_redis_"

func NewRedisStore(expiration time.Duration) captcha.Store {
	// 根据redis配置初始化一个客户端
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "139.224.132.234:6379", // redis地址
		Password: "myredis",              // redis密码，没有则留空
		DB:       1,                      // 默认数据库，默认是0
	})

	redisStore := redisStore{}
	redisStore.expiration = expiration
	redisStore.redisClient = redisClient
	return &redisStore
}

func (s *redisStore) Set(id string, value string) {
	s.redisClient.Set(prefix+id, strings.ToUpper(value), s.expiration)
}

func (s *redisStore) Verify(id, answer string, clear bool) bool {
	v := s.Get(id, clear)
	return v == strings.ToUpper(answer)
}

func (s *redisStore) Get(id string, clear bool) (value string) {
	v := s.redisClient.Get(prefix + id)
	value = strings.ToUpper(v.Val())

	if clear {
		s.redisClient.Del(prefix + id)
	}

	return
}
