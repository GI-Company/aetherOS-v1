package server

import (
	"container/list"
	"sync"
	"time"
)

type cacheEntry struct {
	key       string
	value     interface{}
	expiresAt time.Time
	element   *list.Element
}

type PersistentCache struct {
	mu         sync.Mutex
	capacity   int
	ttl        time.Duration
	items      map[string]*cacheEntry
	lru        *list.List
	bus        *BusServer
	snapshotCh chan string
}

func NewPersistentCache(capacity int, ttl time.Duration, bus *BusServer) *PersistentCache {
	c := &PersistentCache{
		capacity:   capacity,
		ttl:        ttl,
		items:      make(map[string]*cacheEntry),
		lru:        list.New(),
		bus:        bus,
		snapshotCh: make(chan string, 64),
	}
	go c.reaper()
	go c.snapshotWorker()
	return c
}

func (c *PersistentCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if en, ok := c.items[key]; ok {
		en.value = value
		en.expiresAt = time.Now().Add(c.ttl)
		c.lru.MoveToFront(en.element)
		return
	}
	if c.lru.Len() >= c.capacity {
		back := c.lru.Back()
		if back != nil {
			be := back.Value.(*cacheEntry)
			delete(c.items, be.key)
			c.lru.Remove(back)
		}
	}
	entry := &cacheEntry{key: key, value: value, expiresAt: time.Now().Add(c.ttl)}
	entry.element = c.lru.PushFront(entry)
	c.items[key] = entry
}

func (c *PersistentCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if en, ok := c.items[key]; ok {
		if !en.expiresAt.IsZero() && time.Now().After(en.expiresAt) {
			c.lru.Remove(en.element)
			delete(c.items, key)
			return nil, false
		}
		c.lru.MoveToFront(en.element)
		return en.value, true
	}
	return nil, false
}

func (c *PersistentCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if en, ok := c.items[key]; ok {
		c.lru.Remove(en.element)
		delete(c.items, key)
	}
}

func (c *PersistentCache) reaper() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		c.mu.Lock()
		for k, e := range c.items {
			if !e.expiresAt.IsZero() && time.Now().After(e.expiresAt) {
				c.lru.Remove(e.element)
				delete(c.items, k)
			}
		}
		c.mu.Unlock()
	}
}

// Request snapshot for a key (will write to VFS via bus)
func (c *PersistentCache) Snapshot(key string) {
	select {
	case c.snapshotCh <- key:
	default:
		// drop if full
	}
}

func (c *PersistentCache) snapshotWorker() {
	for key := range c.snapshotCh {
		val, ok := c.Get(key)
		if !ok {
			continue
		}
		// prepare payload and publish vfs:write
		payload := map[string]interface{}{"path": "/cache/" + key + ".json", "content": val}
		env := &Envelope{Topic: "vfs:write", From: "kernel.cache", Payload: payload, Time: time.Now()}
		c.bus.Publish(env)
	}
}
