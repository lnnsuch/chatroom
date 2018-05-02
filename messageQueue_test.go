package chatroot

import (
	"testing"
)

func TestMessageQueue1(t *testing.T) {

}

func TestMessageQueue2(t *testing.T) {
	account := Account{}
	client := NewUser(account)
	message1 := NewMessage(1, "")
	message2 := NewMessage(2, "")
	message3 := NewMessage(3, "")
	message4 := NewMessage(0, "aaaaaa")
	PutMessageQueue(client, message1)
	PutMessageQueue(client, message2)
	PutMessageQueue(client, message3)
	PutMessageQueue(client, message4)
}
