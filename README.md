# HTTP Server

This is a simple HTTP server written in Go, supporting concurrent and persistent TCP connections. Instead of using the net/http package, this server is built using the net package for stronger control and to understand the underlying mechanics of HTTP. This project is part of the [Codecrafters HTTP Server Challenge](https://codecrafters.io/challenges/http-server).

** Note: You can try out the persistent connection by running the following command, all requests above are made through the same connection: **

```
curl -v --data "...this is the data sent through the request body..." -H "Content-Type: application/octet-stream" http://http-server.sunray4.hackclub.app/files/file_1111 \
     --next -X GET http://http-server.sunray4.hackclub.app/files/file_1111 \
     --next --header "User-Agent: Mozilla/5.0" http://http-server.sunray4.hackclub.app/user-agent
```

## Features

** curl commands are used to test the server **

- write data from response body of request into a file

```
curl -v --data "...this is the data sent through the request body..." -H "Content-Type: application/octet-stream" http://http-server.sunray4.hackclub.app/files/file_123
```

- read data from an existing file

```
curl -v -X GET http://http-server.sunray4.hackclub.app/files/file_123
```

- return user-agent of request

```
curl -v --header "User-Agent: Mozilla/5.0" http://http-server.sunray4.hackclub.app/user-agent
```

- echo string in url path

```
curl -X GET http://http-server.sunray4.hackclub.app/echo/12345
```

** this returns 12345 **

- echo request data with gzip compression

```
curl -v -H "Accept-Encoding: gzip" http://http-server.sunray4.hackclub.app/echo/12345
```

** Note: running this command with curl will give you a warning because terminal can't display the gzip compressed binary data contained in the response body **

- connection close

```

curl -v -H "Connection: close" http://http-server.sunray4.hackclub.app/

```
