package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/keepondream/article/common"
	"github.com/keepondream/article/library/dt"
)

// List api_list
type List struct {
	APIURI         string `json:"api_uri"`
	APIName        string `json:"api_name"`
	APIProtocol    string `json:"api_protocol"`
	APIRequestType string `json:"api_request_type"`
	APIStatus      string `json:"api_status"`
}

// Content 内容
type Content struct {
	Operation      string `json:"operation"`
	Operator       string `json:"operator"`
	ProjectName    string `json:"project_name"`
	ProjectID      string `json:"project_id"`
	APIList        []List `json:"api_list"`
	APIName        string `json:"api_name,omitempty"`
	APIURI         string `json:"api_uri,omitempty"`
	APIProtocol    string `json:"api_protocol,omitempty"`
	APIRequestType string `json:"api_request_type,omitempty"`
	APIStatus      string `json:"api_status,omitempty"`
}

// Data 返回数据集
type Data struct {
	HookRequestTimestamp int    `json:"hook_request_timestamp"`
	HookProduct          string `json:"hook_product"`
	SpaceID              string `json:"space_id"`
	SpaceName            string `json:"space_name"`
	HookOperator         string `json:"hook_operator"`
	HookEvent            string `json:"hook_event"`
	HookRequestTime      string `json:"hook_request_time"`
	HookOperation        string `json:"hook_operation"`
	Content              `json:"content"`
}

// Eolinker eolinker 接口变动消息
func Eolinker(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	if body == nil {
		common.Failed(c, common.WithMsg("请求body不能为空"))
		return
	}
	var data Data

	err := json.Unmarshal(body, &data)
	if err != nil {
		fmt.Printf("json.Unmarshal failed, err:%v\n", err)
		common.Success(c)
		return
	}

	if data.ProjectName != "易报表" {
		common.Success(c)
		return
	}

	var OutMsg string

	switch data.HookEvent {
	case "api_document":
		// API文档事件
		fmt.Println("API文档事件")
		OutMsg = getAPIDocumentMsg(data)
	case "api_management_project":
		// API 管理项目事件
		fmt.Println("API管理项目事件")
		OutMsg = getManagementProjectMsg(data)
	}
	dt.Post(OutMsg, false)

	fmt.Println(OutMsg)
	//dt.Post(OutMsg, true)
	common.Success(c)
	return
}

// getAPIDocumentMsg 获取API文档事件消息体
func getAPIDocumentMsg(d Data) (OutMsg string) {
	OutMsg += `【` + d.ProjectName + `】` + d.HookOperator
	switch strings.ToLower(d.HookOperation) {
	case "add_api_document":
		OutMsg += `，添加了API，接口名：` + d.APIName + `，
		协议:` + d.APIProtocol + `，请求方式:` + d.APIRequestType + `，
		请求路由:` + d.APIURI + `，
		状态:` + getStatus(d.APIStatus)

	case "update_api_document":
		OutMsg += `，修改了API，接口名：` + d.APIName + `，
		协议:` + d.APIProtocol + `，请求方式:` + d.APIRequestType + `，
		请求路由:` + d.APIURI + `，
		状态:` + getStatus(d.APIStatus)
	case "delete_api_document":
		OutMsg += `，删除了API`
		if len(d.Content.APIList) > 0 {
			OutMsg += `，接口名：` + d.Content.APIList[0].APIName + `，
			协议:` + d.Content.APIList[0].APIProtocol + `，请求方式:` + d.Content.APIList[0].APIRequestType + `，
			请求路由:` + d.Content.APIList[0].APIURI + `，
			状态:` + getStatus(d.Content.APIList[0].APIStatus)
		}

	case "switch_api_status":
		OutMsg += `，切换了API状态`
		if len(d.Content.APIList) > 0 {
			OutMsg += `，接口名：` + d.Content.APIList[0].APIName + `，
			协议:` + d.Content.APIList[0].APIProtocol + `，请求方式:` + d.Content.APIList[0].APIRequestType + `，
			请求路由:` + d.Content.APIList[0].APIURI + `，
			状态:` + getStatus(d.Content.APIList[0].APIStatus)
		}
	case "switch_api_document_version":
		OutMsg += `，切换了API文档版本，接口名：` + d.APIName + `，
		协议:` + d.APIProtocol + `，请求方式:` + d.APIRequestType + `，
		请求路由:` + d.APIURI + `，
		状态:` + getStatus(d.APIStatus)
	case "add_api_document_comment":
		OutMsg += `，评论了API，接口名：` + d.APIName + `，
		协议:` + d.APIProtocol + `，请求方式:` + d.APIRequestType + `，
		请求路由:` + d.APIURI + `，
		状态:` + getStatus(d.APIStatus)
	}

	fmt.Println(OutMsg)

	// dd(d)
	return
}

// getManagementProjectMsg API 管理项目事件消息体
func getManagementProjectMsg(d Data) (OutMsg string) {

	return
}

// getStatus eolinker api_status 数据清洗
func getStatus(s string) (status string) {
	s = strings.ToLower(s)
	status = `
╔───────╗
│     *^_^* 已发布     │　
│　    😀😀😀  　  │ 　
╚──────㊣╝
	`
	switch s {
	case "enable":
		status = `
	╔───────╗
	│     *^_^* 已发布     │　
	│　    😀😀😀  　  │ 　
	╚──────㊣╝
		`
	case "plan":
		status = `
	╔───────╗
	│     *^_^* 设计中     │　
	│　    😀😀😀  　  │ 　
	╚──────㊣╝
		`
	case "pending":
		status = `
	╔───────╗
	│     *^_^* 待确定     │　
	│　    😮😮😮  　  │ 　
	╚──────㊣╝
		`

	case "dev":
		status = `
	╔───────╗
	│     *^_^* 开发中     │　
	│　    💪💪💪  　  │ 　
	╚──────㊣╝
		`
	case "debug":
		status = `
	╔───────╗
	│     *^_^* 对接中     │　
	│　    😤😤😤  　  │ 　
	╚──────㊣╝
		`
	case "test":
		status = `
	╔───────╗
	│     *^_^* 测试中     │　
	│　    😌😌😌  　  │ 　
	╚──────㊣╝
		`
	case "compelete":
		status = `
	╔───────╗
	│     *^_^* 完 成     │　
	│　    😁😁😁  　  │ 　
	╚──────㊣╝
		`
	case "bug":
		status = `
	╔───────╗
	│     *^_^* 异 常     │　
	│　    😱😱😱  　  │ 　
	╚──────㊣╝
		`
	case "maintain":
		status = `
	╔───────╗
	│     *^_^* 维 护     │　
	│　    😫😫😫  　  │ 　
	╚──────㊣╝
		`
	case "ban":
		status = `
	╔───────╗
	│     *^_^* 废 弃     │　
	│　    👿👿👿  　  │ 　
	╚──────㊣╝
		`
	}
	return
}

func dd(i interface{}) {
	kv := make(map[string]interface{})
	vValue := reflect.ValueOf(i)
	vType := reflect.TypeOf(i)

	for i := 0; i < vValue.NumField(); i++ {
		kv[vType.Field(i).Name] = vValue.Field(i)
	}

	fmt.Println("打印数据:")
	for k, v := range kv {
		fmt.Print(k)
		fmt.Print(":")
		fmt.Print(v)
		fmt.Println()
	}
}
