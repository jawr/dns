package cache

import (
	"sync"
)

type Item interface {
	UID() string
}

type Cache struct {
	mutex sync.Mutex
	cache map[string]Item
}

func New() Cache {
	c := Cache{
		mutex: sync.Mutex{},
		cache: make(map[string]Item, 0),
	}
	return c
}

/* tolower or not tolower */

func (c Cache) Check(uid string) (Item, bool) {
	i, ok := c.cache[uid]
	return i, ok
}

func (c *Cache) Add(i Item) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	uid := i.UID()
	if _, ok := c.Check(uid); ok {
		return
	}
	c.cache[uid] = i
}
