package cache

import (
	"sync"
	"time"

	"github.com/miekg/dns"
)

type cacheItem struct {
	msg        *dns.Msg
	expiration time.Time
}

type Cache struct {
	mu    sync.RWMutex
	items map[string]cacheItem
}

func New() *Cache {
	c := &Cache{
		items: make(map[string]cacheItem),
	}
	go c.cleanup()
	return c
}

func (c *Cache) Get(key string) *dns.Msg {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[key]
	if !ok {
		return nil
	}

	if time.Now().After(item.expiration) {
		return nil
	}

	return item.msg.Copy()
}

func (c *Cache) Set(key string, msg *dns.Msg, ttl uint32) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = cacheItem{
		msg:        msg.Copy(),
		expiration: time.Now().Add(time.Duration(ttl) * time.Second),
	}
}

func (c *Cache) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		for key, item := range c.items {
			if time.Now().After(item.expiration) {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}
