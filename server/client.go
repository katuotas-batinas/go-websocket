package main

import (
	"log"
	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn

	send chan []byte
}

func (c *Client) read(incoming chan []byte) {
	defer func() {
    fmt.Println("bye")
		c.conn.Close()
	}()

	for {
		mt, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		// End connection if client message type is not acceptable
		if mt != messageType {
			break
		}

		incoming <- msg
	}
}

func (c *Client) write() {

}

func newClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
		send: make(chan []byte),
	}
}
