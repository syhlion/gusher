package main

import (
	"log"
)

type Room struct {
	id          string
	connections map[*Client]bool
	Boradcast   chan []byte
	Register    chan *Client
	Unregister  chan *Client
}

func NewRoom(id string) *Room {
	return &Room{
		id,
		make(map[*Client]bool),
		make(chan []byte),
		make(chan *Client),
		make(chan *Client),
	}

}

func (r *Room) run() {
	for {
		select {
		case c := <-r.Register:
			r.connections[c] = true
		case c := <-r.Unregister:
			if _, ok := r.connections[c]; ok {
				delete(r.connections, c)
				log.Println("bbb")
				close(c.send)
				log.Println("ccc")
			}
			if len(r.connections) == 0 {
				log.Println("aaa")
				break
			}
		case m := <-r.Boradcast:
			for c := range r.connections {
				select {
				case c.send <- m:
				default:
					close(c.send)
					delete(r.connections, c)
				}
			}
		}
	}
	defer func() {
		close(r.Boradcast)
		close(r.Register)
		close(r.Unregister)
	}()
}
