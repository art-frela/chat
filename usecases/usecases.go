/*Package usecases contains businesslogic of application

Author: Karpov Artem, mailto: art.frela@gmail.com
Date: 2020-03-25
*/
package usecases

import (
	"context"
	"sync"

	"github.com/art-frela/chat/domain"
	"go.uber.org/zap"
)

const (
	msgStopChat string = "got signal to stop chat, sorry..."
)

// ChatRoom..
type ChatRoom struct {
	ID       string
	Title    string
	rcv      chan domain.Message
	chCancel <-chan struct{}
	log      *zap.SugaredLogger
	mu       *sync.Mutex
	members  map[string]chan domain.Message
}

func NewChatRoom(ctx context.Context, id, title string, logger *zap.SugaredLogger) *ChatRoom {
	cr := &ChatRoom{
		ID:       id,
		Title:    title,
		rcv:      make(chan domain.Message, 1),
		chCancel: ctx.Done(),
		log:      logger,
		mu:       &sync.Mutex{},
		members:  make(map[string]chan domain.Message),
	}
	// run canceler
	go cr.canceler()

	// run chatworker
	go cr.chatProcessor()
	return cr
}

func (cr *ChatRoom) chatProcessor() {
	cr.log.Debug("start chatProcessor")
	defer cr.log.Debug("stopped chatProcessor")
	for {
		select {
		case <-cr.chCancel:
			return
		case msg, more := <-cr.rcv:
			if !more {
				return
			}
			cr.reSendMsg(msg)
		}
	}
}

func (cr *ChatRoom) canceler() {
	<-cr.chCancel
	cr.sayGoodBuy(msgStopChat)
}

func (cr *ChatRoom) sayGoodBuy(txt string) {
	msg := domain.Message{
		ID:     "-1",
		Author: "system",
		Body:   txt,
	}
	cr.reSendMsg(msg)
	cr.clearMembers()
}

// AddMember adds member of chatroom to
func (cr *ChatRoom) AddMember(user domain.User, ch chan domain.Message) chan domain.Message {
	cr.mu.Lock()
	cr.members[user.ID] = ch
	cr.mu.Unlock()
	return cr.rcv
}

func (cr *ChatRoom) delMember(key string) {
	cr.mu.Lock()
	delete(cr.members, key)
	cr.mu.Unlock()
}

func (cr *ChatRoom) clearMembers() {
	cr.mu.Lock()
	cr.members = make(map[string]chan domain.Message)
	cr.mu.Unlock()
}

func (cr *ChatRoom) reSendMsg(msg domain.Message) {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	for _, out := range cr.members {
		out <- msg
	}
}
