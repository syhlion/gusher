package main

import (
	"github.com/gorilla/websocket"
	"time"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	Tag  string
	ws   *websocket.Conn
	Send chan []byte
	App  *App
}

func NewClient(tag string, ws *websocket.Conn, app *App) *Client {
	return &Client{
		Tag:  tag,
		ws:   ws,
		Send: make(chan []byte, 1024),
		App:  app,
	}
}

func (c *Client) write(msgType int, msg []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(msgType, msg)
}

/* 暫時用不到的功能
func (c *Client) readPump() {
	defer func() {
		c.ws.Close()
		c.App.Unregister <- c
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			return
		}
		c.Send <- msg
	}

}
*/
func (c *Client) WritePump() {
	t := time.NewTicker(pingPeriod)
	defer func() {
		c.ws.Close()
		c.App.Unregister <- c
		t.Stop()
	}()
	for {
		select {
		case msg, ok := <-c.Send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.ws.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}

		case <-t.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}

		}
	}

}
