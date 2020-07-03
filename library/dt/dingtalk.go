package dt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var uri string = "https://oapi.dingtalk.com/robot/send?access_token=7cb1aab23634d2ebfdd472b37c2aaca6af658951aae07dc3832e44cc1b5a15b0dd"

var secret string = "SECf3146c59bde88be974ce8b667414f8776074849e89c4060a11566c675c70cf6fdd"

// Text 消息内容
type Text struct {
	Content string `json:"content"`
}

var atMobiles []string

// AT @手机号或者所有人
type AT struct {
	ATMobiles []string `json:"atMobiles,omitempty"`
	ISATAll   bool     `json:"isAtAll,omitempty"`
}

// TextParams 文本类型请求
type TextParams struct {
	MsgType string `json:"msgtype"`
	Text    Text   `json:"text"`
	AT      AT     `json:"at"`
}

// getURL 获取加密请求URL
func getURL() string {
	timestamp := getSign()
	stringToSign := fmt.Sprintf("%v\n%v", timestamp, secret)
	fmt.Printf("timestamp : %#v  , stringToSign : %#v\n", timestamp, stringToSign)
	signBase64 := ComputeHmac256(stringToSign, secret)
	fmt.Println(signBase64)
	sign := url.QueryEscape(signBase64)
	fmt.Println(sign)
	uriNew := fmt.Sprintf("%v&timestamp=%v&sign=%s", uri, timestamp, sign)
	fmt.Println(uriNew)
	return uriNew
}

// Post 请求
func Post(s string, isAtAll bool) {
	contentType := "application/json"
	params := TextParams{
		MsgType: "text",
		Text: Text{
			Content: s,
		},
		AT: AT{
			ISATAll: isAtAll,
		},
	}
	data, err := json.Marshal(params)
	if err != nil {
		fmt.Printf("Post data Marshal err : %#v \n", err)
		return
	}
	// fmt.Printf("%v \n", string(data))
	// return
	resp, err := http.Post(getURL(), contentType, strings.NewReader(string(data)))
	if err != nil {
		fmt.Printf("post failed,err:%v \n", err)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("get resp failed, err:%v \n", err)
		return
	}

	fmt.Println(string(b))

	return
}

// getSign 获取当前的时间戳
func getSign() int64 {
	MTime := time.Now().UnixNano() / 1e6
	return MTime
}

// ComputeHmac256 HmacSha256加密
func ComputeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
