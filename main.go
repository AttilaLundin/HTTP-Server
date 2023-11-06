package main

import (
	proxy_server "HTTP_Server/proxy-server"
	web_server "HTTP_Server/web-server"
)

func main() {

	go web_server.StartWebServer()
	go proxy_server.StartProxyServer()
	for {

	}

}
