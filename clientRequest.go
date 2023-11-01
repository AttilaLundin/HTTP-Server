package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"strings"
)

// stateless communication; handle requests not clients per se
func ClientRequestHandler(connection net.Conn) {

	// not sure if we need this
	// client := &http.Client{}

	defer func(connection net.Conn) {
		err := connection.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(connection)

	//TODO: förmodligen introducera en buffer som läser in data
	// Parsing av url, regex*? (Go docs: net/http, ServerMux)

	// man har metoder man kan använda för att få åtkomst till requestens info
	request, err := http.ReadRequest(bufio.NewReader(connection))
	if err != nil {
		fmt.Println("Error handling the request", err)
		return
	}

	// todo: maybe transition to regex l8er
	notImplemented := map[string]struct{}{"PUT": {}, "DELETE": {}, "OPTIONS": {}, "PATCH": {}, "TRACE": {}, "CONNECT": {}}
	supportedFileTypes := map[string]struct{}{"text/html": {}, "text/plain": {}, "text/css": {}, "image/gif": {}, "image/jpeg": {}, "image/jpg": {}}

	if request.Method == "GET" {
		/*
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

		*/

		//TODO: hantera get requesten
	} else if request.Method == "POST" {

		//contentType := request.Header.Get("Content-Type")
		reqBody, err := io.ReadAll(request.Body)
		if err != nil {
			println(err)
			return
		}
		contentType := http.DetectContentType(reqBody)
		// slice away everything after the ; so we simply get e.g. text/plain or image/gif without "; utf-8" then trim spaces
		contentType = strings.TrimSpace(strings.Split(contentType, ";")[0])
		pl(contentType)

		if _, ok := supportedFileTypes[contentType]; !ok || contentType == "application/octet-stream" {
			pl("this is not a supported type")
			//TODO: appropriate response
		}

		//pl("Content Type: ", contentType)
		//pl(reflect.TypeOf(contentType))

		// Get the file from the request
		file, header, err := request.FormFile("file")
		if err != nil {
			pl(err)
			return
		}

		defer func(file multipart.File) {
			err := file.Close()
			if err != nil {
				pl(err)
			}
		}(file)

		emptyFile, err := os.Create(header.Filename)
		if err != nil {
			pl(err)
			return
		}
		_, err = io.Copy(emptyFile, file)
		if err != nil {
			pl(err)
		}

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
