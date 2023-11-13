package proxy_server

import (
	"fmt"
	"log"
	"net"
	"os"
)

// Represents the http status
type code int

func StartProxyServer() {
	// listener for incoming connection
	incomingConnectionListener := setupListener()

	//main loop of the proxy
	for {
		// accept incoming TCP connections
		incomingConnection, err := incomingConnectionListener.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go func() {
			// function for handle incoming TCP connection in a separate goroutine
			handleProxyRequest(incomingConnection)
		}()
	}
}

func setupListener() *net.TCPListener {

	ip := os.Getenv("IP")
	if ip == "" {
		log.Fatal("Not a valid ip address")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Not a valid port")
	}

	address := ip + ":" + port

	// retrieve address from an existing tcp connection
	tcpAddress, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	// listen to retrieved address from tcp conn
	tcpListener, err := net.ListenTCP("tcp", tcpAddress)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Proxy Server now listening to address:", address)
	return tcpListener
}

//<3 in memory of daString <3
