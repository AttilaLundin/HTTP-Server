package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type CODE int32

var supportedFileTypes = map[string]struct{}{"text/html": {}, "text/plain": {}, "text/css": {}, "image/gif": {}, "image/jpeg": {}, "image/jpg": {}}

// TODO: Kolla om vi får använda http.ResponseWriter

// stateless communication; handle requests not clients per se
func ClientRequestHandler(connection net.Conn, lock *sync.Mutex) {

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
		respond(CODE(408))
		return
	}

	// man har metoder man kan använda för att få åtkomst till requestens info
	request, readRequesError := http.ReadRequest(bufio.NewReader(connection))
	if readRequesError != nil {
		respond(CODE(500))
		return
	}

	switch request.Method {
	case "GET":
		handleGet(request, lock)
	case "POST":
		code := handlePOST(request, lock)
		respond(code)

	case "PUT", "DELETE", "OPTIONS", "PATCH", "TRACE", "CONNECT":
		respond(CODE(500))
	default:
		respond(CODE(400))
	}

}

func handlePOST(request *http.Request, emptyFileMutex *sync.Mutex) CODE {

	//TODO: kolla om det är ok att vå läser in request två gånger

	// Get the file from the request
	file, header, formFileError := request.FormFile("file")
	if formFileError != nil {
		return CODE(500)
	}
	defer file.Close()

	//contentType := request.Header.Get("Content-Type")
	reqBody, err := io.ReadAll(request.Body)
	if err != nil {
		return CODE(500)
	}

	contentType := http.DetectContentType(reqBody)

	// slice away everything after the ; so we simply get e.g. text/plain or text/css without "; utf-8" then trim spaces
	contentType = strings.TrimSpace(strings.Split(contentType, ";")[0])

	if _, ok := supportedFileTypes[contentType]; !ok || contentType == "application/octet-stream" {
		return CODE(408)
	}

	pl("request.URL.Path", request.URL.Path)
	if !strings.HasPrefix(request.URL.Path, "/storage/") {
		return CODE(400)
	}

	lock.Lock()
	defer lock.Unlock()
	emptyFile, creationError := os.Create(request.URL.Path[1:] + "/" + header.Filename)
	if creationError != nil {
		return CODE(500)
	}
	defer emptyFile.Close()

	_, copyError := io.Copy(emptyFile, file)
	if copyError != nil {
		return CODE(500)
	}
	return CODE(200)
}

func handleGet(request *http.Request, emptyFileMutex *sync.Mutex) {
	path := request.URL.Path
	fmt.Println("------------")
	fmt.Println(path)
	emptyFileMutex.Lock()
	defer emptyFileMutex.Unlock()

	if !strings.HasPrefix(request.URL.Path, "/storage/") {
		pl("prefix fel")
		respond(CODE(400))
	}
	path = path[1:]
	_, statErr := os.Stat(path)
	if statErr != nil {
		pl("The file doesnt exist!!!?!?!?!?!?!?! what frågetecken")
	}
	fmt.Println("this is path after reg", path)

	var responseBody bytes.Buffer

	// Create a new multipart writer
	responseBodyWriter := multipart.NewWriter(&responseBody)
	pl(responseBodyWriter.FormDataContentType())
	// Create a form file field with the file name "image.gif"

	pl(responseBodyWriter.FormDataContentType())

	file, openError := os.Open(path)
	if openError != nil {
		pl("openError")
		respond(CODE(404))
	}

	defer file.Close()
	// Create a form file field with the file name "image.gif"

	fmt.Println("-----------dafgs")

	/*fileContents, readError := io.ReadAll(file)
	fmt.Println("-----------dafgs numero 2")
	if readError != nil {
		pl(path)
		pl("AAAAHG OUGA BOOGA")
		log.Fatal(readError)
	}
	*/
	response := &http.Response{
		Status:     "200 OK",           // Setting the status text
		StatusCode: http.StatusOK,      // Setting the status code
		Proto:      "HTTP/1.1",         // Setting the protocol version
		ProtoMajor: 1,                  // Major protocol version
		ProtoMinor: 1,                  // Minor protocol version
		Header:     make(http.Header),  // Initializing the Header map
		Body:       io.NopCloser(file), // Setting the response body
	}
	response.Header.Set("Content-Type", multipart.NewWriter(&responseBody).FormDataContentType())

}

// TODO: Skriv om så att vi bara har en funktion och lägger in koder

func respond(code CODE) {
	switch code {
	case 200:
		pl("should respond with HTTP Status Code 200 a OK")
		//respond with HTTP Status Code 200 a OK
	case 500:
		pl("should respond with HTTP Status Code 501 NotImplemented")
		////respond with HTTP Status Code 501 NotImplemented
	case 400:
		pl("should respond with HTTP Status Code 400 BadRequest")
		//respond with HTTP Status Code 400 BadRequest
	case 408:
		pl("should respond with HTTP Status Code 408 RequestTimeout")
		//respond with HTTP Status Code 408 RequestTimeout
	default:
		pl("should respond with HTTP Status Code 500 InternalServerError")
		//respond with HTTP Status Code 500 InternalServerError

	}
}
