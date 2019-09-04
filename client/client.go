package main

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
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

func StartClient(ip string) error {
	// Open a readWriter that is connected to the server
	log.Println("Opening connection")
	addr := ip + port
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return errors.Wrap(err, "Dialing "+addr+" failed")
	}
	log.Println("Connection open.")

	// Receive a message from the server
	message := receiveMessage(conn)
	fmt.Print(message)

	go sendMessageRoutine(conn)

	// Send a message to the client
	//sendStringMessage(conn, "Acknowledged\n")
	//log.Println("Message sent.")

	message = receiveMessage(conn)
	fmt.Print(message)

	// Close connection
	log.Println("Closing connection...")
	err = conn.Close()
	if err != nil {
		return errors.Wrap(err, "Failed to close connection.")
	}
	log.Println("Connection closed.")

	return nil
}

func main() {
	StartClient("127.0.0.1")
}

