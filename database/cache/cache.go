package cache

import (
	"sync"
)

type Item interface {
	UID() string
}

type Cache struct {
	sync.Mutex
	cache map[string]Item
}

func New() Cache {
	c := Cache{
		cache: make(map[string]Item, 0),
	}
	return c
}

func (c Cache) Check(uid string) (Item, bool) {
	i, ok := c.cache[uid]
	return i, ok
}

func (c *Cache) Add(i Item) {
	c.Lock()
	defer c.Unlock()
	uid := i.UID()
	if _, ok := c.Check(uid); ok {
		return
	}
	c.cache[uid] = i
}

type CacheInt32 struct {
	sync.Mutex
	cache map[string]map[int32]Item
}

func NewCacheInt32() CacheInt32 {
	return CacheInt32{
		cache: make(map[string]map[int32]Item, 0),
	}
}

func (c CacheInt32) Check(id1 string, id2 int32) (Item, bool) {
	_, ok := c.cache[id1]
	if ok {
		l2, ok := c.cache[id1][id2]
		return l2, ok
	}
	return nil, ok
}

func (c *CacheInt32) Add(i Item, id1 string, id2 int32) {
	if len(id1) == 0 || id2 == 0 {
		return
	}
	c.Lock()
	defer c.Unlock()
	if _, ok := c.cache[id1]; !ok {
		c.cache[id1] = make(map[int32]Item, 1)
	}
	c.cache[id1][id2] = i
}

type CacheString struct {
	sync.Mutex
	cache map[string]map[string]Item
}

func NewCacheString() CacheString {
	return CacheString{
		cache: make(map[string]map[string]Item, 0),
	}
}

func (c CacheString) Check(id1 string, id2 string) (Item, bool) {
	_, ok := c.cache[id1]
	if ok {
		l2, ok := c.cache[id1][id2]
		return l2, ok
	}
	return nil, ok
}

func (c *CacheString) Add(i Item, id1 string, id2 string) {
	if len(id1) == 0 || len(id2) == 0 {
		return
	}
	c.Lock()
	defer c.Unlock()
	if _, ok := c.cache[id1]; !ok {
		c.cache[id1] = make(map[string]Item, 1)
	}
	c.cache[id1][id2] = i
}
