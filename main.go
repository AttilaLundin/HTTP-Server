package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// start server and proxy in separate goroutines

	//TODO: uncomment the line below before building the web server
	//go web_server.StartWebServer()

	//TODO: uncomment the line below before building the proxy server
	//go proxy_server.StartProxyServer()

	//TODO: uncomment the line blow and run main if you want to run test on the web server
	//make sure to add the ip address of the webserver and the relevant port
	//address := "<IP>:<PORT>"
	//go test.INITTEST(address)

	// we create a channel to wait for an interrupt signal
	quit := make(chan os.Signal, 1)

	// notify the quit channel on SIGINT (CTRL+C on windows) or SIGTERM (termination signal on operative system level)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// wait until quit channel has received signal SIGINT or SIGTERM
	<-quit
}
