package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
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
			clientRequestHandler(tcpConnection)
			<-clientsPool // removes an entry from clientsPool, allowing another to proceed
		}()
	}

}

// stateless communication; handle requests not clients per se
func clientRequestHandler(connection net.Conn) {

	defer func(connection net.Conn) {
		err := connection.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(connection)

	//TODO: förmodligen introducera en buffer som läser in data
	// Parsing av url, regex*? (Go docs: net/http, ServerMux)

	request, err := http.ReadRequest(bufio.NewReader(connection))
	if err != nil {
		pl("Error handling the request", err)
		return
	}

	notImplemented := map[string]struct{}{"PUT": {}, "DELETE": {}, "OPTIONS": {}, "PATCH": {}, "TRACE": {}, "CONNECT": {}}

	if request.Method == "GET" {
		path := request.URL.Path
		file, err := os.Open(path)
		if err != nil {
			//TODO: If the requested file does not exist, your server should return a well-formed 404 "Not Found" code.
		}

		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(file)

		//file2, err := os.OpenFile("data.txt", os.O_RDONLY, 0644)
		fileContents, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal(err)
		}

		connection.Write(fileContents)

		//TODO: hantera get requesten
	} else if request.Method == "POST" {
		//TODO: hantera post requesten
	} else if _, ok := notImplemented[request.Method]; ok {
		//TODO: "Not Implemented" (501)
	} else {
		//TODO: "Bad Request" (400)
	}

	/*
		err := connection.SetReadDeadline(time.Now().Add(time.Second * 100))
		if err != nil {
			pl("Error: request timed out")
			//http response 408 request timeout
			return
		}*/

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

func client() {
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
