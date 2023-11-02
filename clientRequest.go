package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type CODE int32

var supportedFileTypes = map[string]struct{}{"text/html": {}, "text/plain": {}, "text/css": {}, "image/gif": {}, "image/jpeg": {}, "image/jpg": {}}

// TODO: Kolla om vi får använda http.ResponseWriter

// stateless communication; handle requests not clients per se
func ClientRequestHandler(connection net.Conn) {

	defer func(connection net.Conn) {
		connectionError := connection.Close()
		if connectionError != nil {
			log.Fatal(connectionError)
			//TODO: ändra
		}
	}(connection)

	timeoutError := connection.SetReadDeadline(time.Now().Add(time.Second * 100))
	if timeoutError != nil {
		log.Println("Error: request timed out")
		respondRequestTimeout()
		return
	}

	// man har metoder man kan använda för att få åtkomst till requestens info
	request, readRequesError := http.ReadRequest(bufio.NewReader(connection))
	if readRequesError != nil {
		respondInternalServerError()
		return
	}

	switch request.Method {
	case "GET":
		handleGet(request)
	case "POST":
		handlePOST(request)
	case "PUT", "DELETE", "OPTIONS", "PATCH", "TRACE", "CONNECT":
		respond(CODE(500))
	default:
		respond(CODE(400))
	}
}

func handlePOST(request *http.Request) {

	//TODO: kolla om det är ok att vå läser in request två gånger

	// Get the file from the request
	file, header, formFileError := request.FormFile("file")
	if formFileError != nil {
		respondInternalServerError()
		return
	}
	defer file.Close()

	//contentType := request.Header.Get("Content-Type")
	reqBody, err := io.ReadAll(request.Body)
	if err != nil {
		respondInternalServerError()
		return
	}

	contentType := http.DetectContentType(reqBody)

	// slice away everything after the ; so we simply get e.g. text/plain or text/css without "; utf-8" then trim spaces
	contentType = strings.TrimSpace(strings.Split(contentType, ";")[0])

	if _, ok := supportedFileTypes[contentType]; !ok || contentType == "application/octet-stream" {
		fmt.Println("this is not a supported type")
		respondBadRequest()
		return
	}

	pl("request.URL.Path", request.URL.Path)
	if !strings.HasPrefix(request.URL.Path, "/storage/") {
		pl("BadRequest")
		respondBadRequest()
		return
	}

	emptyFile, creationError := os.Create(request.URL.Path[1:] + "/" + header.Filename)
	if creationError != nil {
		respondInternalServerError()
		return
	}
	defer emptyFile.Close()

	_, copyError := io.Copy(emptyFile, file)
	if copyError != nil {
		respondInternalServerError()
		return
	}
}

func handleGet(request *http.Request) { /*
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
}

// TODO: Skriv om så att vi bara har en funktion och lägger in koder

func respond(code CODE) {
	switch code {
	case 200:
		//respond with HTTP Status Code 200 a OK
	case 500:
		////respond with HTTP Status Code 501 NotImplemented
	case 400:
		//respond with HTTP Status Code 400 BadRequest
	case 408:
		//respond with HTTP Status Code 408 RequestTimeout
	default:
		//respond with HTTP Status Code 500 InternalServerError

	}
}

func respondRequestNotImplemented() {
	//respond with HTTP Status Code 500
}
func respondBadRequest() {
	//respond with HTTP Status Code 400
}
func respondRequestTimeout() {
	//respond with HTTP Status Code 408
}
func respondInternalServerError() {
	//respond with HTTP Status Code 500
}
