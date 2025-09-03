package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}
type lruCache struct {
	Cache    // Remove me after realization.
	mu       sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}
type lruCacheValue struct {
	lruKey Key
	val    interface{}
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	item, ok := lc.items[key]
	if ok {
		lc.queue.MoveToFront(item)
		return item.Value.(lruCacheValue).val, ok
	}
	return nil, false
}

func (lc *lruCache) Set(key Key, value interface{}) bool {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	item, isset := lc.items[key]
	if isset {
		item.Value = lruCacheValue{
			lruKey: key,
			val:    value,
		}
		lc.queue.MoveToFront(item)
	} else {
		if lc.queue.Len() == lc.capacity {
			deleteItem := lc.queue.Back()
			lc.queue.Remove(deleteItem)
			keyDelete := deleteItem.Value.(lruCacheValue).lruKey
			delete(lc.items, keyDelete)
		}
		item = lc.queue.PushFront(lruCacheValue{
			lruKey: key,
			val:    value,
		})
	}
	lc.items[key] = item
	return isset
}

func (lc *lruCache) Clear() {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	lc.queue = NewList()
	lc.items = make(map[Key]*ListItem, lc.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
