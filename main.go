package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/juggleim/jugglechat-server/jim"
)

func main() {
	server := &jim.JuggleChatServer{}
	server.Startup(map[string]interface{}{})

	closeChan := make(chan struct{})
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigChan
		server.Shutdown(true)
		signal.Stop(sigChan)
		close(closeChan)
	}()
	<-closeChan
}
