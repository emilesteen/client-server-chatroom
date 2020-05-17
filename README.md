# client-server-chatroom

Implementation of a client-server chatroom

## Installation
```
go get -u github.com/emilesteen/client-server-chatroom/server
go get -u github.com/emilesteen/client-server-chatroom/client
```

## Usage
To start the server:<br/>
```
cd $GOPATH/src/github.com/emilesteen/client-server-chatroom
go run server/server.go
```

|Option|Description|Default|
|--|--|--|
|-port|Port where the server should listen for new connections|8001|


To start a new client:<br/>
```
go run client/client.go
```
|Option|Description|default|
|--|--|--|
|-ip|IP address of the chat server|127.0.0.1|
|-port|Port where the server is listening for new connections|8001|
