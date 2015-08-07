package main

import (
	"github.com/gorilla/websocket"
	"time"
)

const (
	pingWait       = 60 * time.Second
	writeWait      = 10 * time.Second
	readWait       = 10 * time.Second
	maxMessageSize = 512
)

type Client struct {
	ws   *websocket.Conn
	uid  string
	send chan Message
	room *Room
}

func (c *Client) Write(msgType int, msg []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(msgType, msg)
}
func (c *Client) readPump() {
	defer func() {
		c.ws.Close()
		c.room.Unregister <- c
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(readWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(readWait)); return nil })
	for {
		m := Message{}
		err := c.ws.ReadJSON(&m)
		if err != nil {
			return
		}
	}

}

func (c *Client) writePump() {
	t := time.NewTimer(pingWait)
	defer func() {
		c.ws.Close()
		c.room.Unregister <- c
		t.Stop()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				c.Write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.ws.WriteJSON(msg); err != nil {
				return
			}

		case <-t.C:
			if err := c.Write(websocket.PingMessage, []byte{}); err != nil {
				return
			}

		}
	}

}
