package proxy_server

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

// Represents the http status
type code int

func StartProxyServer() {
	// listener for incoming connection
	incomingConnectionListener := setupListener()
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

func handleProxyRequest(incomingConnection *net.TCPConn) {
	// read incoming http request
	request, readRequestError := http.ReadRequest(bufio.NewReader(incomingConnection))
	if readRequestError != nil {
		//if error send response code 500
		code(500).makeAndSendResponse(incomingConnection)
	} else if request.Method != "GET" {
		// if request is not get then send response 501
		code(501).makeAndSendResponse(incomingConnection)
	}
	// retrieve IP address from the request host
	remoteAddress, err := net.ResolveTCPAddr("tcp", request.Host)
	if err != nil {
		code(400).makeAndSendResponse(incomingConnection)
	}
	// establish an outgoing connection with the remote address
	outgoingConnection, err := establishOutgoingConnection(5, remoteAddress)
	if err != nil {
		code(500).makeAndSendResponse(incomingConnection)
	}

	defer func(outgoingConnection *net.TCPConn) {
		err := outgoingConnection.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(outgoingConnection)
	// wend the request to the server
	sendRequest(request, outgoingConnection)
	// set timeout for reading the response
	err = outgoingConnection.SetReadDeadline(time.Now().Add(time.Second * 5))
	if err != nil {
		panic(err)
	}
	// read response from server
	response, err := http.ReadResponse(bufio.NewReader(outgoingConnection), nil)
	if err != nil {
		code(500).makeAndSendResponse(incomingConnection)
	}
	//Send the response back to the client
	sendResponse(response, incomingConnection)

}

// attempts to establish connection to the server by retrying if it fails
// although in our program, practically, if it fails first time it will likely fail on reattempts
func establishOutgoingConnection(nrOfAttempts int, remoteAddress *net.TCPAddr) (*net.TCPConn, error) {
	var err error = nil
	attempt := 1

	for i := 0; i < nrOfAttempts; i++ {
		outgoingConnection, err := net.DialTCP("tcp", nil, remoteAddress)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Proxy server failed to connect to web server... attempt #: ", attempt)
			attempt++
			time.Sleep(time.Second)
			continue
		}
		return outgoingConnection, nil
	}
	return nil, err
}

// creates a http response based on the code and sends it to the client connection
func (code code) makeAndSendResponse(clientConnection *net.TCPConn) {
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

	// convert body to array of bytes so that we can write it to client through connection
	buf := bytes.Buffer{}
	err := response.Write(&buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	//write the bytes to the client connection
	_, err = connection.Write(buf.Bytes())
	if err != nil {
		fmt.Println("sendResponse error: ", err)
	}
}

func sendRequest(request *http.Request, connection *net.TCPConn) {

	// convert body to array of bytes so that we can write it to client through connection
	buf := bytes.Buffer{}
	if err := request.Write(&buf); err != nil {
		fmt.Println("Error 1 is: ", err)
	}
	//write the bytes to server connection
	if _, err := connection.Write(buf.Bytes()); err != nil {
		//pl("Conn is: ", connection)
		fmt.Println("sendRequest error: ", err)
	}
}

func setupListener() *net.TCPListener {
	for {
		// init a reader
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter the ip and port you want the proxy to listen to in the format  <ip>:<port>  below:\n")

		// some ports on windows don't work depending on machine, e.g. 5433. We use 5431 instead.
		address, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			// restart loop until successful
			continue
		}

		// remove delimiter /n or bugs
		address = address[:len(address)-2]

		// retrieve address from an existing tcp connection
		tcpAddress, err := net.ResolveTCPAddr("tcp", address)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// listen to retrieved address from tcp conn
		tcpListener, err := net.ListenTCP("tcp", tcpAddress)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Proxy Server now listening to address:", address)
			return tcpListener
		}
	}
}

//<3 in memory of daString <3
