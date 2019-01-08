package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn

	read chan []byte

	send chan []byte

	disconnect chan bool
}

func (c *Client) listen() {
	defer func() {
    c.disconnect <- true
		c.conn.Close()
	}()

	for {
		mt, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println(err)
			}
			break
		}

		// End connection if client message type is not acceptable
		if mt != messageType {
			break
		}

		c.read <- msg
	}
}

func (c *Client) write() {
	defer func() {
    c.disconnect <- true
		c.conn.Close()
	}()

	for {
		select {
		case msg := <-c.send:
			err := c.conn.WriteMessage(messageType, msg)
			if err != nil {
				log.Println(err)
				break
			}
		case <-c.disconnect:
			return
		}
	}
}

func newClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
		read: make(chan []byte),
		send: make(chan []byte),
		disconnect: make(chan bool),
	}
}
