package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
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
	formFile, err := writer.CreateFormFile("file", "testfile.txt")
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
	//storage / textfilestestfile.txt

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", "http://localhost:5431/storage/text/plain", &requestBody)
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

func SendHTMLContent() {
	// Create a buffer to store the POST request body
	var requestBody bytes.Buffer

	// Create a new multipart writer
	writer := multipart.NewWriter(&requestBody)

	// Create a form file field with the file name "testfile.html"
	formFile, err := writer.CreateFormFile("file", "testfile.html")
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return
	}

	// Write the HTML content to the form file
	htmlContent := "<!DOCTYPE html><html><body><h1>Hello, World!</h1></body></html>"
	_, err = formFile.Write([]byte(htmlContent))
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
	req, err := http.NewRequest("POST", "http://localhost:5431/storage/text/html", &requestBody)
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

	// Here you might want to check the response status, handle the response, etc.
	fmt.Println("HTML file sent successfully, status code:", resp.Status)
}

func SendCSSContent() {
	// Create a buffer to store the POST request body
	var requestBody bytes.Buffer

	// Create a new multipart writer
	writer := multipart.NewWriter(&requestBody)

	// Create a form file field with the file name "stylesheet.css"
	formFile, err := writer.CreateFormFile("file", "stylesheet.css")
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return
	}

	// Write the CSS content to the form file
	cssContent := "body { font-family: Arial, sans-serif; } h1 { color: #333366; }"
	_, err = formFile.Write([]byte(cssContent))
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
	req, err := http.NewRequest("POST", "http://localhost:5431/storage/text/css", &requestBody)
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

	// Here you might want to check the response status, handle the response, etc.
	fmt.Println("CSS file sent successfully, status code:", resp.Status)
}

func SendJPGContent() {
	// Open the JPG file from the testimages directory
	file, err := os.Open("testimages/Cat03.jpg")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a buffer to store the POST request body
	var requestBody bytes.Buffer

	// Create a new multipart writer
	writer := multipart.NewWriter(&requestBody)

	// Create a form file field with the file name "image.jpg"
	formFile, err := writer.CreateFormFile("file", "Cat03.jpg")
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return
	}

	// Copy the file content to the form file
	_, err = io.Copy(formFile, file)
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
	req, err := http.NewRequest("POST", "http://localhost:5431/storage/image/jpg", &requestBody)
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

	// Here you might want to check the response status, handle the response, etc.
	fmt.Println("JPG file sent successfully, status code:", resp.Status)
}

func SendJPEGContent() {
	// Open the JPEG file from the testimages directory
	file, err := os.Open("testimages/astronaut-with-pencil-pen-tool-created-clipping-path-included-jpeg-easy-composite.jpeg")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a buffer to store the POST request body
	var requestBody bytes.Buffer

	// Create a new multipart writer
	writer := multipart.NewWriter(&requestBody)

	// Create a form file field with the file name "image.jpeg"
	formFile, err := writer.CreateFormFile("file", "astronaut-with-pencil-pen-tool-created-clipping-path-included-jpeg-easy-composite.jpeg")
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return
	}

	// Copy the file content to the form file
	_, err = io.Copy(formFile, file)
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
	req, err := http.NewRequest("POST", "http://localhost:5431/storage/image/jpeg", &requestBody)
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

	// Here you might want to check the response status, handle the response, etc.
	fmt.Println("JPEG file sent successfully, status code:", resp.Status)
}

func SendGIFContent() {
	// Open the GIF file from the testimages directory
	file, err := os.Open("testimages/skeleton.gif")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a buffer to store the POST request body
	var requestBody bytes.Buffer

	// Create a new multipart writer
	writer := multipart.NewWriter(&requestBody)

	// Create a form file field with the file name "image.gif"
	formFile, err := writer.CreateFormFile("file", "skeleton.gif")
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return
	}

	// Copy the file content to the form file
	_, err = io.Copy(formFile, file)
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
	req, err := http.NewRequest("POST", "http://localhost:5431/storage/image/gif", &requestBody)
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

	// Here you might want to check the response status, handle the response, etc.
	fmt.Println("GIF file sent successfully, status code:", resp.Status)
}
