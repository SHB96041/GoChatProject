package main

import (
	"chapter/project/trace"
	"log"
	"net/http"

	"chapter/project/github.com/gorilla/websocket"
)


type room struct {
	//forward는 수신 메시지를 보관하는 채널입니다.
	forward chan []byte
	// join은 방에 들어오려는 클라이언트를 위한 채널입니다.
	join chan *client
	// leave는 방을 나가길 원하는 클라이언트를 위한 채널입니다.
	leave chan *client
	// clients는 현재 채팅방에 있는 모든 클라이언트를 나타냅니다.
	clients map[*client]bool
	// tracer는 방 안에서 활동의 추적 정보를 수신한다.
	tracer trace.Tracer
}

// 새로운 방을 만드는 함수입니다.
func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer: trace.Off(),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// 입장
			r.clients[client] = true
			r.tracer.Trace("New Client Joined")
		case client := <-r.leave:
			// 퇴장
			delete(r.clients,client)
			close(client.send)
			r.tracer.Trace("Client Left")
		case msg := <-r.forward:
			r.tracer.Trace("Message Received: ", string(msg))
			// 모든 클라이언트에게 메시지 전달
			for client := range r.clients {
				client.send <- msg
				r.tracer.Trace("-- Sent to Client")
			}
		}
	}
}

//상수 정의
const (
	socketBufferSize = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize,
WriteBufferSize: socketBufferSize}
func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServerHTTP:",err)
		return
	}
	client := &client {
		socket: socket,
		send: make(chan []byte, messageBufferSize),
		room: r,
	}
	r.join <- client
	defer func() { r.leave <- client}()
	go client.write()
	client.read()
}