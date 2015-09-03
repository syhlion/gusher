package core

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
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

func (c *Client) ReadPump() {
	defer func() {
		c.ws.Close()
		c.App.Unregister <- c
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := c.ws.ReadMessage()
		if err != nil {
			return
		}
		//暫時不實做推送到 送出頻道 目前是 readonly
		//c.Send <- msg
	}

}
func (c *Client) WritePump() {
	t := time.NewTicker(pingPeriod)
	defer func() {
		log.Debug(c.ws.RemoteAddr().String(), " ", c.Tag, " dissconect")
		c.ws.Close()
		c.App.Unregister <- c
		t.Stop()
	}()
	for {
		select {
		case msg, ok := <-c.Send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				log.Debug(c.ws.RemoteAddr().String(), " ", c.Tag, " Send Channel Error")
				return
			}
			if err := c.write(websocket.TextMessage, msg); err != nil {
				log.Debug(c.ws.RemoteAddr().String(), " ", c.Tag, " Send Message Error", err)
				return
			}

		case <-t.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				log.Debug(c.ws.RemoteAddr().String(), " ", c.Tag, " Send PingMessage Error")
				return
			}

		}
	}

}
