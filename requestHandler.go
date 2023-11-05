package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type CODE int

type getResponse struct {
	fileInBytes []byte
	fileInfo    os.FileInfo
}

// todo: dubbelkolla så att vi inte missar något när vi checkar "properly formatted" http
var supportedFileTypes = map[string]struct{}{"text/html": {}, "text/plain": {}, "text/css": {}, "image/gif": {}, "image/jpeg": {}, "image/jpg": {}}

// stateless communication; handle requests not clients per se
func requestHandler(connection *net.TCPConn, lock *sync.Mutex) {
	//defer connection.Close()
	timeoutError := connection.SetReadDeadline(time.Now().Add(time.Second * 60))
	if timeoutError != nil {
		log.Println("Error: request timed out")
		CODE(400).respond(connection)

		return
	}

	// man har metoder man kan använda för att få åtkomst till requestens info
	request, readRequestError := http.ReadRequest(bufio.NewReader(connection))
	if readRequestError != nil {
		CODE(500).respond(connection)
	}

	switch request.Method {
	case "GET":
		code, getResponse := handleGet(request, lock)
		if code == http.StatusOK {
			getResponse.respond(connection)
		} else {
			code.respond(connection)
		}

	case "POST":
		code := handlePOST(request, lock)
		code.respond(connection)
	case "PUT", "DELETE", "OPTIONS", "PATCH", "TRACE", "CONNECT":
		CODE(500).respond(connection)
	default:
		CODE(400).respond(connection)
	}
}

func handlePOST(request *http.Request, lock *sync.Mutex) CODE {

	if !strings.HasPrefix(request.URL.Path, "/web-server/storage/") {
		return CODE(400)
	}

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

	// slice away everything after the ; so we simply get e.g. text/plain or text/css without "; utf-8" then trim spaces
	contentType := http.DetectContentType(reqBody)
	contentType = strings.TrimSpace(strings.Split(contentType, ";")[0])

	if _, ok := supportedFileTypes[contentType]; !ok || contentType == "application/octet-stream" {
		return CODE(408)
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

func handleGet(request *http.Request, lock *sync.Mutex) (CODE, getResponse) {

	path := request.URL.Path

	if !strings.HasPrefix(request.URL.Path, "/web-server/storage/") {
		return CODE(400), getResponse{}
	}

	path = path[1:]
	fileInfo, statErr := os.Stat(path)
	if statErr != nil {
		return CODE(400), getResponse{}
	}

	lock.Lock()
	defer lock.Unlock()

	fileInBytes, readError := os.ReadFile(path)
	if readError != nil {
		pl("openError")
		return CODE(404), getResponse{}
	}

	return CODE(200), getResponse{fileInBytes: fileInBytes, fileInfo: fileInfo}

}

func (gr getResponse) respond(connection net.Conn) {
	response := &http.Response{
		Status:     "200 OK",                                      // Setting the status text
		StatusCode: http.StatusOK,                                 // Setting the status code
		Proto:      "HTTP/1.1",                                    // Setting the protocol version
		ProtoMajor: 1,                                             // Major protocol version
		ProtoMinor: 1,                                             // Minor protocol version
		Header:     make(http.Header),                             // Initializing the Header map
		Body:       io.NopCloser(bytes.NewReader(gr.fileInBytes)), // Setting the response body
	}
	response.Header.Set("Content-Type", http.DetectContentType(gr.fileInBytes))
	response.Header.Set("Connection", "Closed")
	response.Header.Set("Content-Length", strconv.Itoa(int(gr.fileInfo.Size())))
	response.Header.Set("Last-Modified", gr.fileInfo.ModTime().String())
	sendResponse(response, connection)
}

func (code CODE) respond(connection net.Conn) {
	response := &http.Response{
		Status:     http.StatusText(int(code)),
		StatusCode: int(code),
		Proto:      "HTTP/1.1",        // Setting the protocol version
		ProtoMajor: 1,                 // Major protocol version
		ProtoMinor: 1,                 // Minor protocol version
		Header:     make(http.Header), // Initializing the Header map
	}
	sendResponse(response, connection)

}

func sendResponse(response *http.Response, connection net.Conn) {
	// convert body to array of bytes so we can write it to client through connection
	buf := bytes.Buffer{}
	if err := response.Write(&buf); err != nil {
		panic(err)
	}
	if _, err := connection.Write(buf.Bytes()); err != nil {
		panic(err)
	}
}
