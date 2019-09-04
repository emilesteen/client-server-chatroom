package main

import (
	"fmt"
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

func listen() (net.Listener, error) {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to listen on port %s\n", port)
	}
	log.Println("Listening on port: " + ln.Addr().String() + "\n")
	return ln, nil
}

func acceptConnections(ln net.Listener) error {
	for {
		conn, err := ln.Accept()
		log.Println("Accepting a connection request.")
		//log.Println("Client connected with address: " + conn.LocalAddr().String())
		if err != nil {
			errors.Wrap(err, "Failed to accept connection request")
			continue
		}
		go handleClient(conn)
	}
}

func getClientName(conn net.Conn) (clientName string) {
	// Send a message to the client
	sendMessage(conn, "You are connected to the server, choose a username.\n")
	log.Println("Message sent.")

	for {
		// Receive a message from the client
		clientName = receiveMessage(conn)
		clientName = strings.TrimRight(clientName, "\n")
		fmt.Print(clientName)
		lock.Lock()
		_, in := clients[clientName]
		if !in {
			clients[clientName] = conn
			lock.Unlock()
			break
		}
		lock.Unlock()
		sendMessage(conn, "The name is already taken, please choose another one.\n")
		clientName = receiveMessage(conn)
	}

	sendMessage(conn, "Welcome to the room, "+clientName+"\n")
	return
}

func closeConnection(conn net.Conn, clientName string) {
	// Close the connection
	log.Println("Closing connection...")
	err := conn.Close()
	if err != nil {
		log.Println("Closing connection failed.")
	}
	log.Println("Connection closed")
	lock.Lock()
	delete(clients, clientName)
	lock.Unlock()
}

func echoMessage(conn net.Conn, clientName string) {
	message := ""
	pre := clientName + ": "
	for message != "!q\n" {
		message = receiveMessage(conn)
		sendMessage(conn, pre+message)
	}
}

func handleClient(conn net.Conn) {
	log.Println("Handling client...")

	clientName := getClientName(conn)

	echoMessage(conn, clientName)

	closeConnection(conn, clientName)
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

func StartServer() error {
	ln, err := listen()
	if err != nil {
		return errors.Wrap(err, "Unable to start server because of listen error")
	}
	err = acceptConnections(ln)
	return err
}

func main() {
	println("*-START-*\n")
	err := StartServer()
	if err != nil {
		log.Println("Error:", errors.WithStack(err))
	}
	println("*-END-*")
}
