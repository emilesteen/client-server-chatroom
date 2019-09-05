package main

import (
	"github.com/marcusolsson/tui-go"
	"log"
	"net"
	"os"
	"sync"
)

const (
	port = ":8001"
)

var buf [512]byte
var lock sync.RWMutex

func startClientUI(ip string) {
	conn := openConnection(ip)
	ui, messageArea := initUI(conn)

	go uiReceiveMessagesRoutine(conn, ui, messageArea)
	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
	closeConnection(conn)
}

func initUI(conn net.Conn) (tui.UI, *tui.Box) {
	userList := tui.NewVBox(tui.NewLabel("User list:"), tui.NewSpacer())
	userList.SetBorder(true)

	messageArea := tui.NewVBox()
	messageAreaScroll := tui.NewScrollArea(messageArea)
	messageAreaBox := tui.NewVBox(messageAreaScroll)
	messageAreaBox.SetBorder(true)

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	chat := tui.NewVBox(messageAreaBox, inputBox)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)

	input.OnSubmit(func(e *tui.Entry) {
		message := e.Text()
		sendMessage(conn, message)
		input.SetText("")
	})

	root := tui.NewHBox(userList, chat)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() {
		ui.Quit()
		sendMessage(conn, "!q\n")
	})

	return ui, messageArea
}

func uiReceiveMessagesRoutine(conn net.Conn, ui tui.UI, messageArea *tui.Box) {
	for {
		message := receiveMessage(conn)
		if message == "!q\n" {
			break
		}
		ui.Update(func() {
			messageArea.Append(tui.NewHBox(tui.NewLabel(message), tui.NewSpacer()))
		})
	}
}

func openConnection(ip string) (conn net.Conn) {
	log.Println("Opening connection...")
	addr := ip + port
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println("Dialing " + addr + " failed.")
	}
	log.Println("Connection open.")
	return
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
		log.Println("Cannot read from connection.")
		os.Exit(5)
	}
	message = string(buf[0:n])
	return
}

func closeConnection(conn net.Conn) {
	log.Println("Closing connection...")
	err := conn.Close()
	if err != nil {
		log.Println("Failed to close connection.")
	}
	log.Println("Connection closed.")
}

func main() {
	startClientUI("127.0.0.1")
}
