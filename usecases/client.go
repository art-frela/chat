package usecases

import (
	"strconv"

	"github.com/art-frela/chat/domain"

	uuid "github.com/satori/go.uuid"
)

type ChatClient struct {
	domain.User
	chat      *ChatRoom
	rcv       chan domain.Message
	sessionID string
	counter   uint64
}

func NewChatClient(user domain.User, room *ChatRoom) *ChatClient {
	cc := &ChatClient{user, room, make(chan domain.Message, 1), uuid.NewV4().String(), 0}
	room.AddMember(user, cc.rcv)
	return cc
}

func (cc *ChatClient) SendMsg(txt string) {
	msg := domain.Message{
		ID:     cc.sessionID + "-" + strconv.FormatUint(cc.counter, 32),
		Author: cc.Nick,
		Body:   txt,
	}
	cc.counter += 1
	cc.chat.rcv <- msg
}

func (cc *ChatClient) RecieveMsg() (domain.Message, bool) {
	msg, more := <-cc.rcv
	return msg, more
}

func (cc *ChatClient) LeaveChat() {
	cc.chat.delMember(cc.User.ID)
}
