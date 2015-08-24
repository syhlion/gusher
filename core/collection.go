package core

import (
	"github.com/syhlion/gopusher/model"
	"github.com/syhlion/gopusher/module/log"
	"sync"
	"time"
)

var (
	Collection *collection = nil
)

func init() {
	if Collection == nil {
		Collection = newCollection()
	}
}

type collection struct {
	lock *sync.RWMutex
	apps map[string]*App
}

type errorCollection struct {
	s string
}

func (e *errorCollection) Error() string {
	return e.s
}
func newCollection() *collection {
	return &collection{new(sync.RWMutex), make(map[string]*App)}
}

func (c *collection) Join(app_key string) (room *App, err error) {
	if !model.AppData.IsExist(app_key) {
		err = &errorCollection{"app_key no exist"}
		return
	}
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

func (c *collection) Get(app_key string) (room *App, err error) {
	if !model.AppData.IsExist(app_key) {
		err = &errorCollection{"Error app_key, Please Register App_key"}
		log.Logger.Debug(app_key, " ", err)
		return
	}
	c.lock.RLock()
	defer c.lock.RUnlock()
	if val, ok := c.apps[app_key]; ok {
		room = val
	} else {
		err = &errorCollection{"No User In the App"}
		log.Logger.Debug(app_key, " ", err)
	}
	return

}

//定時掃除空的app集合
func (c *collection) Run() {

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
