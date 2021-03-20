package utils

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Resp 返回
type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type EmptyData struct {
}

// ErrorResp 错误返回值
func ErrorResp(ctx *gin.Context, code int, msg string, data ...interface{}) {
	resp(ctx, code, msg, data...)
}

// SuccessResp 正确返回值
func SuccessResp(ctx *gin.Context, msg string, data ...interface{}) {
	resp(ctx, 0, msg, data...)
}

// resp 返回
func resp(ctx *gin.Context, code int, msg string, data ...interface{}) {
	resp := Resp{
		Code: code,
		Msg:  msg,
		Data: data,
	}

	if len(data) == 1 {
		resp.Data = data[0]
	}

	if len(data) == 0 {
		resp.Data = &EmptyData{}
	}

	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Origin", "*")                                       // 这是允许访问所有域
	ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
	//  header的类型
	ctx.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
	//              允许跨域设置                                                                                                      可以返回其他子段
	ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
	ctx.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
	ctx.Header("Access-Control-Allow-Credentials", "false")                                                                                                                                                  //  跨域请求是否需要带cookie信息 默认设置为true
	ctx.Set("content-type", "application/json")                                                                                                                                                              // 设置返回格式是json

	ctx.JSON(http.StatusOK, resp)
}

//负责调用panic触发外部的panic处理函数
func CheckError(error error, message ...string) {
	var msg string

	if len(message) == 0 {
		msg = ""
	}

	if error != nil {
		msg = strings.Join(message, " ")
		msg = msg + " error:" + error.Error()
		error = NewError(200, 1, msg)
		panic(error)
	}
}
