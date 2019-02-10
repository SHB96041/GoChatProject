package main

import "chapter/project/github.com/gorilla/websocket"

// client는 한 명의 채팅 사용자를 나타내는 구조체 입니다.
type client struct {
		//socket은 이 클라이언트의 웹 소켓입니다.
		socket *websocket.Conn
		// send는 메시지가 전송되는 채널입니다.
		send chan []byte
		// room은 클라이언트가 채팅하는 방입니다.
		room *room
}


// 사용자의 메시지를 읽는 함수입니다.
func (c *client) read() {
	defer c.socket.Close()

	for { // 무한 루프로 채널 내의 메시지를 읽어옵니다.
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

// 사용자가 메시지를 작성하는 함수입니다.
func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}