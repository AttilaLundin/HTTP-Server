package web_server

import (
	"fmt"
	"net"
	"sync"
)

const maxClients = 10

var lock sync.Mutex
var testLock sync.Mutex

// TODO: double check later if error handling is appropriate

func StartWebServer() {
	//TODO: comment out when not testing

	// empty structure because value does not matter
	requestChannel := make(chan struct{}, maxClients)

	// start listening to a port
	tcpListener := setupListener()

	for {

		//create an empty anonymous struct, the value or content of the struct does not matter
		requestChannel <- struct{}{} // will block if there is maxClients in the clientsPool

		tcpConnection, err := tcpListener.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go func() { // create a concurrent request

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
			fmt.Println(err)
			continue
		}

		tcpListener, err := net.ListenTCP("tcp", tcpAddress)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Now listening to Address:", address)
			return tcpListener
		}
	}
}
