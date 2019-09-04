package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const (
	port = ":8001"
)

var buf [512]byte

func sendStringMessage(conn net.Conn, message string) {
	_, err := conn.Write([]byte(message))
	if err != nil {
		log.Println("Cannot write to connection.")
	}
}

func receiveMessage(conn net.Conn) (message string) {
	n, err := conn.Read(buf[0:])
	if err != nil {
		log.Println("Cannot read from connection.")
		os.Exit(5)
	}
	message = string(buf[0:n])
	return
}

func sendMessageRoutine(conn net.Conn) {
	clReader := bufio.NewReader(os.Stdin)
	for {
		message, err := clReader.ReadString('\n')
		if err != nil {
			log.Println("Command line read error")
		}

		sendStringMessage(conn, message)

		if message == "!q\n" {
			break
		}
	}
}

func receiveMessageRoutine(conn net.Conn) {
	for {
		message := receiveMessage(conn)
		if message == "!q\n" {
			break
		}
		fmt.Print(message)
	}
}

func openConnection(ip string) (conn net.Conn) {
	log.Println("Opening connection...")
	addr := ip + port
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println("Dialing " + addr + " failed")
	}
	log.Println("Connection open.")
	return
}

func closeConnection(conn net.Conn) {
	log.Println("Closing connection...")
	err := conn.Close()
	if err != nil {
		log.Println("Failed to close connection")
	}
	log.Println("Connection closed.")
}

func client(ip string) {
	conn := openConnection(ip)

	// Receive welcome message
	message := receiveMessage(conn)
	fmt.Print(message)

	go sendMessageRoutine(conn)
	receiveMessageRoutine(conn)
	closeConnection(conn)
}

func main() {
	client("127.0.0.1")
}
