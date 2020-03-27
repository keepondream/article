package common

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 定义自定义的参数集合
type Option struct {
	Code int                    `json:"code"`
	Msg  string                 `json:"msg"`
	Time int64                  `json:"time"`
	Data map[string]interface{} `json:"data"`
}

// 定义修改默认参数的钩子函数
type ModOption func(option *Option)

// 实际修改data默认参数的函数
func WithData(data map[string]interface{}) ModOption {
	return func(option *Option) {
		option.Data = data
	}
}

// 实际修改msg默认参数的函数
func WithMsg(msg string) ModOption {
	return func(option *Option) {
		option.Msg = msg
	}
}

// 实际修改code默认参数的函数
func WithCode(code int) ModOption {
	return func(option *Option) {
		option.Code = code
	}
}

// 成功的响应json数据体
func Success(c *gin.Context, modOptions ...ModOption) {
	code := 200
	msg := "请求成功!"
	time := time.Now().Unix()
	data := make(map[string]interface{})
	option := Option{
		Code: code,
		Msg:  msg,
		Time: time,
		Data: data,
	}

	for _, fn := range modOptions {
		fn(&option)
	}

	c.JSON(http.StatusOK, option)
}

// 错误响应json数据体
func Failed(c *gin.Context, modOptions ...ModOption) {
	code := 400
	msg := "请求失败!~"
	time := time.Now().Unix()
	data := make(map[string]interface{})
	option := Option{
		Code: code,
		Msg:  msg,
		Time: time,
		Data: data,
	}
	for _, fn := range modOptions {
		fn(&option)
	}
	c.JSON(http.StatusOK, option)
}
