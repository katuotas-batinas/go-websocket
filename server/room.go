package main

type Room struct {
	publisher *Client

	subscribers map[*Client]*Client

	subscribe chan *Client

	onSubscribe chan *Client

	unsubscribe chan *Client

	onUnsubscribe chan *Client

	broadcast chan []byte
}

func (r *Room) listen() {
	for {
		select {
		case subscriber := <-r.subscribe:
			r.subscribers[subscriber] = subscriber
			r.onSubscribe <- subscriber
		case subscriber := <-r.unsubscribe:
			if _, ok := r.subscribers[subscriber]; ok {
				delete(r.subscribers, subscriber)
				r.onUnsubscribe <- subscriber
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
		publisher:     publisher,
		subscribers:   make(map[*Client]*Client),
		subscribe:     make(chan *Client),
		onSubscribe:   make(chan *Client),
		unsubscribe:   make(chan *Client),
		onUnsubscribe: make(chan *Client),
		broadcast:     make(chan []byte),
	}
}
