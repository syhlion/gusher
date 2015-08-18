package main

type Room struct {
	id          string
	connections map[*Client]bool
	Assign      chan map[string][]byte
	Boradcast   chan []byte
	Register    chan *Client
	Unregister  chan *Client
}

func NewRoom(id string) *Room {
	return &Room{
		id,
		make(map[*Client]bool),
		make(chan map[string][]byte),
		make(chan []byte),
		make(chan *Client),
		make(chan *Client),
	}

}

func (r *Room) run() {
	for {
		select {
		case client := <-r.Register:
			r.connections[client] = true
		case client := <-r.Unregister:
			if _, ok := r.connections[client]; ok {
				delete(r.connections, client)
				close(client.send)
			}
			if len(r.connections) == 0 {
				break
			}
		case message := <-r.Boradcast:
			for client := range r.connections {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(r.connections, client)
				}
			}
		case arr := <-r.Assign:
			for client := range r.connections {
				for rule, message := range arr {

					//TODO尚未實作規則辨識
					if client.token == rule {
						client.send <- message
					}
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
