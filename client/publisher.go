package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	firstSubscriberMessage = "FIRST_SUBSCRIBER"
	noSubscribersMessage   = "NO_SUBSCRIBERS"
	broadcastMessage       = "Welcome?"
)

func read(conn *websocket.Conn, incoming chan<- []byte) {
	defer func() {
		conn.Close()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		incoming <- msg
	}
}

func broadcast(conn *websocket.Conn, stop chan bool) {
	fmt.Println("Start broadcast")
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			conn.WriteMessage(websocket.TextMessage, []byte(broadcastMessage))
		case <-stop:
			fmt.Println("End broadcast")
			ticker.Stop()
			return
		}
	}
}

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

	incoming := make(chan []byte)
	stopBroadcast := make(chan bool)
	go read(conn, incoming)

	for {
		select {
		case msg := <-incoming:
			switch string(msg) {
			case firstSubscriberMessage:
				go broadcast(conn, stopBroadcast)
			case noSubscribersMessage:
				stopBroadcast <- true
			}
		}
	}
}