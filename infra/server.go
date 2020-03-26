package infra

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/art-frela/chat/domain"
	"github.com/art-frela/chat/usecases"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	mockChatRoomID    string = "1234567"
	mockChatRoomTitle string = "Room#1"
)

type Server struct {
	log      *zap.SugaredLogger
	config   *viper.Viper
	mux      *echo.Echo
	upgrader *websocket.Upgrader
	chatRoom *usecases.ChatRoom
}

// NewServer builder main server
func NewServer(ctx context.Context) *Server {
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	s := &Server{
		upgrader: upgrader,
	}
	// fill config and logger
	s.setConfig()
	s.setLogger()
	//
	s.chatRoom = usecases.NewChatRoom(ctx, mockChatRoomID, mockChatRoomTitle, s.log)

	return s
}

// Run is running Server
func (s *Server) Run() {
	s.registerRoutes()
	port := s.config.GetString("httpd.port")
	host := s.config.GetString("httpd.host") + ":" + port
	s.log.Infof("http server starting on the [%s] tcp port", host)
	go func() {
		if err := s.mux.Start(host); err != http.ErrServerClosed {
			s.log.Fatalf("http server error: %v", err)
		}
	}()

}

// Stop is stopping Server
func (s *Server) Stop() {
	s.log.Infof("got signal to stopping server")
	stopDuration := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), stopDuration)
	defer cancel()
	if err := s.mux.Shutdown(ctx); err != nil {
		s.log.Fatal(err)
	}
}

func (s *Server) registerRoutes() {
	e := echo.New()
	e.HidePort = true
	e.HideBanner = true // hide banner ECHO
	e.Use(middleware.Recover())
	// metric handler
	e.Static("/", "assets")
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	e.GET("/ws", s.serveWS)

	s.mux = e
}

func (s *Server) serveWS(c echo.Context) error {
	// use middleware to extract token from requets and valid them and exract user
	//
	s.log.Debugf("connected new client, remoteAddr: %s, agent: %s", c.Request().RemoteAddr, c.Request().UserAgent())
	ws, err := s.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	// TODO: replace to user repo extractor
	user := domain.User{
		ID:   uuid.NewV4().String(),
		Nick: uuid.NewV4().String()[:7],
	}

	client := usecases.NewChatClient(user, s.chatRoom)
	defer client.LeaveChat()

	wg := &sync.WaitGroup{}
	// Write to ws
	wg.Add(1)
	go s.writeToWS(client, ws, wg)

	// Read from ws
	wg.Add(1)
	go s.readFromWS(client, ws, wg)

	wg.Wait()
	return nil
}

func (s *Server) writeToWS(client *usecases.ChatClient, ws *websocket.Conn, wg *sync.WaitGroup) {
	log := s.log.With(zap.String("clientID", client.ID))
	log.Debug("start writer to ws")
	defer wg.Done()
	defer log.Debug("stop writer to ws")
	for {
		msg, more := client.RecieveMsg()
		if !more {
			log.Warn("closed client chan, aborted...")
			return
		}
		err := ws.WriteJSON(msg)
		if err != nil {
			if err == websocket.ErrCloseSent {
				log.Errorf("ws error, %v, aborted...", err)
				return
			}
			log.Errorf("ws error, %v", err)
		}
	}
}

func (s *Server) readFromWS(client *usecases.ChatClient, ws *websocket.Conn, wg *sync.WaitGroup) {
	log := s.log.With(zap.String("clientID", client.ID))
	log.Debug("start reader from ws")
	defer wg.Done()
	defer log.Debug("stop reader from ws")
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Errorf("read message from ws error, %v, aborted...", err)
			return
		}
		client.SendMsg(string(msg))
	}
}
