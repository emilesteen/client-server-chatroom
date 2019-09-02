package main

import (
	"bufio"
	"github.com/pkg/errors"
	"log"
	"net"
)

const port = ":8001"

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
		log.Println("Handle client connection")
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	log.Println("Handling client.")
	connRW := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	// Send a message to the client
	log.Println("Sending message.")
	sendStringMessage(connRW, "You are connected to the server.\n")
	log.Print("Message sent.")

	// Receive a message from the client
	log.Println("Receiving message.")
	message, err := connRW.ReadString('\n')
	if err != nil {
		log.Println("Cannot read from connection.")
	}
	log.Println("Message received")
	log.Println(message)

	// Close the connection
	err = conn.Close()
	if err != nil {
		log.Println("Closing connection failed.")
	}
}

func sendStringMessage(connRW *bufio.ReadWriter, message string) {
	_, err := connRW.WriteString(message)
	if err != nil {
		log.Println("Cannot write to connection.")
	}
	err = connRW.Flush()
	if err != nil {
		log.Println("Flushing connection read writer failed.")
	}
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