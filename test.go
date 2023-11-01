package main

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net"
	"net/http"
)

func Client() {
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
func SendTextContent() {
	// Create a buffer to store the POST request body
	var requestBody bytes.Buffer

	// Create a new multipart writer
	writer := multipart.NewWriter(&requestBody)

	// Create a form file field with the file name "testfile.txt"
	formFile, err := writer.CreateFormFile("file", "shazbigfartings.txt")
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return
	}

	// Write the text content "test" to the form file
	_, err = formFile.Write([]byte("test"))
	if err != nil {
		fmt.Println("Error writing to form file:", err)
		return
	}

	// Close the multipart writer to finalize the POST request body
	err = writer.Close()
	if err != nil {
		fmt.Println("Error closing writer:", err)
		return
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", "http://localhost:5431", &requestBody)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	// Set the Content-Type header to the multipart form's content type
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the HTTP request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	defer resp.Body.Close()

}
