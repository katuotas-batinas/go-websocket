# WebSocket server
## Installing dependencies
This application requires gorilla/websocket library, to install it run command:

    go get github.com/gorilla/websocket

## Running example
To run example, start the server:

    go run server.go <port>
And the client:

    go run client.go <url>
URL must containt full path to the websocket endpoint, for example:

    go run client.go ws://localhost:3000/ws

## Test
To test server, run command:

    go test server.go server_test.go -v

## Libraries used

 - [https://github.com/gorilla/websocket](https://github.com/gorilla/websocket)
