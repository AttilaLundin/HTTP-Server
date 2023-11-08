package main

import (
	proxy_server "HTTP_Server/proxy-server"
	web_server "HTTP_Server/web-server"
)

func main() {
	web_server.StartWebServer()
	proxy_server.StartProxyServer()

}

//docker run -p 8080:8080 -e PORT=8080 http-server
