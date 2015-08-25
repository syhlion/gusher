package core

import (
	"github.com/syhlion/gopusher/module/log"
	"sync"
	"time"
)

type Collection struct {
	lock *sync.RWMutex
	apps map[string]*App
}

type ErrorCollection struct {
	s string
}

func (e *ErrorCollection) Error() string {
	return e.s
}
func NewCollection() *Collection {
	return &Collection{new(sync.RWMutex), make(map[string]*App)}
}

func (c *Collection) Join(app_key string) (room *App, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, ok := c.apps[app_key]; !ok {
		c.apps[app_key] = NewApp(app_key)
		go c.apps[app_key].run()
	}
	log.Logger.Debug("Join ", app_key)
	room = c.apps[app_key]

	return
}

func (c *Collection) Get(app_key string) (room *App, err error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if val, ok := c.apps[app_key]; ok {
		room = val
	} else {
		err = &ErrorCollection{"No User In the App"}
		log.Logger.Debug(app_key, " ", err)
	}
	return

}

//定時掃除空的app集合
func (c *Collection) Run() {

	ticker := time.NewTicker(10 * time.Minute)
	for {
		select {
		case <-ticker.C:
			c.lock.Lock()
			for app_key, app := range c.apps {
				if len(app.Connections) == 0 {
					log.Logger.Debug("clear empty app", app_key)
					delete(c.apps, app_key)
				}
			}
			c.lock.Unlock()
		}
	}

}
