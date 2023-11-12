package main

import (
	proxy_server "HTTP_Server/proxy-server"
	web_server "HTTP_Server/web-server"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// start server and proxy in separate goroutines
	go web_server.StartWebServer()
	go proxy_server.StartProxyServer()

	// we create a channel to wait for an interrupt signal
	quit := make(chan os.Signal, 1)

	// notify the quit channel on SIGINT (CTRL+C on windows) or SIGTERM (termination signal on operative system level)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// wait until quit channel has received signal SIGINT or SIGTERM
	<-quit
}
