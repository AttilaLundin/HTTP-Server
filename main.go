package main

import (
	"HTTP_Server/test"
	web_server "HTTP_Server/web-server"
)

func main() {
	go web_server.StartWebServer()
	//go proxy_server.StartProxyServer()
	go test.INITTEST()

	for {
	}
}
