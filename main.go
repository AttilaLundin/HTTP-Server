package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
)

const MAX_CLIENTS = 10

var pl = fmt.Println

var lock sync.Mutex
var responseWriter http.ResponseWriter

// TODO: double check later if error handling is appropriate

func main() {
	//TODO: comment out when not testing
	go INITTEST()

	// start listening to a port
	tcpListener := setupListener()

	// empty structure because value does not matter
	requestChannel := make(chan struct{}, MAX_CLIENTS)
	for {
		tcpConnection, err := tcpListener.AcceptTCP()
		err = tcpConnection.SetKeepAlive(true)
		if err != nil {
			pl(err)
		}

		//create an empty anonymous struct, the value or content of the struct does not matter
		requestChannel <- struct{}{} // will block if there is MAX_CLIENTS in the clientsPool

		go func() { // create a concurrent request

			ClientRequestHandler(tcpConnection, &lock)
			<-requestChannel // removes an entry from clientsPool, allowing another to proceed
			err := tcpConnection.Close()
			if err != nil {
				return
			}
		}()

	}
}

func setupListener() *net.TCPListener {
	for {
		// TODO: remove comments when not in testing
		/*reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter the port you want to listen to:")
		// some ports on windows don't work depending on machine, e.g. 5433. We use 5431 instead.
		address, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			// restart loop until successful
			continue
		}

		// remove delimiter /n or bugs
		address = address[:len(address)-2]
		*/
		address := "localhost:5431" /*+ address*/

		tcpAddress, err := net.ResolveTCPAddr("tcp", address)
		if err != nil {
			pl(err)
			continue
		}

		tcpListener, err := net.ListenTCP("tcp", tcpAddress)
		if err != nil {
			fmt.Println(err)
		} else {
			pl("Now listening to Address:", address)
			return tcpListener
		}
	}
}
