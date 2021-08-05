package lru

// 并发安全的lru cache

import (
	"container/list"
	"errors"
	"sync"
)

type Key interface{}

type Cache struct {
	MaxCapacity int
	linklist *list.List //链表存储LRU元素
	Callback func(key Key, value interface{})  // key失效时的回调
	cache map[interface{}]*list.Element  // map的key可以为任何类型，map的value是指向linklist的指针

	NeedLock bool   // 并发安全LRU需要锁
	lock sync.Mutex

	HitCnt int  // 命中的次数
	GetCnt int  // 请求读的次数
}

type Kv struct {
	k Key
	v interface{}
}

func NewLruCache(maxCapacity int, needLock bool, fn func(key Key, value interface{})) *Cache {
	if maxCapacity <= 0 {
		return nil
	}

	return &Cache {
		MaxCapacity:maxCapacity,
		linklist:list.New(),
		cache:make(map[interface{}]*list.Element),
		NeedLock:needLock,
		Callback:fn,
	}
}

// 添加元素到lru cache
func (c *Cache) Add(key Key, val interface{}) {
	if c.NeedLock {
		c.lock.Lock()
	}
	defer c.lock.Unlock()
	if c.cache == nil {
		c.linklist = list.New()
	}
	// 检查是否到达容量最大值，到达最大值就需要移除链表尾部的元素
	if c.Len() >= c.MaxCapacity {
		c.removeOldest()
	}
	// 把元素添加到头部
	elem := c.linklist.PushFront(&Kv{k:key, v:val})
	c.cache[key] = elem
}

// 从lru cache读取元素
func (c *Cache) Get(key Key) (val interface{}, err error) {
	if c.NeedLock {
		c.lock.Lock()
	}
	defer c.lock.Unlock()
	c.GetCnt += 1
	if c.cache == nil {
		return nil, errors.New("lru cache is not exist")
	}
	elem, ok := c.cache[key] 
	if !ok {
		return nil, errors.New("key not in cache")
	}
	c.HitCnt += 1
	c.linklist.MoveToFront(elem)
	return elem.Value.(*Kv).v, nil
}

func (c *Cache) Len() int {
	if c.cache == nil {
		return 0
	}
	return c.linklist.Len()
}

// 移除链表中最后的那个元素
func (c *Cache) removeOldest() {
	if c.cache == nil {
		return
	}
	oldest := c.linklist.Back()
	if oldest == nil {
		return
	}
	c.linklist.Remove(oldest)
	kv := oldest.Value.(*Kv)
	delete(c.cache, kv.k)
	if c.Callback != nil {
		c.Callback(kv.k, kv.v)
	}
}