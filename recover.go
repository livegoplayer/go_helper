package helper

//处理全局panic的返回值，重写gin.Recover中间件的内容
import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	mylogger "github.com/livegoplayer/go_logger"
	"github.com/spf13/viper"
)

// 错误处理的结构体
type Error struct {
	//只是说他们两个差不多，没用到
	Resp       `json:"-"`
	StatusCode int         `json:"-"`
	Code       int         `json:"code"`
	Msg        string      `json:"msg"`
	Data       interface{} `json:"data"`
}

var (
	ServerError = NewError(http.StatusInternalServerError, 1, "系统异常，请稍后重试!")
	NotFound    = NewError(http.StatusNotFound, 1, http.StatusText(http.StatusNotFound))
)

func OtherError(message string) *Error {
	return NewError(http.StatusForbidden, 100403, message)
}

func (e *Error) Error() string {
	return e.Msg
}

func NewError(statusCode, Code int, msg string) *Error {
	return &Error{
		StatusCode: statusCode,
		Code:       Code,
		Msg:        msg,
		Data:       &EmptyData{},
	}
}

// 404处理
func HandleNotFound(c *gin.Context) {
	err := NotFound
	c.JSON(err.StatusCode, err)
	c.Abort()
	return
}

// 服务异常处理
func HandleServerError(c *gin.Context) {
	err := ServerError
	c.JSON(err.StatusCode, err)
	c.Abort()
	return
}

func ErrHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var Err *Error
				//如果是通过本文件定义的Error，如果是调试模式，则输出所有的错误内容，否则，只输出自定义内容
				if e, ok := err.(*Error); ok {
					Err = e
					if !gin.IsDebugging() {
						msg := GetSubStringBetween(Err.Msg, " error:", "")
						Err.Msg = msg
					}
				} else if e, ok := err.(error); ok {
					if !gin.IsDebugging() {
						Err = ServerError
					} else {
						Err = OtherError(e.Error())
					}
					//这种程度的error, 输出到数据库
					mylogger.Error(e.Error())
				} else {
					Err = ServerError
					mylogger.Error(Err.Error())
				}
				// 记录一个错误的日志
				c.JSON(Err.StatusCode, Err)
				c.Abort()
				return
			}
		}()
		c.Next()
	}
}

//
////// 跨域
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               //请求方法
		origin := c.Request.Header.Get("Origin") //请求头部
		var headerKeys []string                  // 声明请求头keys
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}

		//获取配置文件中的host
		var accessControlAllowOrigin string
		host := viper.GetString("host")
		port := viper.GetString("port")
		if host == "" || port == "" {
			accessControlAllowOrigin = "*"
		}
		accessControlAllowOrigin = "http://" + host + port
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", accessControlAllowOrigin)                  // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//              允许跨域设置                                                                                                      可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
			//允许设置cookie
			c.Header("Access-Control-Allow-Credentials", "true") //  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")            // 设置返回格式是json
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next() //  处理请求
	}
}
