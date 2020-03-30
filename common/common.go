package common

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
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

// 结构体转map1 保留 Model
func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	data := make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

// 结构体转map2
func StructToMapViaJson(data interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	//struct 转json
	j, _ := json.Marshal(data)
	// log.Fatalln(j)
	//json 转map
	json.Unmarshal(j, &m)
	return m
}

// JSONTime format json time field by myself
type JSONTime struct {
	time.Time
}

// MarshalJSON on JSONTime format Time field with %Y-%m-%d %H:%M:%S
func (t JSONTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

// Value insert timestamp into mysql need this function.
func (t JSONTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan valueof time.Time
func (t *JSONTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JSONTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
