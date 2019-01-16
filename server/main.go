package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	"runtime"

	"github.com/gorilla/websocket"
)

const (
	messageType                  = websocket.TextMessage
	searchChar                   = "?"
	replaceChar                  = "!"
	undefinedRoomMessage         = "You must define room name."
	roomNameTakenMessage         = "This room name is already taken."
	nonExistentRoomMessage       = "Room %s does not exist."
	firstSubscriberMessage       = "FIRST_SUBSCRIBER"
	noSubscribersMessage         = "NO_SUBSCRIBERS"
	newSubscriberMessage         = "New subscriber connected. You have %d subscribers."
	publisherDisconnectedMessage = "PUBLISHER_DISCONNECTED"
)

var mutex = &sync.Mutex{}
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
	mutex.Lock()
	_, ok := rooms[roomName]
	mutex.Unlock()
	if ok {
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
	mutex.Lock()
	rooms[roomName] = room
	mutex.Unlock()
	go room.listen()

	for {
		select {
		case msg := <-client.read:
			msg = bytes.Replace(msg, []byte(searchChar), []byte(replaceChar), -1)
			room.broadcast <- msg
		case <-client.disconnect:
			room.broadcast <- []byte(publisherDisconnectedMessage)
			close(room.stop)
			mutex.Lock()
			delete(rooms, roomName)
			mutex.Unlock()
			return
		case <-room.onSubscribe:
			if len(room.subscribers) == 1 {
				client.send <- []byte(firstSubscriberMessage)
			} else {
				client.send <- []byte(fmt.Sprintf(newSubscriberMessage, len(room.subscribers)))
			}
		case <-room.onUnsubscribe:
			if len(room.subscribers) == 0 {
				client.send <- []byte(noSubscribersMessage)
			}
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
	mutex.Lock()
	room, ok := rooms[roomName]
	mutex.Unlock()
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

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-ticker.C:
				fmt.Println(runtime.NumGoroutine())
			}
		}
	}()

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
