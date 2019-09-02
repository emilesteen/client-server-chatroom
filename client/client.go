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

func openConnection(addr string) (*bufio.ReadWriter, error) {
	log.Println("Dial " + addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, errors.Wrap(err, "Dialing "+addr+" failed")
	}
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}

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
	addr := ip + port
	log.Println("Dialing address: " + addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return errors.Wrap(err, "Dialing "+addr+" failed")
	}
	log.Println("Connected to server.")
	connRW := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	// Receive a message from the server
	log.Println("Receiving message.")
	message, err := connRW.ReadString('\n')
	if err != nil {
		log.Println("Cannot read from connection.")
	}
	log.Println("Message received.")
	log.Print(message)
	// Send a message to the client
	log.Println("Sending message.")
	sendStringMessage(connRW, "Acknowledged\n")
	log.Println("Message sent.")

	// Close connection
	err = conn.Close()
	if err != nil {
		return errors.Wrap(err, "Failed to close connection.")
	}

	return nil
}

func main() {
	StartClient("127.0.0.1")
}

