package main

import (
	"sync"
	"time"
)

type Collection struct {
	lock  *sync.RWMutex
	rooms map[string]*Room
}

type errorCollection struct {
	s string
}

func (e *errorCollection) Error() string {
	return e.s
}
func NewCollection() *Collection {
	return &Collection{new(sync.RWMutex), make(map[string]*Room)}
}

func (c *Collection) Join(id string) (room *Room, err error) {
	//DB驗證在這邊驗證尚未實作
	if id != "test" {
		return nil, &errorCollection{"empt"}
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, ok := c.rooms[id]; !ok {
		c.rooms[id] = NewRoom(id)
		go c.rooms[id].run()
	}
	room = c.rooms[id]

	return
}

func (c *Collection) Get(id string) (room *Room, err error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if val, ok := c.rooms[id]; ok {
		room = val
	} else {
		err = &errorCollection{"no colleciotn"}
	}
	return

}

func (c *Collection) run() {

	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			c.lock.Lock()
			for id, room := range c.rooms {
				if len(room.connections) == 0 {
					delete(c.rooms, id)
				}
			}
			c.lock.Unlock()
		}
	}

}
