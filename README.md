# broker

`broker` implements a simple file opening broker and a corresponding client for testing the broker. 
The server accepts requests to open a file over a unix domain socket.
It opens the file on behalf of the client process and passes the resulting file descriptor back to the client over the unix domain socket.

A simple Makefile is included for building, running, and testing this package.

## Build

### Server

`make build-server`

### Client

`make build-client`

## Running 

### Server
`./bin/server --sockAddr <socket-address>`

`<socket-address>` is any file path that the server may use for creating a unix domain socket.

### Client

`./bin/client --sockAddr <socket-address> -fname <filename>`

`<socket-address>` is a path to the server's open unix domain socket.

`<fname>` is the file path that the client is requesting to open. 
Its contents will be echoed to the command line.

