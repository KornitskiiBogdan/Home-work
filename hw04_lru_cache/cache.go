package hw04lrucache

import (
	"hash/crc32"
	"sync"
)

const defaultShardCount = 10

type Key string

type cacheItem struct {
	Key
	Value interface{}
}

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type shard struct {
	cache Cache
	mu    sync.Mutex
}

type shardedCache struct {
	shards []*shard
}

func NewCache(capacity int) Cache {
	shardCount := defaultShardCount
	if capacity < defaultShardCount {
		shardCount = capacity
	}

	shards := make([]*shard, shardCount)

	baseCapacity := capacity / shardCount
	extraCapacity := capacity % shardCount

	for i := 0; i < shardCount; i++ {
		capPerShard := baseCapacity
		if i < extraCapacity {
			capPerShard++
		}
		shards[i] = &shard{
			cache: NewLruCache(capPerShard),
		}
	}

	return &shardedCache{shards: shards}
}

func (c *shardedCache) Set(key Key, value interface{}) bool {
	s := c.shardFor(key)
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.cache.Set(key, value)
}
func (c *shardedCache) Get(key Key) (interface{}, bool) {
	s := c.shardFor(key)
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.cache.Get(key)
}

func (c *shardedCache) Clear() {
	for i := range c.shards {
		c.shards[i].mu.Lock()
		c.shards[i].cache.Clear()
		c.shards[i].mu.Unlock()
	}
}

func (c *shardedCache) shardFor(key Key) *shard {
	hash := crc32.ChecksumIEEE([]byte(key))

	index := hash % uint32(len(c.shards))
	return c.shards[index]
}

func NewLruCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	if item, ok := c.items[key]; ok {
		item.Value.(*cacheItem).Value = value
		c.queue.MoveToFront(item)
		return true
	}
	newItem := c.queue.PushFront(&cacheItem{
		Key:   key,
		Value: value,
	})
	c.items[key] = newItem

	if c.queue.Len() > c.capacity {
		backItem := c.queue.Back()
		excluded := backItem.Value.(*cacheItem)
		c.queue.Remove(backItem)
		delete(c.items, excluded.Key)
	}
	return false

}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)
		return item.Value.(*cacheItem).Value, ok
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.items = make(map[Key]*ListItem, c.capacity)
	c.queue = NewList()
}
