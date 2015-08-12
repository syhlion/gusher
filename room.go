package main

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
		make(chan *Client)}
}

func (r *Room) run() {
	for {
		select {
		case c := <-r.Register:
			r.connections[c] = true
		case c := <-r.Unregister:
			if _, ok := r.connections[c]; ok {
				delete(r.connections, c)
				close(c.send)
			}
			if len(r.connections) == 0 {
				break
			}
		}
	}
}
