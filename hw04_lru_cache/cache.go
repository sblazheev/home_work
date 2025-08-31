package hw04lrucache

import (
	"sync"
)

type Key string

var mu sync.Mutex

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}
type lruCache struct {
	Cache // Remove me after realization.

	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	mu.Lock()
	defer mu.Unlock()

	item, ok := lc.items[key]
	if ok {
		lc.queue.MoveToFront(item)
		return item.Value.(struct {
			lruKey Key
			val    interface{}
		}).val, ok
	}
	return nil, false
}

func (lc *lruCache) Set(key Key, value interface{}) bool {
	mu.Lock()
	defer mu.Unlock()

	item, isset := lc.items[key]
	if isset {
		item.Value = struct {
			lruKey Key
			val    interface{}
		}{
			lruKey: key,
			val:    value,
		}
		lc.queue.MoveToFront(item)
	} else {
		if lc.queue.Len() == lc.capacity {
			deleteItem := lc.queue.Back()
			lc.queue.Remove(deleteItem)
			keyDelete := deleteItem.Value.(struct {
				lruKey Key
				val    interface{}
			}).lruKey
			delete(lc.items, keyDelete)
		}
		item = lc.queue.PushFront(struct {
			lruKey Key
			val    interface{}
		}{
			lruKey: key,
			val:    value,
		})
	}
	lc.items[key] = item
	return isset
}

func (lc *lruCache) Clear() {
	mu.Lock()
	defer mu.Unlock()

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
