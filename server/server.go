package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	messageType = websocket.TextMessage
	searchChar  = "?"
	replaceChar = "!"
)

var upgrader = websocket.Upgrader{}
var rooms map[string]*Room

func servePublisher(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

  // defer func () {
  // fmt.Println("Stop")
  //   conn.Close()
  // }()

  // Create client
	client := newClient(conn)

  // Create room
  room := newRoom(client)
  rooms["test"] = room

  incoming := make(chan []byte)
	go client.read(incoming)

  for {
    select {
    case msg := <-incoming:
        fmt.Println("ATEJO %s", msg)
    }
  }
}

func main() {
	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Fatal("Must specify port")
	}

	var port = flag.Args()[0]

  rooms = make(map[string]*Room)

	http.HandleFunc("/publisher", servePublisher)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
