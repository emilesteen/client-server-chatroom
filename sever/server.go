package main

import (
	"github.com/pkg/errors"
	"log"
	"net"
	"strings"
	"sync"
)

const (
	port = ":8001"
)

var clients = make(map[string]net.Conn)
var lock = sync.RWMutex{}
var buf [512]byte

func startServer() error {
	ln, err := listen()
	if err != nil {
		return errors.Wrap(err, "Unable to start server because of listen error.")
	}
	err = acceptConnections(ln)
	return err
}

func listen() (net.Listener, error) {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to listen on port %s.\n", port)
	}
	log.Println("Listening on port: " + ln.Addr().String())
	return ln, nil
}

func acceptConnections(ln net.Listener) error {
	for {
		conn, err := ln.Accept()
		log.Println("Client connected with address: " + conn.LocalAddr().String())
		if err != nil {
			log.Println("Failed to accept connection request.")
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	log.Println("Handling client...")
	clientName := getClientName(conn)
	echoMessages(conn, clientName)
	closeConnection(conn, clientName)
}

func getClientName(conn net.Conn) (clientName string) {
	sendMessage(conn, "You are connected to the server, choose a username.\n")

	for {
		clientName = receiveMessage(conn)
		clientName = strings.TrimRight(clientName, "\n")
		lock.Lock()
		_, in := clients[clientName]
		if !in {
			clients[clientName] = conn
			lock.Unlock()
			break
		}
		lock.Unlock()
		sendMessage(conn, "The name is already taken, please choose another one.\n")
	}

	sendMessage(conn, "Welcome to the room, "+clientName+"\n")
	lock.RLock()
	for name, conn := range clients {
		if name != clientName {
			sendMessage(conn, clientName+" joined the room.\n")
		}
	}
	lock.RUnlock()

	return
}

func echoMessages(conn net.Conn, clientName string) {
	message := ""
	pre := clientName + ": "
	for {
		message = receiveMessage(conn)
		if message == "!q\n" {
			sendMessage(conn, "!q\n")
			break
		} else {
			broadcastMessage(pre + message)
		}
	}
}

func closeConnection(conn net.Conn, clientName string) {
	log.Println("Closing connection...")
	err := conn.Close()
	if err != nil {
		log.Println("Closing connection failed.")
	}
	log.Println("Connection closed.")

	lock.Lock()
	delete(clients, clientName)
	lock.Unlock()

	broadcastMessage(clientName + " left the room.\n")
}

func sendMessage(conn net.Conn, message string) {
	_, err := conn.Write([]byte(message))
	if err != nil {
		log.Println("Cannot write to connection.")
	}
}

func receiveMessage(conn net.Conn) (message string) {
	n, err := conn.Read(buf[0:])
	if err != nil {
		log.Println("Read error.")
	}
	message = string(buf[0:n])
	return
}

func broadcastMessage(message string) {
	lock.RLock()
	for _, conn := range clients {
		sendMessage(conn, message)
	}
	lock.RUnlock()
}

func main() {
	err := startServer()
	if err != nil {
		log.Println("Error:", errors.WithStack(err))
	}
}
