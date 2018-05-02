package chatroot

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
	"fmt"
)

var HttpResultArr = map[int]interface{}{
	-2: "系统通知",
	-1: []string{"您未登录,请先登录", "/public/login.html"},
	0: "成功",
	1: "用户名密码错误",
	2: "链接错误",
	3: "用户不存在",
}

const (
	PublicGroupId = 1

	SysMessageStatusCode = -2
	NotLoginMessageStatusCode = -1
	OkMessageStatusCode = 0

	GroupMessageTypeCode = 1
	PrivateMessageTypeCode = 2
)

type (
	Message struct {
		Status int `json:"status"` // HttpResultArr
		Info string `json:"info"` // 信息
		Url string `json:"url,omitempty"`
		Name string `json:"name"` // 用户名
		Type int `json:"type"` // 信息类型(1：群聊 2：私聊 3：群聊总数)
		Id int `json:"id"` // id
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

// 发送系统公告 info: 系统通知 status：状态位 t：消息类型 client：需要发送的用户
func sendGroupMessage(info, name string, status int, t int, client map[int]*Client) {
	m := NewNullMessage()
	m.Status = status
	m.Type = t
	m.Name = name
	m.Id = PublicGroupId
	m.Info = info
	for _, client := range client {
		PutMessageQueue(client, m)
	}
}

func sendMessage(message Message, conn *websocket.Conn) {
	err := conn.WriteJSON(message)
	if err != nil {
		fmt.Println("--消息发送错误--", err)
	}
}

