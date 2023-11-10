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

type CODE int

var pl = fmt.Println

/*
	func ProxyRequestHandler(connection net.Conn, lock *sync.Mutex) {
		fmt.Println("-----")
		request, readRequestError := http.ReadRequest(bufio.NewReader(connection))
		if readRequestError != nil {
			log.Fatal(readRequestError)
		}
		HandleProxyRequest(connection, request, lock)

}
*/
func StartProxyServer() {

	incomingConnectionListener := setupListener()
	for {
		incomingConnection, err := incomingConnectionListener.AcceptTCP()
		if err != nil {
			continue
		}

		//err = incomingConnection.SetKeepAlive(true)
		if err != nil {
			pl(err)
		}

		go func() {
			HandleProxyRequest(incomingConnection)
		}()
	}
}

func HandleProxyRequest(incomingConnection *net.TCPConn) {

	request, readRequestError := http.ReadRequest(bufio.NewReader(incomingConnection))
	if readRequestError != nil {
		CODE(500).respond(incomingConnection)
	} else if request.Method != "GET" {
		CODE(501).respond(incomingConnection)
	}
	pl("request.Header: ", request.Header)
	pl("request.Header: ", request.Header)
	pl("request.Body: ", request.Body)

	//pl("request.host", request.Host)
	//pl("request.URL.Host", request.URL.Host)
	//pl("request.URL.RequestURI()", request.URL.RequestURI())
	//pl("request.URL.Port()", request.URL.Port()) //
	//
	//conn, err := net.Dial("tcp", "localhost:5431")
	//if err != nil {
	//	pl("Error in net.Dial: ", err)
	//}
	//pl("..........")
	//pl("connected with regular dial hoohwooo", conn.RemoteAddr())
	//pl("..........")
	//conn.Close()

	remoteAddress, err := net.ResolveTCPAddr("tcp", request.Host)
	if err != nil {
		pl("WE ARE HERE: ", err)
		//	hantera
	}

	outgoingConnection, err := establishOutgoingConnection(5, remoteAddress)
	if err != nil {
		CODE(500).respond(incomingConnection)
	}
	defer func(outgoingConnection *net.TCPConn) {
		err := outgoingConnection.Close()
		if err != nil {
			panic(err)
		}
	}(outgoingConnection)

	sendRequest(request, outgoingConnection)

	err = outgoingConnection.SetReadDeadline(time.Now().Add(time.Second * 5))
	if err != nil {
		panic(err)
	}

	response, err := http.ReadResponse(bufio.NewReader(outgoingConnection), nil)
	if err != nil {
		//problem med att läsa in responsen
	}

	pl(response)

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

func (code CODE) respond(clientConnection *net.TCPConn) {
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
	if err := response.Write(&buf); err != nil {
		panic(err)
	}
	if _, err := connection.Write(buf.Bytes()); err != nil {
		panic(err)
	}
}
func sendRequest(request *http.Request, connection *net.TCPConn) {
	// convert body to array of bytes so we can write it to client through connection
	buf := bytes.Buffer{}
	if err := request.Write(&buf); err != nil {
		panic(err)
	}
	if _, err := connection.Write(buf.Bytes()); err != nil {
		panic(err)
	}
}

func setupListener() *net.TCPListener {

	ip := os.Getenv("IP")
	if ip == "" {
		log.Fatal("Can't start proxy without a valid ip number")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Can't start proxy without a valid port number")
	}

	address := ip + ":" + port

	tcpAddress, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	tcpListener, err := net.ListenTCP("tcp", tcpAddress)
	if err != nil {
		log.Fatal(err)
	}

	pl("Proxy Server now listening to Address:", address)
	return tcpListener

}

//<3 in memory of daString <3
