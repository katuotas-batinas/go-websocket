# WebSocket server
## Installing dependencies
This application requires gorilla/websocket library, to install it run command:

    go get github.com/gorilla/websocket

## Running example
To run example, start the server:

    go run server/*.go <port>
Publisher client:

    go run client/publisher.go ws://<ip>:<port>/publish?room=<name>
Publisher client:

    go run client/subscriber.go ws://<ip>:<port>/subscribe?room=<name>
Example:

    go run server/*.go 3000
    go run client/publisher.go ws://localhost:3000/publish?room=test
    go run client/subscriber.go ws://localhost:3000/subscribe?room=test


## Libraries used

 - [https://github.com/gorilla/websocket](https://github.com/gorilla/websocket)
