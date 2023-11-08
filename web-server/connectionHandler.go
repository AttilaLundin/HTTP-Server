package web_server

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

const maxClients = 10

var pl = fmt.Println
var lock sync.Mutex

// TODO: double check later if error handling is appropriate

func StartWebServer() {

	// start listening to a port
	tcpListener := setupListener()

	// empty structure because value does not matter
	requestChannel := make(chan struct{}, maxClients)

	for {

		//create an empty anonymous struct, the value or content of the struct does not matter
		requestChannel <- struct{}{} // will block if there is maxClients in the clientsPool

		tcpConnection, err := tcpListener.AcceptTCP()
		if err != nil {
			pl(err)
			continue
		}

		go func() { // create a concurrent request

			requestHandler(tcpConnection, &lock)
			<-requestChannel // removes an entry from clientsPool, allowing another to proceed
			err := tcpConnection.Close()
			if err != nil {
				return
			}
		}()

	}
}

func setupListener() *net.TCPListener {

	port := os.Getenv("PORT")
	//Bind to All Interfaces Inside the Container: Inside the Docker container,
	//you should bind your server to 0.0.0.0.
	//This makes the server listen on all interfaces inside the container.
	ip := "0.0.0.0"
	address := ip + ":" + port

	tcpAddress, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	tcpListener, err := net.ListenTCP("tcp", tcpAddress)
	if err != nil {
		log.Fatal(err)
	}

	pl("Now listening to Address:", address)
	return tcpListener

}

//docker run -p 8080:8080 -e PORT=8080 http-server
