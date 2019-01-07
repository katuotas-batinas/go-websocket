package main

import (
	"bytes"
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
	publishersMax = 1
	subscribersMax = 0
)

var upgrader = websocket.Upgrader{}
var publishers map[*websocket.Conn]bool

func serveWs(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Add connection to publishers map
	publishers[conn] = true

	log.Printf("Publishers connected: %d", len(publishers))

	defer conn.Close()
	defer delete(publishers, conn)

	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		// End connection if client message type is not acceptable
		if mt != messageType {
			break
		}

		log.Printf("Received: %s", msg)

		// Transform message and send back to client
		msg = bytes.Replace(msg, []byte(searchChar), []byte(replaceChar), -1)
		err = conn.WriteMessage(messageType, msg)
		if err != nil {
			log.Println(err)
			break
		}

		log.Printf("Sent: %s", msg)
	}
}

func main() {
	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Fatal("Must specify port")
	}

	var port = flag.Args()[0]

	http.HandleFunc("/publisher", servePublisher)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
