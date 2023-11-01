package main

import (
	"fmt"
	"net"
)

func Client() {
	conn, err := net.Dial("tcp", "localhost:5431")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte("GET / HTTP/1.1\r\nHost: localhost\r\n\r\n"))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
}
