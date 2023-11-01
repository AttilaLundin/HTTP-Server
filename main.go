package main

import (
	"fmt"
	"net"
)

const MAX_CLIENTS = 10

var pl = fmt.Println

// TODO: double check later if error handling is appropriate
func main() {
	// read_line = strings.TrimSuffix(read_line, "\n")
	// start listening to a port
	listener := setupListener()

	// empty structure because value does not matter
	clientsPool := make(chan struct{}, MAX_CLIENTS)
	for {
		//TODO: avkommentera vid testning
		//go func() {
		//	client()
		//}()

		tcpConnection, err := listener.Accept()
		if err != nil {
			pl(err)
			continue
		}
		// create an empty anonymous struct, the value or content of the struct does not matter
		clientsPool <- struct{}{} // will block if there is MAX_CLIENTS in the clientsPool

		go func() { // create a concurrent request
			ClientRequestHandler(tcpConnection)
			<-clientsPool // removes an entry from clientsPool, allowing another to proceed
		}()
	}

}

func setupListener() net.Listener {
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

		listener, err := net.Listen("tcp", address)
		if err != nil {
			fmt.Println(err)
		} else {
			pl("Now listening to Address:", address)
			return listener
		}
	}
}
