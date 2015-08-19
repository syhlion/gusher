package main

import (
	"sync"
	"time"
)

type Collection struct {
	lock *sync.RWMutex
	apps map[string]*App
}

type errorCollection struct {
	s string
}

func (e *errorCollection) Error() string {
	return e.s
}
func NewCollection() *Collection {
	return &Collection{new(sync.RWMutex), make(map[string]*App)}
}

func (c *Collection) Join(app_key string) (room *App, err error) {

	//app_key := keys[0]
	//user_token := keys[1]
	//DB驗證在這邊驗證尚未實作

	//DB 驗證結束
	if app_key != "test" {
		return nil, &errorCollection{"empty"}
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, ok := c.apps[app_key]; !ok {
		c.apps[app_key] = NewApp(app_key)
		go c.apps[app_key].run()
	}
	room = c.apps[app_key]

	return
}

func (c *Collection) Get(id string) (room *App, err error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if val, ok := c.apps[id]; ok {
		room = val
	} else {
		err = &errorCollection{"no colleciotn"}
	}
	return

}

//定時掃除空的app集合
func (c *Collection) run() {

	ticker := time.NewTicker(10 * time.Minute)
	for {
		select {
		case <-ticker.C:
			c.lock.Lock()
			for id, room := range c.apps {
				if len(room.Connections) == 0 {
					delete(c.apps, id)
				}
			}
			c.lock.Unlock()
		}
	}

}
