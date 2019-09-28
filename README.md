# client-server-chatroom

Basic client server Terminal UI chatroom app

## Usage
```
go get -u github.com/emilesteen/client-server-chatroom/server
go get -u github.com/emilesteen/client-server-chatroom/client
cd ~/go/src/github.com/emilesteen/client-server-chatroom
```

To start the server:<br/>
```
go run server/server.go
```

To start a new client:<br/>
```
go run client/client.go [ip address]
```
[ip address] -> server address, if no argument is given, the address defaults to localhost
