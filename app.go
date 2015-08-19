package main

import (
	"regexp"
)

type App struct {
	id          string
	Connections map[*Client]bool
	Assign      chan map[string][]byte
	Boradcast   chan []byte
	Register    chan *Client
	Unregister  chan *Client
}

func NewApp(id string) *App {
	return &App{
		id,
		make(map[*Client]bool, 1024),
		make(chan map[string][]byte, 1024),
		make(chan []byte, 1024),
		make(chan *Client, 1024),
		make(chan *Client, 1024),
	}
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
			for client := range a.Connections {
				client.Send <- message
			}
		case ruleMsg := <-a.Assign:

			//迴圈跑所有連線
			for client := range a.Connections {

				//跑規則map
				for rule, message := range ruleMsg {

					//檢查正規式
					if vailed, err := regexp.Compile(rule); err == nil {
						if vailed.MatchString(client.Tag) {

							client.Send <- message
						}
					}
				}
			}
		}
	}
	defer func() {
		close(a.Boradcast)
		close(a.Register)
		close(a.Unregister)
	}()
}
