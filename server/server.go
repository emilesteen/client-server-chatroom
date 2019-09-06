package main

import (
	"fmt"
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
		return fmt.Errorf("\nunable to start server: %v\n", err)
	}
	err = acceptConnections(ln)
	return err
}

func listen() (net.Listener, error) {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		return nil, fmt.Errorf("\nunable to listen on port: %s because of listen error: %v", port, err)
	}
	log.Println("Listening on address: " + ln.Addr().String())
	return ln, nil
}

func acceptConnections(ln net.Listener) error {
	for {
		conn, err := ln.Accept()
		log.Println("Client connected with address: " + conn.RemoteAddr().String())
		if err != nil {
			log.Println("Failed to accept connection request.")
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	log.Println("Handling client...")

	clientName, err := getClientName(conn)
	if err != nil {
		log.Println("Error with connection: + " + conn.RemoteAddr().String() + ": " + err.Error())
		return
	}
	if clientName == "!q\n" {
		sendMessage(conn, "!q\n")
		log.Println("Closing connection...")
		err := conn.Close()
		if err != nil {
			log.Println("Closing connection failed.")
		}
		log.Println("Connection with address: " + conn.RemoteAddr().String() + " closed.")
		return
	}

	err = echoMessages(conn, clientName)
	if err != nil {
		log.Println("Error with connection: + " + conn.RemoteAddr().String() + ": " + err.Error())
		return
	}

	closeConnection(conn, clientName)
}

func getClientName(conn net.Conn) (string, error) {
	clientName := ""
	sendMessage(conn, "You are connected to the server, choose a username. Press esc to quit.")

	for {
		receivedName, err := receiveMessage(conn)
		if err != err {
			return "", err
		}
		if receivedName == "!q\n" {
			return "!q\n", nil
		}
		receivedName = strings.TrimRight(receivedName, "\n")
		lock.Lock()
		_, in := clients[receivedName]
		if !in {
			clients[receivedName] = conn
			lock.Unlock()
			clientName = receivedName
			break
		}
		lock.Unlock()
		sendMessage(conn, "The name is already taken, please choose another one.")
	}

	sendMessage(conn, "Welcome to the room, "+clientName)
	lock.RLock()
	for name, conn := range clients {
		if name != clientName {
			sendMessage(conn, clientName+" joined the room.")
		}
	}
	lock.RUnlock()

	return clientName, nil
}

func echoMessages(conn net.Conn, clientName string) error {
	pre := clientName + ": "
	for {
		message, err := receiveMessage(conn)
		if err != nil {
			lock.Lock()
			delete(clients, clientName)
			lock.Unlock()

			broadcastMessage(clientName + " left the room.")
			return err
		}
		if message == "!q\n" {
			sendMessage(conn, "!q\n")
			return nil
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
	log.Println("Connection with address: " + conn.RemoteAddr().String() + " closed.")

	lock.Lock()
	delete(clients, clientName)
	lock.Unlock()

	broadcastMessage(clientName + " left the room.")
}

func broadcastMessage(message string) {
	lock.RLock()
	for _, conn := range clients {
		sendMessage(conn, message)
	}
	lock.RUnlock()
}

func sendMessage(conn net.Conn, message string) {
	_, err := conn.Write([]byte(message))
	if err != nil {
		log.Println("Cannot write to connection.")
	}
}

func receiveMessage(conn net.Conn) (message string, err error) {
	n, err := conn.Read(buf[0:])
	if err != nil {
		return
	}
	message = string(buf[0:n])
	return
}

func main() {
	err := startServer()
	if err != nil {
		log.Fatal(err)
	}
}
