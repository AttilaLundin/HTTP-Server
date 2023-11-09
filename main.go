package main

import (
	proxy_server "HTTP_Server/proxy-server"
	"HTTP_Server/test"
	web_server "HTTP_Server/web-server"
)

func main() {
	go web_server.StartWebServer()
	go proxy_server.StartProxy()
	go test.INITTEST()

	for {
	}
}
