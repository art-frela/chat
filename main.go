package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/art-frela/chat/infra"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	quit := make(chan os.Signal)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	serv := infra.NewServer(ctx)
	serv.Run()
	<-quit
	cancel()
	serv.Stop()
}
