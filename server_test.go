package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestServer(t *testing.T) {
	// Create test server
	s := httptest.NewServer(http.HandlerFunc(serveWs))
	defer s.Close()

	// Convert http URL to ws
	url := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer conn.Close()

	// Send message to server and check response
	testResponse("?", "!", conn, t)
	testResponse("hello???", "hello!!!", conn, t)
	testResponse("???-hey!!!", "!!!-hey!!!", conn, t)
	testResponse("? ? ?", "! ! !", conn, t)
}

func testResponse(send string, expect string, conn *websocket.Conn, t *testing.T) {
	err := conn.WriteMessage(websocket.TextMessage, []byte(send))
	if err != nil {
		t.Errorf("Write: %v", err)
	}

	_, msg, err := conn.ReadMessage()
	if err != nil {
		t.Errorf("Read: %v", err)
	}

	if string(msg) != expect {
		t.Errorf("message transformation failed: got %v want %v", string(msg), expect)
	}
}
