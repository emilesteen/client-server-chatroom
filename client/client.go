package main

import (
	"bufio"
	"github.com/pkg/errors"
	"log"
	"net"
)

const (
	port = ":8001"
)

func sendStringMessage(connRW *bufio.ReadWriter, message string) {
	_, err := connRW.WriteString(message)
	if err != nil {
		log.Println("Cannot write to connection.")
	}
	// Flush the read writer
	err = connRW.Flush()
	if err != nil {
		log.Println("Failed to flush print writer")
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
	connRW := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	log.Println("Connection open.")

	// Receive a message from the server
	message, err := connRW.ReadString('\n')
	if err != nil {
		log.Println("Cannot read from connection.")
	}
	print(message)
	// Send a message to the client
	sendStringMessage(connRW, "Acknowledged\n")
	log.Println("Message sent.")

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

