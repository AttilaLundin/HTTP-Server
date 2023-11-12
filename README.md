# **HTTP Web Server and Proxy Server**

In order to start both the HTTP and Proxy server you need to run main. 
You will be prompted to enter the ip address and port for the web server and proxy server.
We have lazily implemented it so that both are started by goroutines immediately after each other: 
sometimes you are prompted for the web server first and sometimes for the proxy server first, 
the prompt that appears first is the one that will also read your input first! 
We generally use port 5431 for the web server and port 5430 for the proxy server.
However, in the cloud we have containerized the web server and proxy server separately using Docker so that they run in separate instances.


## **Testing**
There are some tests that you can start by uncommenting the following line in row 26 in the connectionhhandler:

`//go test.INITTEST()`

It is important that you have the correct repository structure if you wish to use these tests or properly GET/POST in general(!)

* The storage folder is initially empty, this will result in GET tests failing as intended.
* Then we POST all file types 50 times in a for loop with debug statements printing success/fail statements. 
* Finally we send GET requests for all the different types and store the files in the clienttest folder. 

## **The Cloud and Docker**

We have moved our solution for the web server and proxy server to the cloud using the student granted AWS accounts. 
We containerized our web server and proxy server using Docker and uploaded them to a public repo in Dockerhub.
From our cloud server, we pulled the two containers and can now run them separately on the cloud with communication being possible between the
web server and proxy server.
#### **NOTE that the containers each expose port 5431**

* Dockerhub web server: https://hub.docker.com/r/attilalundin/web-server
* Dockerhub proxy server: https://hub.docker.com/r/attilalundin/proxy-server

**To run the code in the cloud using Docker, first pull the web server and proxy server from the cloud server using these two commands:**

    sudo docker pull attilalundin/web-server
    sudo docker pull attilalundin/proxy-server


**make sure you have the images called attilalundin/web-server:golang, attilalundin/web-server:latest, attilalundin/proxy-server:latest.**
        
    sudo docker images

**then run Docker containers using the following commands:**

    sudo docker run -d -p 5431:5431 -e IP=0.0.0.0 -e PORT=5431 attilalundin/web-server:latest
    sudo docker run -d -p 5430:5431 -e IP=0.0.0.0 -e PORT=5431 attilalundin/proxy-server:latest

**in order to send requests through the terminal you can use the following commands:**
    
    GET request 
    curl -v http://<WEB IPADRESS>:<WEB PORT>/web-server/storage/testfile.txt --output -

    GET  request through proxy
    curl -v -X GET http://<WEB IPADRESS>:<WEB PORT>/web-server/storage/testfile.txt -x  <PROXY IP>:<PROXY PORT> --output -

    POST request
    curl -v -X POST -F "file=@<absolute path>" http://<WEB IPADRESS>:<WEB PORT>/web-server/storage/

    POST request through proxy
    curl -v -X POST -F "file=@<absolute path>" http://<WEB IPADRESS>:<WEB PORT>/web-server/storage/ -x  <PROXY IP>:<PROXY PORT> --output -