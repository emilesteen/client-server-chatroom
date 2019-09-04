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
		log.Println("Cannot write to connection.")
	}
	message = string(buf[0:n])
	return
}

func sendMessageRoutine(conn net.Conn) {
	clReader := bufio.NewReader(os.Stdin)
	in := ""
	for in != "!quit\n"  {
		in, err :=clReader.ReadString('\n')
		if err != nil {
			log.Println("Command line read error")
		}
		sendStringMessage(conn, in)
		log.Println("Message sent")
	}
}

func closeConnection(conn net.Conn) {
	// Close connection
	log.Println("Closing connection...")
	err := conn.Close()
	if err != nil {
		log.Println("Failed to close connection")
	}
	log.Println("Connection closed.")
}

func Client(ip string) {
	// Open a readWriter that is connected to the server
	log.Println("Opening connection")
	addr := ip + port
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println("Dialing "+addr+" failed")
	}
	log.Println("Connection open.")

	// Receive a message from the server
	message := receiveMessage(conn)
	fmt.Print(message)

	go sendMessageRoutine(conn)

	for message != "!q\n"  {
		message = receiveMessage(conn)
		fmt.Print(message)
	}

	closeConnection(conn)
}

func main() {
	Client("127.0.0.1")
}

