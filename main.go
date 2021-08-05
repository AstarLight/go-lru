package main

import (
	"fmt"
	"go-lru/lru"
	"time"
)

var hotCache *lru.Cache

// LRU key 被删除时的回调通知
func keyDeleteNotify(key lru.Key, val interface{}) {
	fmt.Printf("keyDeleteNotify, key=%+v, val=%+v\n", key, val)
}

func main() {
	hotCache = lru.NewLruCache(3, true, keyDeleteNotify) // 启用一个容量为3，并发安全的LRU CACHE
	for i := 0; i < 1000000; i++ {
		go func(index int) {
			key := fmt.Sprintf("key_%d", index%5)
			hotCache.Add(key, index)
			fmt.Printf("cache add, key=%+v, val=%+v\n", key, index)
			val, err := hotCache.Get(key)
			if err != nil {
				//fmt.Printf("cache get error, key=%+v, err=%+v\n", key, err)
				return
			}
			fmt.Printf("cache get OK, key=%+v, val=%+v\n", key, val)
		}(i)
	}
	time.Sleep(time.Second * time.Duration(10))
	fmt.Printf("cache GetCnt=%d, HitCnt=%d, HitPercent=%f\n", hotCache.GetCnt, hotCache.HitCnt, float64(hotCache.HitCnt)/float64(hotCache.GetCnt))
}
