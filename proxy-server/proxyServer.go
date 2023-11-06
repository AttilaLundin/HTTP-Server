package proxy_server

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/http"
	"time"
)

type CODE int

var pl = fmt.Println

func ProxyMain() {

	incomingConnectionListener := setupListener()
	for {
		incomingConnection, err := incomingConnectionListener.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go func() {
			handleProxyRequest(incomingConnection)
		}()
	}
}

func handleProxyRequest(incomingConnection *net.TCPConn) {

	request, readRequestError := http.ReadRequest(bufio.NewReader(incomingConnection))
	if readRequestError != nil {
		CODE(500).makeAndSendResponse(incomingConnection)
	} else if request.Method != "GET" {
		CODE(501).makeAndSendResponse(incomingConnection)
	}

	remoteAddress, err := net.ResolveTCPAddr("tcp", request.Host)
	if err != nil {
		CODE(400).makeAndSendResponse(incomingConnection)
	}

	outgoingConnection, err := establishOutgoingConnection(5, remoteAddress)
	if err != nil {
		CODE(500).makeAndSendResponse(incomingConnection)
	}
	defer func(outgoingConnection *net.TCPConn) {
		err := outgoingConnection.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(outgoingConnection)

	sendRequest(request, outgoingConnection)

	err = outgoingConnection.SetReadDeadline(time.Now().Add(time.Second * 5))
	if err != nil {
		panic(err)
	}

	response, err := http.ReadResponse(bufio.NewReader(outgoingConnection), nil)
	if err != nil {
		CODE(500).makeAndSendResponse(incomingConnection)
	}

	//TODO: Ta bort
	fmt.Println("The response is", response, "with response", response.StatusCode)

	sendResponse(response, incomingConnection)

}

func establishOutgoingConnection(nrOfAttempts int, remoteAddress *net.TCPAddr) (*net.TCPConn, error) {
	var err error = nil
	attempt := 1

	for i := 0; i < nrOfAttempts; i++ {
		outgoingConnection, err := net.DialTCP("tcp", nil, remoteAddress)
		if err != nil {
			pl(remoteAddress)
			pl(err)
			pl("Proxy server failed to connect to web server... attempt #: ", attempt)
			attempt++
			time.Sleep(time.Second)
			continue
		}
		return outgoingConnection, nil
	}
	return nil, err
}

func (code CODE) makeAndSendResponse(clientConnection *net.TCPConn) {
	response := &http.Response{
		Status:     http.StatusText(int(code)),
		StatusCode: int(code),
		Proto:      "HTTP/1.1",        // Setting the protocol version
		ProtoMajor: 1,                 // Major protocol version
		ProtoMinor: 1,                 // Minor protocol version
		Header:     make(http.Header), // Initializing the Header map
	}
	sendResponse(response, clientConnection)
}

func sendResponse(response *http.Response, connection *net.TCPConn) {

	// convert body to array of bytes so we can write it to client through connection
	buf := bytes.Buffer{}
	err := response.Write(&buf)
	if err != nil {
		pl("Error 1 is: ", err)
	}

	_, err = connection.Write(buf.Bytes())
	if err != nil {
		pl("THE CONNECTION IN SENDREQUEST IS: ", connection)
		pl("Error 2 is: ", err)
	}
}

func sendRequest(request *http.Request, connection *net.TCPConn) {

	// convert body to array of bytes so we can write it to client through connection
	buf := bytes.Buffer{}
	if err := request.Write(&buf); err != nil {
		pl("Error 1 is: ", err)
	}
	if _, err := connection.Write(buf.Bytes()); err != nil {
		//pl("Conn is: ", connection)
		pl("Error 2 is: ", err)
	}
}

func setupListener() *net.TCPListener {
	for {

		// TODO: remove comments when not in testing
		/*reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter the ip and port you want to listen to in the format  <ip>:<port>  below:")
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
		address := "localhost:5430" /*+ address*/

		tcpAddress, err := net.ResolveTCPAddr("tcp", address)
		if err != nil {
			pl(err)
			continue
		}

		tcpListener, err := net.ListenTCP("tcp", tcpAddress)
		if err != nil {
			fmt.Println(err)
		} else {
			pl("Proxy Server now listening to Address:", address)
			return tcpListener
		}
	}
}

//<3 in memory of daString <3
