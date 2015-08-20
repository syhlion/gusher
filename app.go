package main

import (
	"regexp"
)

type App struct {
	key               string
	Connections       map[*Client]bool
	AssignTotalResult chan int
	Assign            chan map[string][]byte
	Boradcast         chan []byte
	Register          chan *Client
	Unregister        chan *Client
}

func NewApp(app_key string) *App {
	return &App{
		key:               app_key,
		Connections:       make(map[*Client]bool, 1024),
		AssignTotalResult: make(chan int),
		Assign:            make(chan map[string][]byte, 1024),
		Boradcast:         make(chan []byte, 1024),
		Register:          make(chan *Client, 1024),
		Unregister:        make(chan *Client, 1024),
	}
}
func (a *App) run() {
	for {
		select {
		case client := <-a.Register:
			a.Connections[client] = true
		case client := <-a.Unregister:
			if _, ok := a.Connections[client]; ok {
				log.Info(client.ws.RemoteAddr().String(), " ", client.Tag, " disconnect")
				delete(a.Connections, client)
				close(client.Send)
			}
			if len(a.Connections) == 0 {
				log.Info("This Connections is 0 Break this foreach")
				break
			}
		case message := <-a.Boradcast:
			log.Info(a.key, " Boradcast start")
			for client := range a.Connections {
				client.Send <- message
			}
		case ruleMsg := <-a.Assign:
			log.Info(a.key, " Assign Start")
			i := 0
			//迴圈跑所有連線
			for client := range a.Connections {

				//跑規則map
				for rule, message := range ruleMsg {

					//檢查正規式
					if vailed, err := regexp.Compile(rule); err == nil {
						if vailed.MatchString(client.Tag) {

							client.Send <- message
							i++
						}
					}
				}
			}
			a.AssignTotalResult <- i
		}
	}
	defer func() {
		close(a.Boradcast)
		close(a.Register)
		close(a.Unregister)
	}()
}
