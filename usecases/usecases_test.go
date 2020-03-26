package usecases

import (
	"context"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/art-frela/chat/domain"
	"go.uber.org/zap"
)

var (
	log     *zap.SugaredLogger = zap.NewNop().Sugar()
	timeout time.Duration      = 2 * time.Second
)

func TestAddMember(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cr := NewChatRoom(ctx, "001", "testChat", log)
	var client1 *ChatClient = NewChatClient(domain.User{ID: "1", Nick: "client1"}, cr)
	//client2 *ChatClient = NewChatClient(domain.User{ID: "2", Nick: "client2"})
	//client3 *ChatClient = NewChatClient(domain.User{ID: "3", Nick: "client3"})

	cr.AddMember(client1.User, client1.rcv)
	if _, ok := cr.members[client1.ID]; !ok {
		t.Error("expected exists member of chatRoom, but not")
	}
}

func TestDelMember(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cr := NewChatRoom(ctx, "001", "testChat", log)
	var client3 *ChatClient = NewChatClient(domain.User{ID: "3", Nick: "client3"}, cr)
	cr.AddMember(client3.User, client3.rcv)
	if _, ok := cr.members[client3.ID]; !ok {
		t.Error("expected exists member of chatRoom, but not")
	}
	cr.delMember(client3.ID)
	if _, ok := cr.members[client3.ID]; ok {
		t.Error("expected not exists member of chatRoom, but exists")
	}
}

// этим тестом мы проверяем что останавливаются все горутины которые у вас были и нет утечек
// некоторый запас ( goroutinesPerTwoIterations*5 ) остаётся на случай рантайм горутин
func TestChatRoomLeak(t *testing.T) {
	goroutinesStart := runtime.NumGoroutine()
	TestDelMember(t)
	goroutinesPerTwoIterations := runtime.NumGoroutine() - goroutinesStart

	goroutinesStart = runtime.NumGoroutine()
	goroutinesStat := []int{}
	for i := 0; i <= 25; i++ {
		TestDelMember(t)
		goroutinesStat = append(goroutinesStat, runtime.NumGoroutine())
	}
	goroutinesPerFiftyIterations := runtime.NumGoroutine() - goroutinesStart
	if goroutinesPerFiftyIterations > goroutinesPerTwoIterations*5 {
		t.Fatalf("looks like you have goroutines leak: %+v", goroutinesStat)
	}
}

func TestReSendMsg(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cr := NewChatRoom(ctx, "001", "testChat", log)
	const (
		testMsg string = "testMessage"
	)
	var client1 *ChatClient = NewChatClient(domain.User{ID: "1", Nick: "client1"}, cr)
	var client2 *ChatClient = NewChatClient(domain.User{ID: "2", Nick: "client2"}, cr)
	client1Msg := domain.Message{
		ID:     client1.sessionID + "-0",
		Author: client1.Nick,
		Body:   testMsg,
	}
	client2Msg := domain.Message{
		ID:     client2.sessionID + "-0",
		Author: client2.Nick,
		Body:   testMsg,
	}

	// send client1
	client1.SendMsg(testMsg)
	reply1, _ := client1.RecieveMsg()
	reply2, _ := client2.RecieveMsg()
	if !reflect.DeepEqual(reply1, client1Msg) {
		t.Errorf("expected got message %+v, but got %+v", client1Msg, reply1)
	}
	if !reflect.DeepEqual(reply2, client1Msg) {
		t.Errorf("expected got message %+v, but got %+v", client1Msg, reply2)
	}
	client2.SendMsg(testMsg)
	reply1, _ = client1.RecieveMsg()
	reply2, _ = client2.RecieveMsg()
	if !reflect.DeepEqual(reply1, client2Msg) {
		t.Errorf("expected got message %+v, but got %+v", client2Msg, reply1)
	}
	if !reflect.DeepEqual(reply2, client2Msg) {
		t.Errorf("expected got message %+v, but got %+v", client2Msg, reply2)
	}
}

func TestSayGoodBuy(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	cr := NewChatRoom(ctx, "001", "testChat", log)
	var client3 *ChatClient = NewChatClient(domain.User{ID: "3", Nick: "client3"}, cr)
	cr.AddMember(client3.User, client3.rcv)
	if _, ok := cr.members[client3.ID]; !ok {
		cancel()
		t.Error("expected exists member of chatRoom, but not")
	}
	cancel() // close context.Done()
	goodBuy, _ := client3.RecieveMsg()
	if goodBuy.Body != msgStopChat {
		t.Errorf("expeced goodBuy message=%s, but got %s", msgStopChat, goodBuy.Body)
	}
}
