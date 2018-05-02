package chatroot

import (
	"sync"
	"errors"
	"fmt"
)

type (
	messageQueue struct {
		sync.Mutex
		queue []*queue
	}
	queue struct {
		client *Client
		message Message
	}
)

var list *messageQueue

func init() {
	list = NewMessageQueue()
	go list.sendMessage()
}

// 新建消息队列
func NewMessageQueue() *messageQueue {
	return &messageQueue{}
}

// 循环发送消息队列的消息
func (q *messageQueue) sendMessage() {
	for {
		if q.length() > 0 {
			if l, err := q.get(); err != nil {
				fmt.Println(err)
			} else {
				sendMessage(l.message, l.client.conn)
				//fmt.Printf("%+v\n", l)
			}
		}
	}
}

// 写入
func (q *messageQueue) put(client *Client, message Message) {
	q.Lock()
	defer q.Unlock()
	l := &queue{client, message}
	q.queue = append(q.queue, l)
}

// 获取第一条
func (q *messageQueue) get() (*queue, error) {
	q.Lock()
	defer q.Unlock()
	if q.queue == nil {
		return nil, errors.New("queue null")
	}
	l := q.queue[0]
	q.queue = q.queue[1:]
	return l, nil
}

// 记录数
func (q *messageQueue) length() int {
	q.Lock()
	defer q.Unlock()
	return len(q.queue)
}

// 获取消息的第一条记录
func GetMessageQueue() (*queue, error) {
	return list.get()
}

// 写入消息队列
func PutMessageQueue(client *Client, message Message) {
	list.put(client, message)
}

// 获取消息队列的长度
func lenMessageQueue() int {
	return list.length()
}
