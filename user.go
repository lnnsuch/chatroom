package chatroot

import (
	"github.com/gorilla/websocket"
	"fmt"
	"sync"
	"errors"
	"time"
	"encoding/json"
)

type (
	// 用户信息
	Account struct {
		id int
		user string
		name string
		password string
		cookie string
	}
	// 用户客户端信息
	Client struct {
		conn *websocket.Conn
		Account
	}
	// 所有用户
	Clients struct {
		client map[int]*Client
		sync.Mutex
	}
)

var Accounts = []Account{
	{id :1, user: "admin1",name: "admin1", password: "111111"},
	{id :2, user: "admin2",name: "admin2", password: "222222"},
	{id :3, user: "admin3",name: "admin3", password: "333333"},
	{id :4, user: "admin4",name: "admin4", password: "444444"},
	{id :5, user: "admin5",name: "admin5", password: "555555"},
}

var clients = Clients{client: make(map[int]*Client)}

// 获取所有客户端
func getAllClient() map[int]*Client {
	clients.Lock()
	defer clients.Unlock()
	return clients.client
}

// 获取客户端
func getClient (key int) *Client {
	clients.Lock()
	defer clients.Unlock()
	if v, ok := clients.client[key]; ok {
		return v
	} else {
		return nil
	}
}

// 添加客户端
func (c *Client) addUser() {
	clients.Lock()
	defer clients.Unlock()
	clients.client[c.id] = c
}

// 删除客户端
func (c *Client) delUser() {
	clients.Lock()
	defer clients.Unlock()
	delete(clients.client, c.id)
	// 如果该cookie10秒内未重新登录,则删除该cookie
	go func() {
		time.Sleep(time.Second * 10)
		if getClient(c.id) == nil {
			DelCookie(c.cookie)
		}
	}()
}

// 新建用户
func NewUser(account Account) *Client {
	go sendGroupMessage(account.name + " 加入", "", SysMessageStatusCode, GroupMessageTypeCode, getAllClient())
	client := &Client{Account: account}
	return client
}

// 获取用户信息
func getAccount(id int) (Account, error) {
	for _, val := range Accounts {
		if val.id == id {
			return val, nil
		}
	}
	return Account{}, errors.New("用户不存在")
}

func getAccountByName(name string) (Account, error) {
	for _, val := range Accounts {
		if val.name == name {
			return val, nil
		}
	}
	return Account{}, errors.New("用户不存在")
}

// 设置用户的cookie
func (a *Account) SetCookie(cookie string)  {
	a.cookie = cookie
}

// 读取用户发送的消息
func (c *Client) readMes() {
	for {
		mt, message, err := c.conn.ReadMessage()
		if err != nil {
			if mt == -1 {
				c.delUser()
				go sendGroupMessage(c.name + " 退出", "", SysMessageStatusCode, GroupMessageTypeCode, getAllClient())
			} else {
				fmt.Println("消息获取失败: ", err, "\n消息类型", mt)
			}
			break
		}
		go c.MessageHandle(message)
	}
}

// 消息处理
func (c *Client) MessageHandle(m []byte) {
	ms := NewNullMessage()
	err := json.Unmarshal(m, &ms)
	if err != nil {
		fmt.Println(err)
		message := NewMessage(-2, "消息发送失败")
		PutMessageQueue(c, message)
		return
	}
	switch ms.Type {
	case 1:
		c.GroupMessageHandle(ms)
	case 2:
		c.PrivateMessageHandle(ms)
	}
}

// 发送私聊信息
func (c *Client) sendPrivateMessageHandle (message Message, client *Client)  {
	message.Name = c.name
	message.Id = client.id
	PutMessageQueue(c, message)

	message.Id = c.id
	PutMessageQueue(client,message)
}

// 私聊信息处理
func (c *Client) PrivateMessageHandle (ws Message)  {
	if client := getClient(ws.Id); client == nil {
		fmt.Println("消息发送失败:对方已下线")
	} else {
		c.sendPrivateMessageHandle(ws, client)
	}
}

// 群聊信息处理
func (c *Client) GroupMessageHandle (m Message) {
	if m.Id == PublicGroupId {
		sendGroupMessage(m.Info, c.name, OkMessageStatusCode, GroupMessageTypeCode, getAllClient())
	}
}
