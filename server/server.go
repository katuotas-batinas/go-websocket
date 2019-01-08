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
	messageType          = websocket.TextMessage
	searchChar           = "?"
	replaceChar          = "!"
	undefinedRoomMessage = "You must define room name."
	roomNameTakenMessage = "This room name is already taken."
	nonExistentRoomMessage = "Room %s does not exist."
	newSubscriberMessage = "New subscriber connected. You have %d subscribers."
	subscriberDisconnectedMessage = "Subscriber has disconnected. You have %d subscribers."
	publisherDisconnectedMessage = "Publisher of %s has disconnected."
)

var upgrader = websocket.Upgrader{}
var rooms map[string]*Room

func servePublisher(w http.ResponseWriter, r *http.Request) {
	// Check if room parameter is set
	roomName := r.URL.Query().Get("room")
	if roomName == "" {
		fmt.Fprint(w, undefinedRoomMessage)
		return
	}

	// Check if room name is available
	if _, ok := rooms[roomName]; ok {
		fmt.Fprint(w, roomNameTakenMessage)
		return
	}

	// Upgrade HTTP connection to WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Create client nad start I/O goroutines
	client := newClient(conn)
	go client.listen()
	go client.write()

	// Create room
	room := newRoom(client)
	rooms[roomName] = room
	go room.listen()

	for {
		select {
		case msg := <-client.read:
			msg = bytes.Replace(msg, []byte(searchChar), []byte(replaceChar), -1)
			room.broadcast <- msg
		case <-client.disconnect:
			fmt.Printf(publisherDisconnectedMessage+"\n", roomName)
			delete(rooms, roomName)
			return
		case <-room.onSubscribe:
			client.send <-[]byte(fmt.Sprintf(newSubscriberMessage, len(room.subscribers)))
		case <-room.onUnsubscribe:
			client.send <-[]byte(fmt.Sprintf(subscriberDisconnectedMessage, len(room.subscribers)))
		}
	}
}

func serveSubscriber(w http.ResponseWriter, r *http.Request) {
	// Check if room parameter is set
	roomName := r.URL.Query().Get("room")
	if roomName == "" {
		fmt.Fprint(w, undefinedRoomMessage)
		return
	}

	// Check if room exists
	room, ok := rooms[roomName]
	if !ok {
		fmt.Fprintf(w, nonExistentRoomMessage, roomName)
		return
	}

	// Upgrade HTTP connection to WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Create client and start I/O goroutine
	client := newClient(conn)
	go client.listen()
	go client.write()

	room.subscribe <- client

	for {
		select {
		case msg := <-client.read:
			msg = bytes.Replace(msg, []byte(searchChar), []byte(replaceChar), -1)
			room.broadcast <- msg
		case <-client.disconnect:
			room.unsubscribe <- client
			return
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

	http.HandleFunc("/publish", servePublisher)
	http.HandleFunc("/subscribe", serveSubscriber)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
