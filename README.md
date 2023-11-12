# **HTTP Web Server and Proxy Server**

In order to start both the HTTP and Proxy server you need to run main. 
You will be prompted to enter the ip address and port for the web server and proxy server.
We have lazily implemented it so that both are started by goroutines immediately after each other: 
sometimes you are prompted for the web server first and sometimes for the proxy server first, 
the prompt that appears first is the one that will also read your input first! 
We generally use port 5431 for the web server and port 5430 for the proxy server.
However, in the cloud we have containerized the web server and proxy server separately using Docker so that they run in separate instances.




### **Testing**
There are some tests that you can start by uncommenting the following line in row 26 in the connectionhandler:

`//go test.INITTEST()`

It is important that you have the correct repository structure if you wish to use these tests(!) 

Within the webserver folder there must be a subfolder  

    HTTP-Server/web-server/storage








# **The Basics (10 points)**

Your task is to build a web sever capable of accepting HTTP requests and returning response data from locally stored files to a client. The server will be implemented in Go and MUST handle concurrent requests by creating a Go routine for each new client request. You will only be responsible for implementing the GET and POST methods. All other request methods received by the server should elicit a "Not Implemented" (501) error (see RFC 1945Links to an external site. section 9.5 - Server Error).

Your web server should compile and run without errors or warnings, producing a binary called http_server that takes as its first argument a port to listen from. Don't use a hard-coded port number. You shouldn't assume that your server will be running on a particular IP address, or that clients will be coming from a pre-determined IP.

## **Listening**

When your server starts, the first thing that it will need to do is establish a socket connection that it can use to listen for incoming connections. Your server should listen on the port specified from the command line and wait for incoming client connections. Each new client request is accepted, and a new Go routine is spawned to handle the request. To avoid overwhelming your server, you should not create more than a reasonable number of child processes (for this assignment, use at most 10). In case an additional child process would break this limit, your server should wait until one of its ongoing child processes exits before forking a new one to handle the new request.
Once a client has connected, the server should read data from the client and then check for a properly-formatted HTTP request. Your server should accept requests for files ending in html, txt, gif, jpeg, jpg, or css and transmit them to the client with a Content-Type of text/html, text/plain, image/gif, image/jpeg, image/jpeg, or text/css, respectively. If the client requests a file with any other extension, the web server must respond with a well-formed 400 "Bad Request" code. An invalid request from the client should be answered with an appropriate error code, i.e. "Bad Request" (400) or "Not Implemented" (501) for valid HTTP methods other than GET. If the requested file does not exist, your server should return a well-formed 404 "Not Found" code. Similarly, if headers are not properly formatted for parsing or any other error condition not listed before, your server should also generate a type-400 message.  For POST requests, please make sure that you store the files and make them accessible with a subsequent GET request.

## **Parsing and Networking Libraries in Go**

For this assignment, you should use the package `net` for the networking, for example using `net.Listen("tcp", address)` to listen for incoming TCP connections. You can also use the package `net/http`, but ONLY for parsing and working with HTTP request objects, and not the networking part. You should not use e.g., `http.ListenAndServe` which trivializes the assignment (the same goes for `http.Listen`, and `http.Serve`).

## **Testing**

There are no included tests for this assignment. However, to make sure your code works, you should set up some way of testing its functionalities yourself.



# **Bonus part (10 points):**

Ordinarily, HTTP is a client-server protocol. The client (usually your web browser) communicates directly with the server (the web server software). However, in some circumstances it may be useful to introduce an intermediate entity called a proxy. Conceptually, the proxy sits between the client and the server. In the simplest case, instead of sending requests directly to the server, the client sends all of its requests to the proxy. The proxy then opens a connection to the server, and passes on the client's request. The proxy receives the reply from the server, and then sends that reply back to the client. Notice that the proxy is essentially acting like both an HTTP client (to the remote server) and an HTTP server (to the initial client).

Why use a proxy? There are a few possible reasons:

* **Performance:** By saving a copy of the pages that it fetches, a proxy can reduce the need to create connections to remote servers. This can reduce the overall delay involved in retrieving a page, particularly if a server is remote or under heavy load.
* **Content Filtering and Transformation:** While in the simplest case the proxy merely fetches a resource without inspecting it, there is nothing that says that a proxy is limited to blindly fetching and serving files. The proxy can inspect the requested URL and selectively block access to certain domains, reformat web pages (for instances, by stripping out images to make a page easier to display on a handheld or other limited-resource client), or perform other transformations and filtering.
* **Privacy:** Normally, web servers log all incoming requests for resources. This information typically includes at least the IP address of the client, the browser, or other client program that they are using (called the User-Agent), the date and time, and the requested file. If a client does not wish to have this personally identifiable information recorded, routing HTTP requests through a proxy is one solution. All requests coming from clients using the same proxy appear to come from the IP address and User-Agent of the proxy itself, rather than the individual clients. If a number of clients use the same proxy (say, an entire business or university), it becomes much harder to link a particular HTTP transaction to a single computer or individual.
Your task is to build a web proxy capable of accepting HTTP requests, forwarding requests to remote (origin) servers, and returning response data to a client. The proxy will be implemented in Go and MUST handle concurrent requests by creating a Go routine for each new client request. You will only be responsible for implementing the GET method. All other request methods received by the proxy should elicit a "Not Implemented" (501) error (see RFC 1945Links to an external site. section 9.5 - Server Error).

Your proxy implementation should compile and run (using go build) without errors or warnings, producing a binary called proxy that takes as its first argument a port to listen from. Don't use a hard-coded port number.

You shouldn't assume that your proxy will be running on a particular IP address, or that clients will be coming from a pre-determined IP.

Here is an example of how you can test sending a get request to the server, through your proxy, with curl:
$ curl -X GET <server_ip>:<server_port>/<file> -x <proxy_ip>:<proxy_port>

## **Listening**

When your proxy starts, the first thing that it will need to do is establish a socket connection that it can use to listen for incoming connections. Your proxy should listen on the port specified from the command line and wait for incoming client connections. Each new client request is accepted, and a new Go routine is spawned to handle the request.

Once a client has connected, the proxy should read data from the client and then check for a properly-formatted HTTP request. Go provides packages to parse the HTTP request lines and headers. Specifically, you will use the package net/http to ensure that the proxy receives a request that contains a valid request line (see the sever description above for details about HTTP lines and headers). You should NOT use any Proxy method from the http package (http.Proxy).



## **Suggestions to help you build this and reuse it later for further labs**

Look at either Docker, LXC or Vagrant to implement your server. This way you will never have portability issues, and it will make your solution cleaner. Bonus for this will also be considered based on the cleanliness of your code/solution. 