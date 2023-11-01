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

// stateless communication; handle requests not clients per se
func ClientRequestHandler(connection net.Conn) {

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
		fmt.Println("Error handling the request", err)
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
