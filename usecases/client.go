package usecases

import (
	"strconv"

	"github.com/art-frela/chat/domain"

	uuid "github.com/satori/go.uuid"
)

type ChatClient struct {
	domain.User
	rcv       chan domain.Message
	snd       chan domain.Message
	sessionID string
	counter   uint64
}

func NewChatClient(user domain.User) *ChatClient {
	return &ChatClient{
		user, make(chan domain.Message, 1), nil, uuid.NewV4().String(), 0,
	}
}

func (cc *ChatClient) AddSendChan(out chan domain.Message) {
	cc.snd = out
}

func (cc *ChatClient) SendMsg(txt string) {
	msg := domain.Message{
		ID:     cc.sessionID + "-" + strconv.FormatUint(cc.counter, 32),
		Author: cc.Nick,
		Body:   txt,
	}
	cc.counter += 1
	cc.snd <- msg
}

func (cc *ChatClient) RecieveMsg() domain.Message {
	return <-cc.rcv
}
