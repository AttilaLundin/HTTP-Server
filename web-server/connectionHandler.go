package web_server

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

// A shared lock between each go routine in order to avoid interleaving, this can be seen as coarse grained lock for the whole server,
// if you wish to optimize the locking then implement a fine grain lock that only locks the targeted file.
var lock sync.Mutex

// a constant for max nr of simultaneous concurrent client requests the server handles
const maxClients = 10

// main function of the web server, handles the incoming connection and fork them off to different goroutines
func StartWebServer() {

	// empty structure because value does not matter. Empty struct takes no memory space.
	requestChannel := make(chan struct{}, maxClients)

	// start listening to a port
	tcpListener := setupListener()

	// this is the main loop of the web server, handling each connection.
	for {

		//create an empty anonymous struct, the value or content of the struct does not matter
		requestChannel <- struct{}{} // will block if there is maxClients in the clientsPool

		// accept incoming tcp connection, blocks until a client tries to connect
		tcpConnection, err := tcpListener.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go func() { // create a concurrent request handler
			requestHandler(tcpConnection, &lock)
			<-requestChannel // removes an entry from clientsPool, allowing another to proceed
			err := tcpConnection.Close()
			if err != nil {
				fmt.Println(err)
				return
			}
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

	fmt.Println("Web server now listening to address:", address)
	return tcpListener
}
