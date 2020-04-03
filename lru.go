package lru

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

// Cache is an LRU cache. It is not safe for concurrent access.
type Cache struct {
	sync.RWMutex
	// MaxEntries is the maximum number of cache entries before
	// an item is evicted. Zero means no limit.
	MaxEntries int
	// OnEvicted optionally specifics a callback function to be
	// executed when an entry is purged from the cache.
	OnEvicted func(key Key, value interface{})

	ll     *list.List
	cache  map[interface{}]*list.Element
	expire int64
}

// A Key may be any value that is comparable. See http://golang.org/ref/spec#Comparison_operators
type Key interface{}
type entry struct {
	key    Key
	value  interface{}
	expire int64
}

// NewCache creates a new Cache.
// If maxEntries is zero, the cache has no limit and it's assumed
// that eviction is done by the caller.
func NewCache(maxEntries int, expired int64) *Cache {
	c := &Cache{
		MaxEntries: maxEntries,
		ll:         list.New(),
		cache:      make(map[interface{}]*list.Element),
		expire:     expired,
	}
	if c.expire > 0 {
		go c.cleanExpired()
	}
	return c
}

// Check whether entry is expired or not
func (e *entry) isExpired() bool {
	if e.expire == 0 { // entry without expire
		return false
	}
	if e.expire >= time.Now().Unix() {
		return false
	}
	return true
}

// cleans expired entries performing minimal checks
func (c *Cache) cleanExpired() {
	for {
		if c.ll.Len() == 0 {
			time.Sleep(time.Duration(c.expire) * time.Second)
			continue
		}
		ele := c.ll.Back()
		if ele.Value.(*entry).isExpired() {
			c.RLock()
			c.removeElement(ele)
			c.RUnlock()
		} else {
			time.Sleep(time.Duration(c.expire) * time.Second)
		}
	}
}

// Set a value to the cache.
// Key and value is required.
func (c *Cache) Set(key Key, value interface{}) (bool, error) {
	if c.cache == nil {
		return false, errors.New("cache is not initialized")
	}
	//
	c.Lock()
	var expire int64
	if c.expire > 0 {
		expire = time.Now().Unix() + c.expire
	}
	if ee, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ee)
		ee.Value.(*entry).value = value
		ee.Value.(*entry).expire = expire
		c.Unlock()
		return true, nil
	}
	ele := c.ll.PushFront(&entry{key, value, expire})
	c.cache[key] = ele
	if c.MaxEntries != 0 && c.ll.Len() > c.MaxEntries {
		c.removeOldest()
	}
	c.Unlock()
	return true, nil
}

// Get looks up a key's value from the cache.
func (c *Cache) Get(key Key) (value interface{}, ok bool) {
	c.Lock()
	defer c.Unlock()
	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		if ele.Value.(*entry).isExpired() {
			// delete expired elem
			c.removeElement(ele)
			return
		}
		c.ll.MoveToFront(ele)
		return ele.Value.(*entry).value, true
	}
	return
}

// Remove removes the provided key from the cache.
func (c *Cache) Remove(key Key) {
	c.Lock()
	defer c.Unlock()
	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		c.removeElement(ele)
	}
}

// removeOldest removes the oldest item from the cache.
func (c *Cache) removeOldest() {
	if c.cache == nil {
		return
	}
	ele := c.ll.Back()
	if ele != nil {
		c.removeElement(ele)
	}
}

func (c *Cache) removeElement(e *list.Element) {
	c.ll.Remove(e)
	kv := e.Value.(*entry)
	delete(c.cache, kv.key)
	if c.OnEvicted != nil {
		c.OnEvicted(kv.key, kv.value)
	}
}

// Len returns the number of items in the cache.
func (c *Cache) Len() int {
	c.RLock()
	defer c.RUnlock()
	if c.cache == nil {
		return 0
	}
	return c.ll.Len()
}

// Keys return all the keys in cache
func (c *Cache) Keys() []interface{} {
	c.RLock()
	defer c.RUnlock()
	keys := make([]interface{}, 0, c.ll.Len())
	for ele := c.ll.Front(); ele != nil; ele = ele.Next() {
		if ele.Value.(*entry).isExpired() {
			c.removeElement(ele)
		} else {
			keys = append(keys, ele.Value.(*entry).key)
		}
	}
	return keys
}

// Values return all the value in cache
func (c *Cache) Values() []interface{} {
	c.RLock()
	defer c.RUnlock()
	values := make([]interface{}, 0, c.ll.Len())
	for ele := c.ll.Front(); ele != nil; ele = ele.Next() {
		if ele.Value.(*entry).isExpired() {
			c.removeElement(ele)
		} else {
			values = append(values, ele.Value.(*entry).value)
		}
	}
	return values
}

// Flush remove all the keys in cache
func (c *Cache) Flush() {
	c.Lock()
	defer c.Unlock()
	c.ll = list.New()
	c.cache = make(map[interface{}]*list.Element)
}
