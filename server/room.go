package main

import "fmt"

type Room struct {
	publisher *Client

	subscribers map[*Client]*Client

	subscribe chan *Client

	unsubscribe chan *Client

	broadcast chan []byte
}

func (r *Room) listen() {
	for {
		select {
		case subscriber := <-r.subscribe:
			r.subscribers[subscriber] = subscriber
			fmt.Println("New subsriber connected")
		case subscriber := <-r.unsubscribe:
			if _, ok := r.subscribers[subscriber]; ok {
				delete(r.subscribers, subscriber)
				fmt.Println("Subscriber disconnected")
			}
		case msg := <-r.broadcast:
			for _, subscriber := range r.subscribers {
			    subscriber.send <- msg
			}
		}
	}
}

func newRoom(publisher *Client) *Room {
	return &Room{
		publisher: publisher,
		subscribers:      make(map[*Client]*Client),
		subscribe:    make(chan *Client),
		unsubscribe: make(chan *Client),
	}
}
