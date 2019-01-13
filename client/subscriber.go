package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gorilla/websocket"
)

const publisherDisconnectedMessage = "PUBLISHER_DISCONNECTED"

func main() {
	flag.Parse()

	// Check if server URL is specified
	if len(flag.Args()) < 1 {
		log.Fatal("Must specify url")
	}

	var url = flag.Args()[0]

	// Connect to the server
	conn, resp, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Fatal(string(body))
	}

	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		if string(msg) == publisherDisconnectedMessage {
			fmt.Println("Publisher has disconnected.")
			return
		} else {
			fmt.Println("Incoming message:", string(msg))
		}
	}
}
