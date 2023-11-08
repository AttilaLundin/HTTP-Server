package test

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func INITTEST() {
	time.Sleep(time.Second * 3)
	test()
}

func test() {

	go testPostText()
	go testPostGif()
	go testPostHtml()
	go testPostCss()
	go testPostJpg()
	go testPostJpeg()

	time.Sleep(time.Second * 3)
	go testGetGif()
	go testGetText()
	//cba making a channel, so we wait to make sure that the functions have executed
	time.Sleep(time.Second * 2)
}

func testPostText() {
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

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", "http://localhost:5431/web-server/storage/text/plain", &requestBody)
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
	fmt.Println("Plain file sent successfully, status code:", resp.Status)
}

func testPostHtml() {
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
	req, err := http.NewRequest("POST", "http://localhost:5431/web-server/storage/text/html", &requestBody)
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

func testPostCss() {
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
	req, err := http.NewRequest("POST", "http://localhost:5431/web-server/storage/text/css", &requestBody)
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

func testPostJpg() {
	// Open the JPG file from the testimages directory
	file, err := os.Open("test/testimages/Cat03.jpg")
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
	req, err := http.NewRequest("POST", "http://localhost:5431/web-server/storage/image/jpg", &requestBody)
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

func testPostJpeg() {
	// Open the JPEG file from the testimages directory
	file, err := os.Open("test/testimages/astronaut-with-pencil-pen-tool-created-clipping-path-included-jpeg-easy-composite.jpeg")
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

	// Set the content type to the multipart form's content type
	contentType := writer.FormDataContentType()

	// Perform the HTTP POST request
	resp, err := http.Post("http://localhost:5431/web-server/storage/image/jpeg", contentType, &requestBody)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	// Here you might want to check the response status, handle the response, etc.
	fmt.Println("JPEG file sent successfully, status code:", resp.Status)
}

func testPostGif() {
	// Open the GIF file from the testimages directory
	file, err := os.Open("test/testimages/skeleton.gif")
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
	req, err := http.NewRequest("POST", "http://localhost:5431/web-server/storage/image/gif", &requestBody)
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

func testGetGif() {
	// Define the URL to send the GET request to
	url := "http://localhost:5431/web-server/storage/image/gif/skeleton.gif"

	// Create a new HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Error sending HTTP GET request: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println("Statuscode of the response:", resp.StatusCode)

	// Check if the HTTP status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		log.Fatal("Expected HTTP status code 200, got %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body: %v", err)
	}
	emptyFile, creationError := os.Create("test/clienttest/skeleton.gif")
	if creationError != nil {
		log.Fatal(creationError)
	}
	io.Copy(emptyFile, bytes.NewReader(body))
	fmt.Println("200!! GET GIF request and response validation successful")
}

func testGetText() {
	// Define the URL to send the GET request to
	url := "http://localhost:5431/web-server/storage/text/plain/testfile.txt"

	// Create a new HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Error sending HTTP GET request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the HTTP status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		log.Fatal("Expected HTTP status code 200, got %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body: %v", err)
	}

	// Define the expected response body content
	expectedContent := "test"

	// Validate the response body content
	if string(body) != expectedContent {
		log.Fatal("Expected response body to be %q, got %q", expectedContent, body)
	}

	fmt.Println("200!! GET request and response validation successful")
}
