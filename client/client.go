package main

import (
	"github.com/marcusolsson/tui-go"
	"log"
	"net"
)

const (
	port = ":8001"
)

var buf [512]byte

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
	messageArea := tui.NewVBox()
	messageAreaScroll := tui.NewScrollArea(messageArea)
	messageAreaScroll.SetAutoscrollToBottom(true)
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
		err := sendMessage(conn, message)
		if err != nil {
			message = "You are not connected to the server, you can not send messages."
			messageArea.Append(tui.NewHBox(tui.NewLabel(message), tui.NewSpacer()))
		}
		input.SetText("")
	})

	root := tui.NewHBox(chat)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() {
		ui.Quit()
		err = sendMessage(conn, "!q\n")
	})

	return ui, messageArea
}

func uiReceiveMessagesRoutine(conn net.Conn, ui tui.UI, messageArea *tui.Box) {
	for {
		message, err := receiveMessage(conn)
		if err != nil {
			message = "You disconnected from the server, you are no longer receiving messages."
			ui.Update(func() {
				messageArea.Append(tui.NewHBox(tui.NewLabel(message), tui.NewSpacer()))
			})
			return
		}
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
		log.Fatal("Dialing " + addr + " failed.")
	}
	log.Println("Connection open.")
	return
}

func sendMessage(conn net.Conn, message string) error {
	_, err := conn.Write([]byte(message))
	return err
}

func receiveMessage(conn net.Conn) (message string, err error) {
	n, err := conn.Read(buf[0:])
	if err != nil {
		return
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
