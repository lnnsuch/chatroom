package chatroot

import (
	"encoding/json"
	"net/http"
	"fmt"
)

var HttpResultArr = map[int]interface{}{
	-2: "系统通知",
	-1: []string{"您未登录,请先登录", "/public/login.html"},
	0: "成功",
	1: "用户名密码错误",
	2: "链接错误",
}

type (
	Message struct {
		Status int `json:"status"` // HttpResultArr
		Info string `json:"info"` // 信息
		Url string `json:"url,omitempty"`
		Name string `json:"name,omitempty"` // 用户名
		Type int `json:"type,omitempty"` // 信息类型(1：群聊 2：私聊 3：群聊总数)
		Id int `json:"id,omitempty"` // id
	}
)

// 返回ajax信息 w：http返回包 status：信息KEY
func AjaxReturn(w http.ResponseWriter, status int){
	info := NewMessage(status, "")
	str, err := json.Marshal(info)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(str)
}

// 返回空消息
func NewNullMessage() Message {
	return Message{}
}

// 新建消息 status：消息KEY info：发送的信息(当信息为空时,从消息数组中获取信息)
func NewMessage(status int, info string) Message {
	message := Message{Status: status}
	if info != "" {
		message.Info = info
		return message
	}
	if v, ok := HttpResultArr[status]; ok {
		switch v.(type) {
		case string:
			message.Info = v.(string)
		case []string:
			arr := v.([]string)
			message.Info = arr[0]
			message.Url = arr[1]
		default:
			fmt.Println("status", status, "值不符合输出格式,请检查")
		}
	} else {
		message.Info = "未知错误"
	}
	return message
}

// 发送系统公告 info: 系统通知
func sendSysMessage(info string) {
	m := Message{Status: -2, Info: info}
	mes, _ := json.Marshal(m)
	for _, client := range getAllClient() {
		err := client.conn.WriteMessage(1, mes)
		if err != nil {
			fmt.Println("---发送系统消息错误---", err)
		}
	}
}

// 发送文本信息 info：发送的信息 name：发送人的姓名
func sendTextMessage(info, name string) {
	m := NewMessage(0, info)
	m.Name = name
	mes, _ := json.Marshal(m)
	for _, client := range getAllClient() {
		err := client.conn.WriteMessage(1, mes)
		if err != nil {
			fmt.Println("---发送文本消息错误---", err)
		}
	}
}

