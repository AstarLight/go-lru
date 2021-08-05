package main

import (
	"fmt"
	"go-lru/lru"
)

var hotCache *lru.Cache

func main() {
	stop := make(chan int)
	hotCache = lru.NewLruCache(3, true) // 启用一个容量为3，并发安全的LRU CACHE
	for i := 0; i < 100000; i++ {
		go func(index int) {
			key := fmt.Sprintf("key_%d", index%5)
			hotCache.Add(key, index)
			fmt.Printf("cache add, key=%+v, val=%+v\n", key, index)
			val,err := hotCache.Get(key)
			if err != nil {
				fmt.Printf("cache get error, key=%+v, err=%+v\n", key, err)
				return
			}
			fmt.Printf("cache get OK, key=%+v, val=%+v\n", key, val)
		}(i)
	}

	<- stop

}