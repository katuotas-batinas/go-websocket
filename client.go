package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

func main() {
	flag.Parse()

	// Check if server URL is specified
	if len(flag.Args()) < 1 {
		log.Fatal("Must specify url")
	}

	var url = flag.Args()[0]

	// Connect to the server
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Create input reader
	reader := bufio.NewReader(os.Stdin)

	// Wait for first user input and send it to the server
	fmt.Print("Message: ")
	input, _ := reader.ReadBytes('\n')
	conn.WriteMessage(websocket.TextMessage, bytes.TrimRight(input, "\n"))

	// Websocket read loop
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Printf("Received: %s\n", msg)

		fmt.Print("Message: ")
		input, _ := reader.ReadBytes('\n')
		conn.WriteMessage(websocket.TextMessage, bytes.TrimRight(input, "\n"))
	}
}
