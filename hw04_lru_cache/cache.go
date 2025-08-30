package hw04lrucache

type Key string

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
	keys     map[*ListItem]Key
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	item, ok := lc.items[key]
	if ok {
		lc.queue.MoveToFront(item)
		return item.Value, ok
	}
	return nil, false
}

func (lc *lruCache) Set(key Key, value interface{}) bool {
	item, isset := lc.items[key]
	if isset {
		item.Value = value
		lc.queue.MoveToFront(item)
	} else {
		if lc.queue.Len() == lc.capacity {
			deleteItem := lc.queue.Back()
			lc.queue.Remove(deleteItem)
			keyDelete := lc.keys[deleteItem]
			delete(lc.items, keyDelete)
			delete(lc.keys, deleteItem)
		}
		item = lc.queue.PushFront(value)
	}
	lc.items[key] = item
	lc.keys[item] = key
	return isset
}

func (lc *lruCache) Clear() {
	lc.queue = NewList()
	lc.items = make(map[Key]*ListItem, lc.capacity)
	lc.keys = make(map[*ListItem]Key, lc.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		keys:     make(map[*ListItem]Key, capacity),
	}
}
