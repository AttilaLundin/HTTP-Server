package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

var pl = fmt.Println

// TODO: double check later if error handling is appropriate
func main() {
	// read_line = strings.TrimSuffix(read_line, "\n")
	// start listening to a port
	listener := setupListener()

	var waitGroup sync.WaitGroup

	for {
		connection, err := listener.Accept()
		if err != nil {
			pl(err)
			continue
		}

	}

}

func setupListener() net.Listener {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter the port you want to listen to:")
		// some ports on windows don't work depending on machine, e.g. 5433. We use 5431 instead.
		address, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			// restart loop until successful
			continue
		}

		// remove delimiter /n or bugs
		address = address[:len(address)-2]

		// TODO: remove localhost later
		address = "localhost:" + address

		listener, err := net.Listen("tcp", address)
		if err != nil {
			fmt.Println(err)
		} else {
			pl("Now listning to Address:", address)
			return listener
		}
	}
}

/*
When your server starts, the first thing that it will need to do is establish a socket
connection that it can use to listen for incoming connections. Your server should listen
on the port specified from the command line and wait for incoming client connections.
Each new client request is accepted, and a new Go routine is spawned to handle the request.
To avoid overwhelming your server, you should not create more than a reasonable number of
child processes (for this assignment, use at most 10). In case an additional child process
would break this limit, your server should wait until one of its ongoing child processes
exits before forking a new one to handle the new request.

Once a client has connected, the server should read data from the client and then check for a
properly-formatted HTTP request. Your server should accept requests for files ending in html,
txt, gif, jpeg, jpg, or css and transmit them to the client with a Content-Type of text/html,
text/plain, image/gif, image/jpeg, image/jpeg, or text/css, respectively. If the client
requests a file with any other extension, the web server must respond with a well-formed 400
"Bad Request" code. An invalid request from the client should be answered with an appropriate
error code, i.e. "Bad Request" (400) or "Not Implemented" (501) for valid HTTP methods other than GET.
If the requested file does not exist, your server should return a well-formed 404 "Not Found" code.
Similarly, if headers are not properly formatted for parsing or any other error condition not listed
before, your server should also generate a type-400 message.  For POST requests, please make sure
that you store the files and make them accessible with a subsequent GET request.
*/
