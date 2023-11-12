package web_server

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
	"strconv"
	"strings"
	"sync"
	"time"
)

// used for HTTP response code
type code int

// Struct that denotes the contents of a get response body
type getResponse struct {
	fileInBytes []byte
	fileInfo    os.FileInfo
}

// map containing all the supported file types, used to easily check if the contents in the received post is valid.
var supportedFileTypes = map[string]struct{}{"text/html": {}, "text/plain": {}, "text/css": {}, "image/gif": {}, "image/jpeg": {}, "image/jpg": {}}

// stateless communication - handle requests not clients per se
func requestHandler(connection *net.TCPConn, lock *sync.Mutex) {

	// put a deadline on a connection for safe measure
	timeoutError := connection.SetReadDeadline(time.Now().Add(time.Second * 60))
	if timeoutError != nil {
		log.Println("Error: request timed out")
		code(400).makeAndSendResponse(connection)
		return
	}

	// read request from our connection
	request, readRequestError := http.ReadRequest(bufio.NewReader(connection))
	if readRequestError != nil {
		code(500).makeAndSendResponse(connection)
	}

	// handles the request depending on what typ of http request it is.
	switch request.Method {
	case "GET":
		code, getResponse := handleGet(request, lock)
		if code == http.StatusOK {
			getResponse.makeAndSendResponse(connection)
		} else {
			code.makeAndSendResponse(connection)
		}
	case "POST":
		code := handlePOST(request, lock)
		code.makeAndSendResponse(connection)
	case "PUT", "DELETE", "OPTIONS", "PATCH", "TRACE", "CONNECT":
		code(500).makeAndSendResponse(connection)
	default:
		code(400).makeAndSendResponse(connection)
	}
}

func handlePOST(request *http.Request, lock *sync.Mutex) code {
	// make sure we can only post to our storage directory
	if !strings.HasPrefix(request.URL.Path, "/web-server/storage") {
		return code(400)
	}

	// get the file from the request
	file, header, err := request.FormFile("file")
	if err != nil {
		return code(500)
	}

	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// read and store request body
	reqBody, err := io.ReadAll(request.Body)
	if err != nil {
		return code(500)
	}

	contentType := http.DetectContentType(reqBody)

	// slice away everything after the ; so we simply get e.g. text/plain or text/css without "; utf-8"
	// then trim spaces for safe measure
	contentType = strings.TrimSpace(strings.Split(contentType, ";")[0])

	// check if content type is supported, sets supportedContent to true or false
	_, supportedContent := supportedFileTypes[contentType]
	if !supportedContent || contentType == "application/octet-stream" {
		return code(408)
	}

	// lock to avoid interleaving when creating and copying to storage, shared data is being manipulated
	lock.Lock()
	defer lock.Unlock()

	emptyFile, err := os.Create(request.URL.Path[1:] + "/" + header.Filename)
	if err != nil {
		return code(500)
	}

	defer func(emptyFile *os.File) {
		err := emptyFile.Close()
		if err != nil {

		}
	}(emptyFile)

	_, err = io.Copy(emptyFile, file)
	if err != nil {

		return code(500)
	}
	return code(200)
}

// function that handles get requests, returns either an error code and nil, or 200 and the requested file as a struct
func handleGet(request *http.Request, lock *sync.Mutex) (code, getResponse) {

	//relative path of requested file
	path := request.URL.Path

	// makes sure client's requested file is in the designated storage space
	if !strings.HasPrefix(request.URL.Path, "/web-server/storage") {
		return code(400), getResponse{}
	}

	//removes the first "/" in the relative path
	path = path[1:]

	//check if the file exists, if not we respond with error message.
	fileInfo, err := os.Stat(path)
	if err != nil {
		return code(400), getResponse{}
	}

	//lock the central lock to avoid interleaving and race conditions.
	lock.Lock()
	defer lock.Unlock()

	//fetching the requested file
	fileInBytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return code(404), getResponse{}
	}

	return code(200), getResponse{fileInBytes: fileInBytes, fileInfo: fileInfo}

}

// for successful GET response
func (gr getResponse) makeAndSendResponse(connection net.Conn) {
	response := &http.Response{
		Status:     http.StatusText(200),                          // Setting the status text
		StatusCode: http.StatusOK,                                 // Setting the status code ex, 200 for success
		Proto:      "HTTP/1.1",                                    // Setting the protocol version
		ProtoMajor: 1,                                             // Major protocol version
		ProtoMinor: 1,                                             // Minor protocol version
		Header:     make(http.Header),                             // Initializing the Header map
		Body:       io.NopCloser(bytes.NewReader(gr.fileInBytes)), // Setting the response body
	}
	response.Header.Set("Content-Type", http.DetectContentType(gr.fileInBytes))
	response.Header.Set("Connection", "Closed")
	// convert integer to decimal string representation
	response.Header.Set("Content-Length", strconv.Itoa(int(gr.fileInfo.Size())))
	response.Header.Set("Last-Modified", gr.fileInfo.ModTime().String())
	sendResponse(response, connection)
}

// for successful POST request or unsuccessful GET or POST request
func (code code) makeAndSendResponse(connection net.Conn) {
	response := &http.Response{
		Status:     http.StatusText(int(code)), // Setting the status text
		StatusCode: int(code),                  // Setting the status code e.g. 200 for success
		Proto:      "HTTP/1.1",                 // Setting the protocol version
		ProtoMajor: 1,                          // Major protocol version
		ProtoMinor: 1,                          // Minor protocol version
		Header:     make(http.Header),          // Init header map
	}
	sendResponse(response, connection)

}

// sends the response to the client
func sendResponse(response *http.Response, connection net.Conn) {
	// convert body to array of bytes so that we can write it to client through connection
	buf := bytes.Buffer{}
	if err := response.Write(&buf); err != nil {
		panic(err)
	}
	if _, err := connection.Write(buf.Bytes()); err != nil {
		panic(err)
	}
}
