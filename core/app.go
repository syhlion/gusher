package core

import (
	"sync"

	"github.com/syhlion/gusher/module/log"
)

type App struct {
	key         string
	lock        *sync.RWMutex
	Connections map[*Client]bool
	Boradcast   chan []byte
	Register    chan *Client
	Unregister  chan *Client
}

func NewApp(app_key string) *App {
	return &App{
		key:         app_key,
		lock:        new(sync.RWMutex),
		Connections: make(map[*Client]bool, 1024),
		Boradcast:   make(chan []byte, 1024),
		Register:    make(chan *Client, 1024),
		Unregister:  make(chan *Client, 1024),
	}
}

func (a *App) GetAllUserTag() []string {

	a.lock.RLock()
	defer a.lock.RUnlock()
	var list []string
	for client := range a.Connections {
		list = append(list, client.Tag)
	}
	return list
}

func (a *App) run() {
	for {
		select {
		case client := <-a.Register:
			a.Connections[client] = true
		case client := <-a.Unregister:
			if _, ok := a.Connections[client]; ok {
				delete(a.Connections, client)
				close(client.Send)
			}
			if len(a.Connections) == 0 {
				break
			}
		case message := <-a.Boradcast:
			log.Logger.Debug(a.key, " Boradcast start")
			for client := range a.Connections {
				client.Send <- message
			}
		}
	}
	defer func() {
		close(a.Boradcast)
		close(a.Register)
		close(a.Unregister)
	}()
}
